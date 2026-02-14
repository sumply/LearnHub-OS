package config

import (
	"fmt"
	"os"
	"server/internal/adapter/postgres"
	"strconv"
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

type Creator interface {
	CreatePostgresConnection() PostgresConnection
	CreateServeAddress() ServeAddress
}
