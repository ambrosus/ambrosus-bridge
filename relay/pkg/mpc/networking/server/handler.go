package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug().Str("addr", r.RemoteAddr).Msg("New http connection")

	clientID, operation, err := parseHeaders(r)
	if err != nil {
		s.logger.Error().Err(err).Msg("parse headers error")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	connLogger := s.logger.With().Str("clientID", clientID).Logger()

	if bytes.Equal(operation, common.KeygenOperation) {
		connLogger.Info().Msg("Client wants to start keygen")
	} else {
		connLogger.Info().Msg("Client wants to start signing")
	}

	if !bytes.Equal(s.operation.SignMsg, operation) || !s.operation.Started {
		connLogger.Info().Msg("This operation doesn't started by server")
		http.Error(w, "This operation doesn't started by server", http.StatusTooEarly)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		connLogger.Error().Err(err).Msg("Failed to upgrade connection to websocket")
		return
	}

	// register connection (now ready for protocol)
	s.clientConnected(clientID, conn)
	defer s.clientDisconnected(clientID)

	result, err := s.receiveMsgs(conn, connLogger)
	if err != nil {
		connLogger.Error().Err(err).Msg("Server error")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		return
	}

	s.Lock()
	s.results[clientID] = result
	s.Unlock()
	s.resultsWaiter.Done()

	connLogger.Debug().Hex("result", result).Msg("Client finished protocol")
}

func parseHeaders(r *http.Request) (string, []byte, error) {
	clientID := r.Header.Get(common.HeaderTssID)
	if clientID == "" {
		return "", nil, fmt.Errorf("wrong clientID header")
	}

	operation, err := hex.DecodeString(r.Header.Get(common.HeaderTssOperation))
	if err != nil {
		return "", nil, fmt.Errorf("wrong operation header: %w", err)
	}

	return clientID, operation, err
}

func (s *Server) receiveMsgs(conn *websocket.Conn, logger zerolog.Logger) ([]byte, error) {
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("read protocol message: %w", err)
		}

		// client sends result message (end of protocol)
		if bytes.Index(msgBytes, common.ResultPrefix) == 0 { // msg starts with result prefix
			result := msgBytes[len(common.ResultPrefix):]
			return result, nil
		}

		msg := new(tss_wrap.OutputMessage)
		if err := msg.Unmarshal(msgBytes); err != nil {
			return nil, fmt.Errorf("unmarshal protocol message: %w", err)
		}

		// send client message to out sender service
		logger.Debug().Msg("Received msg")
		s.operation.OutCh <- msg
	}

}
