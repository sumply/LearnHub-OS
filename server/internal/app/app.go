package app

import (
	"net/http"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/rest"
	"server/internal/rest/middleware"
	"server/internal/service/generator"
	"server/internal/usecase"
)

func Run() error {
	logger.SetNewFunc(func() logger.Logger {
		return logger.NewFake()
	})
	logger.SetLayer(logger.DEBUG)

	storage := repository.NewStorage()
	repository.PrepareStorage(storage)
	repo := repository.New(
		repository.NewUserMemory(storage),
		repository.NewSubjectMemory(storage),
		repository.NewGroupMemory(storage),
		repository.NewQuizMemory(storage),
	)
	r, err := rest.NewRouter(
		usecase.NewUserReal(generator.NewStub(), repo),
		usecase.NewGroupReal(repo),
		usecase.NewRealSubject(repo),
		usecase.NewQuiz(repo),
		&middleware.TokenParserFake{},
	)
	if err != nil {
		return err
	}

	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
