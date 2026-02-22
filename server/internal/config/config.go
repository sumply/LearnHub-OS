package config

import (
	"fmt"
	"os"
	"server/internal/adapter/postgres"
	"strconv"
	"time"
)

type ServeAddress struct {
	Host string
	Port int
}

func (s ServeAddress) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type PostgresConnection struct {
	User     string
	Password string
	Database string
	SSLMode  string
	Port     int
	Host     string
}

func (p *PostgresConnection) CreateOptions() postgres.Options {
	return postgres.Options{
		User:     p.User,
		Password: p.Password,
		DB:       p.Database,
		SSLMode:  p.SSLMode,
		Port:     p.Port,
		Host:     p.Host,
	}
}

type JWT struct {
	Secret     []byte
	Issuer     string
	AccessDur  time.Duration
	RefreshDur time.Duration
}

type Env struct{}

func (e Env) CreatePostgresConnection() PostgresConnection {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	return PostgresConnection{
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslmode,
		Host:     host,
		Port:     portInt,
	}
}

func (e Env) CreateServeAddress() ServeAddress {
	host := os.Getenv("SERVE_HOST")
	port := os.Getenv("SERVE_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return ServeAddress{
		Host: host,
		Port: portInt,
	}
}

func (e Env) CreateJWT() JWT {
	issuer := os.Getenv("JWT_PAYLOAD_ISSUER")
	secret := os.Getenv("JWT_SECRET_KEY")
	aDurStr := os.Getenv("JWT_ACCESS_DURATION")
	aDur, err := strconv.Atoi(aDurStr)
	if err != nil {
		panic(err)
	}
	rDurStr := os.Getenv("JWT_REFRESH_DURATION")
	rDur, err := strconv.Atoi(rDurStr)
	if err != nil {
		panic(err)
	}
	return JWT{
		Issuer:     issuer,
		Secret:     []byte(secret),
		AccessDur:  time.Duration(aDur),
		RefreshDur: time.Duration(rDur),
	}
}

type Creator interface {
	CreatePostgresConnection() PostgresConnection
	CreateServeAddress() ServeAddress
	CreateJWT() JWT
}
