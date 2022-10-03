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
	// todo validate id
	if oldCon, ok := s.connections[id]; ok {
		oldCon.Close()
	}
	s.connections[id] = conn
	s.connChangeCh <- 1
}

func (s *Server) clientDisconnected(id string) {
	if oldCon, ok := s.connections[id]; ok {
		oldCon.Close()
		delete(s.connections, id)
	}
	s.connChangeCh <- 1
}
