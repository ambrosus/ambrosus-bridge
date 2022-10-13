package server

import (
	"bytes"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"golang.org/x/sync/errgroup"
)

// transmitter send messages from outCh to all clients.
// returns error if any client has error
// returns nil when outCh is closed
func (s *Server) transmitter(outCh chan *tss_wrap.OutputMessage, inCh chan []byte) error {
	s.logger.Debug().Msg("Start transmitter")
	for msg := range outCh {
		err := s.sendMsg(msg, inCh)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to send message")
			return err
			// todo repeat on err
		}
	}
	return nil
}

// sendMsg send message to own Tss or to another client(s)
func (s *Server) sendMsg(msg *tss_wrap.OutputMessage, toMeCh chan []byte) error {
	if msg == nil || msg.SendToIds == nil {
		return fmt.Errorf("nil message")
	}

	for _, id := range msg.SendToIds {
		s.logger.Debug().Str("To", id).Msg("Send message to client")

		// send to own tss
		if id == s.Tss.MyID() {
			toMeCh <- msg.Message
			s.logger.Debug().Str("To", id).Msg("Send message to myself successfully")
			continue
		}

		// send to another client
		conn, ok := s.connections[id]
		if !ok {
			s.logger.Warn().Msgf("Connection with id %s not found, ", id)
			return fmt.Errorf("connection %v not found", id)
			// todo maybe call waitForConnections on this error
		}

		if err := conn.Write(msg.Message); err != nil {
			return fmt.Errorf("writeMessage: %w", err)
		}
		s.logger.Debug().Str("To", id).Msg("Send message to client sucessfully")
	}

	return nil
}

// receiver receive messages from all clients.
// returns error if any client has error
// returns nil when all clients send result message
func (s *Server) receiver(ch chan *tss_wrap.OutputMessage) error {
	s.logger.Debug().Msg("Start receiver")

	eg := new(errgroup.Group)
	for id_ := range s.connections {
		id := id_
		eg.Go(func() error {
			result, err := s.receiveMsgs(id, ch)
			if err != nil {
				s.logger.Error().Str("clientID", id).Err(err).Msg("Receive message from client error")
				return fmt.Errorf("receive message from client %s: %w", id, err)
			}

			s.logger.Debug().Str("clientID", id).Hex("result", result).Msg("Client finished protocol")

			s.Lock()
			s.results[id] = result
			s.Unlock()

			return nil
		})
	}

	return eg.Wait()
}

func (s *Server) receiveMsgs(clientID string, outCh chan *tss_wrap.OutputMessage) ([]byte, error) {
	s.logger.Debug().Str("clientID", clientID).Msg("Start receive messages from client")
	conn := s.connections[clientID]
	for {
		msgBytes, err := conn.Read()
		if err != nil {
			return nil, fmt.Errorf("read protocol message: %w", err)
		}

		// client sends result message (end of protocol)
		if bytes.Index(msgBytes, common.ResultPrefix) == 0 { // msg starts with result prefix
			result := msgBytes[len(common.ResultPrefix):]
			s.logger.Debug().Str("ClientID", clientID).Msg("Receive result message from client")
			return result, nil
		}
		s.logger.Debug().Str("ClientID", clientID).Msg("Receive message from client")

		msg := new(tss_wrap.OutputMessage)
		if err := msg.Unmarshal(msgBytes); err != nil {
			return nil, fmt.Errorf("unmarshal protocol message: %w", err)
		}

		// send client message to out sender service
		outCh <- msg
	}

}
