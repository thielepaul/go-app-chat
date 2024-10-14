package rpc

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Backend struct {
	clients []*Client
}

func NewBackend() *Backend {
	backend := &Backend{}

	rpcServer, err := newServer(backend)
	if err != nil {
		log.Panicf("error creating rpc server: %v", err)
	}

	http.Handle(FrontendToBackendPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := websocket.Accept(w, r, nil)
			if err != nil {
				log.Printf("error accepting websocket f2b: %s", err)
				return
			}
			defer wsConn.CloseNow()
			rpcServer.serveConn(wsConn)
			wsConn.Close(websocket.StatusNormalClosure, "")
		}))

	http.Handle(BackendToFrontendPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := websocket.Accept(w, r, nil)
			if err != nil {
				log.Printf("error accepting websocket b2f: %s", err)
				return
			}
			client := newClient(wsConn)
			backend.clients = append(backend.clients, client)
			log.Println("new client connected")
		}))
	return backend
}

func (c *Backend) AddMessage(message string, _ *struct{}) error {
	for _, client := range c.clients {
		_, _ = Call(client, (&Frontend{}).AddMessage, message)
		// ignore errors because each we do not handle connection close correctly yet
	}
	return nil
}
