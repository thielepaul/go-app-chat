package rpc

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"

	"github.com/coder/websocket"
)

type Backend struct {
	clients map[string]*Client
}

func NewBackend() *Backend {
	backend := &Backend{clients: make(map[string]*Client)}

	rpcServer, err := newServer(backend)
	if err != nil {
		log.Panicf("error creating rpc server: %v", err)
	}

	http.Handle(FrontendToBackendPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("new client f2b connected", r.RemoteAddr)
			wsConn, err := websocket.Accept(w, r, nil)
			if err != nil {
				log.Printf("error accepting websocket f2b: %s", err)
				return
			}
			defer wsConn.CloseNow()
			rpcServer.serveConn(wsConn)
			wsConn.Close(websocket.StatusNormalClosure, "")
			log.Println("client f2b disconnected", r.RemoteAddr)
		}))

	http.Handle(BackendToFrontendPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := websocket.Accept(w, r, nil)
			if err != nil {
				log.Printf("error accepting websocket b2f: %s", err)
				return
			}
			client := newClient(wsConn)
			backend.clients[r.RemoteAddr] = client
			log.Println("new client b2f connected", r.RemoteAddr)
		}))
	return backend
}

func (c *Backend) AddMessage(message string, _ *struct{}) error {
	for addr, client := range c.clients {
		if _, err := Call(client, (&Frontend{}).AddMessage, message); err != nil {
			if errors.Is(err, rpc.ErrShutdown) {
				log.Println("client b2f is disconnected", addr)
				delete(c.clients, addr)
			} else {
				log.Printf("error sending message to %s: %s", addr, err)
			}
		}
	}
	return nil
}
