package rpc

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/coder/websocket"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Frontend struct {
	*Client
	Messages []string
	Ctx      app.Context
}

func NewFrontend(ctx app.Context, remoteUrl *url.URL) (*Frontend, error) {
	address := "ws://" + remoteUrl.Host
	if remoteUrl.Scheme == "https" {
		address = "wss://" + remoteUrl.Host
	}

	frontendToBackendUrl, err := url.JoinPath(address, FrontendToBackendPath)
	if err != nil {
		return nil, fmt.Errorf("error joining url: %w", err)
	}
	frontendToBackendConn, _, err := websocket.Dial(context.Background(), frontendToBackendUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket f2b: %w", err)
	}
	log.Println("rpc client connected")

	backendToFrontendUrl, err := url.JoinPath(address, BackendToFrontendPath)
	if err != nil {
		return nil, fmt.Errorf("error joining url: %w", err)
	}
	backendToFrontendConn, _, err := websocket.Dial(context.Background(), backendToFrontendUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket b2f: %w", err)
	}

	frontend := &Frontend{
		Ctx:    ctx,
		Client: newClient(frontendToBackendConn),
	}

	server, err := newServer(frontend)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc server: %w", err)
	}
	go server.serveConn(backendToFrontendConn)
	log.Println("rpc server connected")

	return frontend, nil
}

func (c *Frontend) AddMessage(message string, _ *struct{}) error {
	log.Println("addmessage:", message)
	c.Messages = append(c.Messages, message)
	c.Ctx.Update()
	return nil
}
