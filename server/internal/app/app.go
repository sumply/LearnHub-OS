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

	repo := repository.New(
		&repository.UserMemory{},
		repository.NewSubjectStub(),
		repository.NewGroupStub(),
		repository.NewSpecialityStub(),
	)
	r, err := rest.NewRouter(
		usecase.NewRealUser(generator.NewStub(), repo),
		usecase.NewGroupStub(),
		usecase.NewRealSubject(),
		usecase.NewRealSpeciality(),
		&middleware.StubTokenParser{},
	)
	if err != nil {
		return err
	}

	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
