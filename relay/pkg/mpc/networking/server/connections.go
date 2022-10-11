package server

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/gorilla/websocket"
)

func (s *Server) waitForConnections() {
	s.logger.Debug().Msg("Wait for connections")
	for {
		if len(s.connections) < s.Tss.Threshold()-1 { // -1 coz server
			<-s.connChangeCh // wait for new connections
			continue
		}
		s.logger.Debug().Msg("All connections established")
		return
	}
}

func (s *Server) disconnectAll() {
	for _, conn := range s.connections {
		common.NormalClose(conn)
	}
}

func (s *Server) clientConnected(id string, conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	// todo validate id
	if oldCon, ok := s.connections[id]; ok {
		oldCon.Close()
	}
	s.connections[id] = conn
	s.connChangeCh <- 1

	s.logger.Debug().Str("id", id).Msg("Client connected")
}

func (s *Server) clientDisconnected(id string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.connections[id]; ok {
		//oldCon.Close()
		delete(s.connections, id)
	}
	s.connChangeCh <- 1

	s.logger.Debug().Str("id", id).Msg("Client disconnected")
}
