package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/attempt/start_attempt"

	"github.com/go-chi/chi/v5"
)

func Route(router chi.Router, p *postgres.Postgres) {
	startUC := start_attempt.New(p)

	router.Post("/attempts", start_attempt.HTTP(startUC))
}
