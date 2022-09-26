package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
)

var keygenOperation = []byte("keygen")

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to upgrade connection to websocket")
		return
	}
	defer conn.Close()

	s.logger.Info().Str("address", conn.RemoteAddr().String()).Msg("New websocket connection")

	err = s.handler(conn)

	if err != nil {
		s.logger.Error().Err(err).Msg("Server error")

		errText := make([]byte, 120) // maxControlFramePayloadSize
		copy(errText, err.Error())

		if err = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(errText))); err != nil {
			s.logger.Error().Err(err).Msg("During sending error to client another error happened")
		}
	}
}

func (s *Server) handler(conn *websocket.Conn) (err error) {
	// ask client for ID
	_, clientIDBytes, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message 1: %w", err)
	}
	clientID := string(clientIDBytes)
	// register connection
	s.clientConnected(clientID, conn)
	defer s.clientDisconnected(clientID)

	// ask client for operation
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message 2: %w", err)
	}

	if bytes.Equal(msg, keygenOperation) {
		s.logger.Info().Msg("Client wants to start keygen")
	} else {
		s.logger.Info().Msg("Client wants to start signing")
	}

	if !bytes.Equal(s.operation.signMsg, keygenOperation) || !s.operation.started {
		return fmt.Errorf("This operation doesn't stated by server")
	}

	for {
		// todo is cycle break when client disconnect?

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read protocol message: %w", err)
		}

		msg := new(tss_wrap.OutputMessage)
		if err := msg.Unmarshal(msgBytes); err != nil {
			return fmt.Errorf("unmarshal protocol message: %w", err)
		}

		// send client message to out sender service
		s.operation.outCh <- msg
	}

}
