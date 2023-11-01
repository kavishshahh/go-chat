package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWs(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client", ws.RemoteAddr())

	s.conns[ws] = true

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)

		if err != nil {
			delete(s.conns, ws)
			fmt.Println("Client disconnected", ws.RemoteAddr())
			return
		}

		fmt.Printf("recieved message from %s %s\n", ws.RemoteAddr(), message)
		websocket.Message.Send(ws, "Hi client!")

		for conn := range s.conns {
			if conn != ws {
				websocket.Message.Send(conn, message)
			}
		}
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWs))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err in listening", err)
	}
}
