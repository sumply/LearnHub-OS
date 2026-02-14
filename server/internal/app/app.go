package app

import (
	"net/http"
	"server/internal/config"
	myhttp "server/internal/controller/http"
)

func Run() error {
	env := &config.Env{}

	addr := env.CreateServeAddress()

	return http.ListenAndServe(addr.String(), myhttp.Router(env))
}
