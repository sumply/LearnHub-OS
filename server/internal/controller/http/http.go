package http

import (
	"net/http"
	"server/internal/adapter/postgres"
	"server/internal/controller/http/group"
	"server/internal/controller/http/subject"
	"server/internal/controller/http/user"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewMux()

	p, err := postgres.New(postgres.Options{
		User:     "postgres",
		Password: "2121",
		DB:       "test",
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}

	user.Route(r, p)
	group.Route(r, p)
	subject.Route(r, p)

	return r
}
