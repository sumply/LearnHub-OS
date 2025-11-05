package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrEmpty   = errors.New("key is empty")
	ErrInvalid = errors.New("value is invalid")
)

type Server struct {
	Addr string
	Port int
}

func NewServer() (*Server, error) {
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
