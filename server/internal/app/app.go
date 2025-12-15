package app

import (
	"net/http"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/transport"
	"server/internal/usecase"
)

func Run() error {
	logger.SetNewFunc(func() logger.Logger {
		return logger.NewMock()
	})

	r, err := transport.NewRouter(
		usecase.NewFakeUser(),
		usecase.NewFakeGroup(),
		usecase.NewFakeSubject(),
		usecase.NewFakeQuiz(),
		usecase.NewFakeQuizResult(),
		&transport.FakeTokenParser{},
	)
	if err != nil {
		return err
	}
	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
