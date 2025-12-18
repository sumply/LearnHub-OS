package app

import (
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/transport"
	"server/internal/usecase"
)

func Run() error {
	app, err := config.NewApp("./configs/app.yaml")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", app)

	logger.SetNewFunc(func() logger.Logger {
		return logger.NewFake()
	})

	uu, err := app.CreateUsecaseUser()
	if err != nil {
		return err
	}
	ua, err := app.CreateUsecaseAnswer()
	if err != nil {
		return err
	}

	r, err := transport.NewRouter(
		uu,
		usecase.NewStubGroup(),
		usecase.NewStubSubject(),
		usecase.NewStubQuiz(),
		ua,
		&transport.StubTokenParser{},
	)
	if err != nil {
		return err
	}
	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
