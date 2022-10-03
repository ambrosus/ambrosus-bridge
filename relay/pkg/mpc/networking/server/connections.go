package server

import "github.com/gorilla/websocket"

func (s *Server) waitForConnections() {
	for {
		if len(s.connections) < s.Tss.Threshold()-1 { // -1 coz server
			<-s.connChangeCh // wait for new connections
			continue
		}
		return
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
