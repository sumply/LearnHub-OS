package config

import (
	"fmt"
	"os"
	"server/archive/internal/auth"
	"server/archive/internal/logger"
	"server/archive/internal/service/client"
	"strconv"
	"time"
)

const (
	JWT_SECRET_KEY       = "JWT_SECRET_KEY"
	JWT_ACCESS_DURATION  = "JWT_ACCESS_DURATION"
	JWT_REFRESH_DURATION = "JWT_REFRESH_DURATION"
	JWT_ISSUER           = "JWT_ISSUER"
	SMTP_FROM            = "SMTP_FROM"
	SMTP_PASSWORD        = "SMTP_PASSWORD"
	SMTP_HOST            = "SMTP_HOST"
	SMTP_PORT            = "SMTP_PORT"
	LOG_LEVEL            = "LOG_LEVEL"
)

type Server struct {
	Addr string
	Port int
}

func (s *Server) String() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

type ENV struct {
	JWTSecretKey  []byte
	JWTAccessDur  time.Duration
	JWTRefreshDur time.Duration
	JWTIssur      string
	SMTPEmailFrom string
	SMTPPassword  string
	SMTPHost      string
	SMTPPort      string
	LOGLevel      string
}

func NewENV() *ENV {
	return &ENV{
		JWTSecretKey:  getByte(JWT_SECRET_KEY),
		JWTAccessDur:  getTimeDuration(JWT_ACCESS_DURATION),
		JWTRefreshDur: getTimeDuration(JWT_REFRESH_DURATION),
		JWTIssur:      getValue(JWT_ISSUER),
		SMTPEmailFrom: getValue(SMTP_FROM),
		SMTPPassword:  getValue(SMTP_PASSWORD),
		SMTPHost:      getValue(SMTP_HOST),
		SMTPPort:      getValue(SMTP_PORT),
		LOGLevel:      getValue(LOG_LEVEL),
	}
}

func (env *ENV) CreateAuthJWT() *auth.JWT {
	return auth.NewJWT(
		env.JWTSecretKey,
		env.JWTAccessDur,
		env.JWTRefreshDur,
		env.JWTIssur,
	)
}

func (env *ENV) CreateSMTPClient() client.SMTP {
	return client.NewSMTPClient(
		env.SMTPEmailFrom,
		env.SMTPPassword,
		env.SMTPHost,
		env.SMTPPort,
	)
}

func (env *ENV) InitLogger() {
	logger.SetNewFunc(func() logger.Logger {
		return logger.NewFake()
	})
	switch env.LOGLevel {
	case "DEBUG":
		logger.SetLayer(logger.DEBUG)
	case "INFO":
		logger.SetLayer(logger.INFO)
	case "WARN":
		logger.SetLayer(logger.WARN)
	case "ERROR":
		logger.SetLayer(logger.ERROR)
	default:
		panic(fmt.Errorf("%s value is invalid: %s", LOG_LEVEL, getValue(LOG_LEVEL)))
	}
}

func getTimeDuration(key string) time.Duration {
	d := getValue(key)
	dInt, err := strconv.Atoi(d)
	if err != nil {
		panic(fmt.Errorf("%s: %w", key, err))
	}
	return time.Minute * time.Duration(dInt)
}

func getByte(key string) []byte {
	return []byte(getValue(key))
}

func getValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("%s is empty", key))
	}
	return value
}
