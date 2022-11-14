package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(common.HeaderTssID) == "" {
		s.fullMsgHandler(w, r)
	} else {
		if !s.isAccessTokenCorrect(&r.Header) {
			http.Error(w, "wrong access token", http.StatusUnauthorized)
			return
		}
		s.registerConnection(w, r)
	}
}

func (s *Server) fullMsgHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug().Str("addr", r.RemoteAddr).Msg("Client ask FullMsg")
	if s.fullMsg == nil {
		http.Error(w, "full msg not set yet", http.StatusTooEarly)
		return
	}
	w.Write(s.fullMsg)
}

func (s *Server) registerConnection(w http.ResponseWriter, r *http.Request) {
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

	if !bytes.Equal(s.operation, operation) {
		connLogger.Info().Msg("This operation doesn't started by server")
		http.Error(w, "This operation doesn't started by server", http.StatusTooEarly)
		return
	}

	conn_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		connLogger.Error().Err(err).Msg("Failed to upgrade connection to websocket")
		return
	}

	// register connection (now ready for protocol)
	conn := &common.Conn{Conn: conn_}
	err = s.clientConnected(clientID, conn)

	if err != nil {
		connLogger.Error().Err(err).Msg("Failed to register connection")
		conn.Close(err)
		return
	}

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

func (s *Server) isAccessTokenCorrect(h *http.Header) bool {
	return h.Get(common.HeaderAccessToken) == s.accessToken
}
