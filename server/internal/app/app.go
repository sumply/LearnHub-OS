package app

import (
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/transport"
	"server/internal/usecase"
)

func Run() error {
	r, err := transport.NewRouter(
		usecase.NewFakeUser(),
		&usecase.FakeGroup{},
		&usecase.FakeSubject{},
		&transport.FakeTokenParser{},
	)
	if err != nil {
		return err
	}
	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Addr, s.Port), r)
}
