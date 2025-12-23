package config

import (
	"errors"
	"fmt"
	"os"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/generator"
	"server/internal/service/sender"
	"server/internal/service/validator"
	"server/internal/transport"
	"server/internal/usecase"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmpty   = errors.New("key is empty")
	ErrInvalid = errors.New("value is invalid")
)

const (
	ImplStub  = "stub"
	ImplFake  = "fake"
	ImplReal  = "real"
	ImplDebug = "debug"
	ImplInfo  = "info"
	ImplWarn  = "warn"
	ImplError = "error"
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

type AppCommonLogger struct {
	Type  string `yaml:"type"`
	Level string `yaml:"level"`
}

type AppCommon struct {
	Logger AppCommonLogger `yaml:"logger"`
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

type AppGenerator struct {
	Auth string `yaml:"auth"`
	Page string `yaml:"page"`
}

func (a *AppGenerator) CreateAuth() (generator.Generator, error) {
	switch a.Auth {
	case ImplStub:
		return generator.NewStub(), nil
	default:
		return nil, fmt.Errorf("%w: generator .auth", ErrInvalid)
	}
}

func (a *AppGenerator) CreatePage() (generator.PageGenerator, error) {
	switch a.Page {
	case ImplStub:
		return generator.NewStubPageGenerator(), nil
	default:
		return nil, fmt.Errorf("%w: generator .page", ErrInvalid)
	}
}

type AppSender struct {
	Mail string `yaml:"mail"`
}

func (a *AppSender) CreateMail() (sender.Mail, error) {
	switch a.Mail {
	case ImplStub:
		return sender.NewStubMail(), nil
	default:
		return nil, fmt.Errorf("%w: sender .mail", ErrInvalid)
	}
}

type AppValidator struct {
	User string `yaml:"user"`
}

func (a *AppValidator) CreateUser() (validator.User, error) {
	switch a.User {
	case ImplStub:
		return validator.NewStubUser(), nil
	default:
		return nil, fmt.Errorf("%w: validator .user", ErrInvalid)
	}
}

type AppStorage struct {
	Type string `yaml:"type"`
}

func (a *AppStorage) CreateStorage() (repository.Repository, error) {
	switch a.Type {
	case ImplStub:
		return repository.NewStub(), nil
	default:
		return nil, fmt.Errorf("%w: storage .type", ErrInvalid)
	}
}

type AppService struct {
	Generator AppGenerator `yaml:"generator"`
	Sender    AppSender    `yaml:"sender"`
	Storage   AppStorage   `yaml:"storage"`
	Validator AppValidator `yaml:"validator"`
}

type App struct {
	Common    AppCommon    `yaml:"common"`
	Usecase   AppUsecase   `yaml:"usecase"`
	Service   AppService   `yaml:"service"`
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
	case ImplReal:
		gen, err := a.Service.Generator.CreateAuth()
		if err != nil {
			return nil, err
		}
		val, err := a.Service.Validator.CreateUser()
		if err != nil {
			return nil, err
		}
		ms, err := a.Service.Sender.CreateMail()
		if err != nil {
			return nil, err
		}
		pgen, err := a.Service.Generator.CreatePage()
		if err != nil {
			return nil, err
		}
		storage, err := a.Service.Storage.CreateStorage()
		if err != nil {
			return nil, err
		}
		return usecase.NewRealUser(gen, val, ms, pgen, storage), nil
	default:
		return nil, fmt.Errorf("%w: usecase .user", ErrInvalid)
	}
}

func (a *App) CreateUsecaseSubject() (usecase.Subject, error) {
	switch a.Usecase.Subject {
	case ImplStub:
		return usecase.NewStubSubject(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, fmt.Errorf("%w: usecase .subject", ErrInvalid)
	}
}

func (a *App) CreateUsecaseGroup() (usecase.Group, error) {
	switch a.Usecase.Group {
	case ImplStub:
		return usecase.NewStubGroup(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, fmt.Errorf("%w: usecase .group", ErrInvalid)
	}
}

func (a *App) CreateUsecaseAnswer() (usecase.Answer, error) {
	switch a.Usecase.Answer {
	case ImplStub:
		return usecase.NewStubAnswer(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, fmt.Errorf("%w: usecase .answer", ErrInvalid)
	}
}

func (a *App) CreateUsecaseQuiz() (usecase.Quiz, error) {
	switch a.Usecase.Quiz {
	case ImplStub:
		return usecase.NewStubQuiz(), nil
	case ImplFake:
		return nil, nil
	default:
		return nil, fmt.Errorf("%w: usecase .quiz", ErrInvalid)
	}
}

func (a *App) CreateLoggerNewFunc() (logger.NewFunc, error) {
	switch a.Common.Logger.Type {
	case ImplFake:
		return func() logger.Logger {
			return logger.NewFake()
		}, nil
	default:
		return nil, fmt.Errorf("%w: logger .type", ErrInvalid)
	}
}

func (a *App) CreateTransportJWTParser() (transport.TokenParser, error) {
	switch a.Transport.JWT.Parser {
	case ImplFake:
		return nil, nil
	case ImplStub:
		return &transport.StubTokenParser{}, nil
	default:
		return nil, fmt.Errorf("%w: logger .parse", ErrInvalid)
	}
}

func (a *App) CreateLoggerLevel() (logger.Level, error) {
	switch a.Common.Logger.Level {
	case ImplDebug:
		return logger.DEBUG, nil
	case ImplInfo:
		return logger.INFO, nil
	case ImplWarn:
		return logger.WARN, nil
	case ImplError:
		return logger.ERROR, nil
	default:
		return 0, ErrInvalid
	}
}
