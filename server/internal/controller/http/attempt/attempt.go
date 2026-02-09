package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/attempt/finish_attempt"
	"server/internal/feature/attempt/start_attempt"

	"github.com/go-chi/chi/v5"
)

func Route(router chi.Router, p *postgres.Postgres) {
	startUC := start_attempt.New(p)
	finishUC := finish_attempt.New(p)

	router.Post("/attempts", start_attempt.HTTP(startUC))
	router.Post("/attempts/{attempt_id}", finish_attempt.HTTP(finishUC))
}
