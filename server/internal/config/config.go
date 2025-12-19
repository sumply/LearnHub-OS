package config

import (
	"errors"
	"fmt"
	"os"
	"server/internal/logger"
	"server/internal/transport"
	"server/internal/usecase"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmpty   = errors.New("key is empty")
	ErrInvalid = errors.New("value is invalid")
)

type Server struct {
	Addr string
	Port int
}

func (s *Server) String() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

func NewServerFromEnv() (*Server, error) {
	addr := os.Getenv("ADDRESS_HOST")
	if addr == "" {
		return nil, fmt.Errorf("ADDRESS_HOST: %w", ErrEmpty)
	}
	port := os.Getenv("PORT_HOST")
	if port == "" {
		return nil, fmt.Errorf("PORT_HOST: %w", ErrEmpty)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("PORT_HOST: %w", err)
	}
	return &Server{
		Addr: addr,
		Port: p,
	}, nil
}

const (
	ImplStub = "stub"
	ImplFake = "fake"
)

type AppCommon struct {
	Logger string `yaml:"logger"`
}

type AppUsecase struct {
	User    string `yaml:"user"`
	Quiz    string `yaml:"quiz"`
	Group   string `yaml:"group"`
	Answer  string `yaml:"answer"`
	Subject string `yaml:"subject"`
}

type AppJWT struct {
	Parser string `yaml:"parser"`
}

type AppTransport struct {
	JWT AppJWT `yaml:"jwt"`
}

type App struct {
	Common    AppCommon    `yaml:"common"`
	Usecase   AppUsecase   `yaml:"usecase"`
	Transport AppTransport `yaml:"transport"`
}

func NewApp(path string) (*App, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var app App
	if err := yaml.NewDecoder(f).Decode(&app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *App) CreateUsecaseUser() (usecase.User, error) {
	switch a.Usecase.User {
	case ImplStub:
		return usecase.NewStubUser(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateUsecaseSubject() (usecase.Subject, error) {
	switch a.Usecase.User {
	case ImplStub:
		return usecase.NewStubSubject(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateUsecaseGroup() (usecase.Group, error) {
	switch a.Usecase.User {
	case ImplStub:
		return usecase.NewStubGroup(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateUsecaseAnswer() (usecase.Answer, error) {
	switch a.Usecase.User {
	case ImplStub:
		return usecase.NewStubAnswer(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateUsecaseQuiz() (usecase.Quiz, error) {
	switch a.Usecase.User {
	case ImplStub:
		return usecase.NewStubQuiz(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateLoggerNewFunc() (logger.NewFunc, error) {
	switch a.Common.Logger {
	case ImplFake:
		return func() logger.Logger {
			return logger.NewFake()
		}, nil
	default:
		return nil, nil
	}
}

func (a *App) CreateTransportJWTParser() (transport.TokenParser, error) {
	switch a.Transport.JWT.Parser {
	case ImplFake:
		return nil, nil
	case ImplStub:
		return &transport.StubTokenParser{}, nil
	default:
		return nil, nil
	}
}
