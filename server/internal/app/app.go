package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	if err := storage.Load(); err != nil {
		return err
	}

	signs := make(chan os.Signal, 1)
	signal.Notify(signs, os.Interrupt)
	go func() {
		<-signs
		fmt.Print("Saving storage data...")
		if err := storage.Save(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		os.Exit(0)
	}()
	repo := repository.New(
		repository.NewUserMemory(storage),
		repository.NewSubjectMemory(storage),
		repository.NewGroupMemory(storage),
		repository.NewQuizMemory(storage),
		repository.NewProgressMemory(storage),
	)
	r, err := rest.NewRouter(
		usecase.NewUserReal(generator.NewReal(), repo),
		usecase.NewGroupReal(repo),
		usecase.NewRealSubject(repo),
		usecase.NewQuiz(repo),
		usecase.NewProgressUsecase(repo),
		&middleware.TokenParserFake{},
	)
	if err != nil {
		return err
	}

	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
