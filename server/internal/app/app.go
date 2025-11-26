package app

import (
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/transport"
)

func Run() error {
	r := transport.NewServer(InitHandler())
	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Addr, s.Port), r)
	return nil
}

func InitHandler() transport.HandlerInterface {
	return transport.NewMockHandler()
}
