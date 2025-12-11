package app

import (
	"net/http"
	"server/internal/config"
	"server/internal/transport"
	"server/internal/usecase"
)

func Run() error {
	r, err := transport.NewRouter(
		usecase.NewFakeUser(),
		usecase.NewFakeGroup(),
		usecase.NewFakeSubject(),
		usecase.NewFakeQuiz(),
		&transport.FakeTokenParser{},
	)
	if err != nil {
		return err
	}
	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
