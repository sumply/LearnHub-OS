package auth

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/auth/login"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	loginUC := login.New(p)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", login.HTTP(loginUC))
	})
}
