package app

import (
	"net/http"
	"server/internal/config"

	controller "server/internal/controller/http"
)

func Run() error {
	env := &config.Env{}

	addr := env.CreateServeAddress()

	return http.ListenAndServe(addr.String(), controller.Router(env))
}
