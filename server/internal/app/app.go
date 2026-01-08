package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/repository"
	"server/internal/rest"
	"server/internal/service/generator"
	"server/internal/usecase"
	"syscall"
)

func Run() error {
	env := config.NewENV()
	env.InitLogger()

	jwt := env.CreateAuthJWT()
	domain.InitGenerateTokenPair(jwt.GenerateTokenPair)
	client := env.CreateSMTPClient()

	storage := repository.NewStorage()
	if err := storage.Load(); err != nil {
		return err
	}

	signs := make(chan os.Signal, 1)
	signal.Notify(signs, os.Interrupt, syscall.SIGTERM)
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
		usecase.NewUserReal(generator.NewReal(), repo, client),
		usecase.NewGroupReal(repo),
		usecase.NewRealSubject(repo),
		usecase.NewQuiz(repo),
		usecase.NewProgressUsecase(repo),
		jwt,
	)
	if err != nil {
		return err
	}

	s := config.Server{Addr: "0.0.0.0", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}
