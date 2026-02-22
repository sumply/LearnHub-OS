package auth

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/auth/login"
	"server/internal/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, g *jwt.Generator, p *postgres.Postgres) {
	loginUC := login.New(g, p)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", login.HTTP(loginUC))
	})
}
