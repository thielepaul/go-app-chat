package rpc

import (
	"context"
	"log"
	"net/rpc"
	"reflect"
	"runtime"
	"strings"

	"github.com/coder/websocket"
)

const FrontendToBackendPath = "/frontend2backend"
const BackendToFrontendPath = "/backend2frontend"

type Client struct {
	client *rpc.Client
	wsConn *websocket.Conn
}

func newClient(wsConn *websocket.Conn) *Client {
	return &Client{
		client: rpc.NewClient(websocket.NetConn(context.Background(), wsConn, websocket.MessageBinary)),
		wsConn: wsConn,
	}
}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return c.wsConn.CloseNow()
}

func Call[T1 any, T2 any](client *Client, call func(argType T1, replyType *T2) error, arg T1) (T2, error) {
	name := runtime.FuncForPC(reflect.ValueOf(call).Pointer()).Name()
	name = strings.NewReplacer("(", "", ")", "", "*", "").Replace(name)
	name = strings.TrimSuffix(name, "-fm")
	var reply T2
	err := client.client.Call(name, arg, &reply)
	return reply, err
}

func register[T any](serv *rpc.Server, service T) error {
	typ := reflect.TypeOf(service)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	name := typ.PkgPath() + "." + typ.Name()
	log.Print("Registering ", name)
	return serv.RegisterName(name, service)
}

type Server struct {
	server *rpc.Server
}

func newServer(types ...interface{}) (*Server, error) {
	rpcServer := rpc.NewServer()
	for _, typ := range types {
		if err := register(rpcServer, typ); err != nil {
			return nil, err
		}
	}
	return &Server{server: rpcServer}, nil
}

func (s *Server) serveConn(conn *websocket.Conn) {
	s.server.ServeConn(websocket.NetConn(context.Background(), conn, websocket.MessageBinary))
}
