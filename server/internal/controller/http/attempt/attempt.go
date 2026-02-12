package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/attempt/finish_attempt"
	"server/internal/feature/attempt/start_attempt"
	"server/internal/pkg/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	startUC := start_attempt.New(p)
	finishUC := finish_attempt.New(p)

	r.Route("/attempts", func(r chi.Router) {
		r.Post("/", start_attempt.HTTP(startUC))
		r.Route("/"+param.AttemptID.Path(), func(r chi.Router) {
			r.Post("/finish", finish_attempt.HTTP(finishUC))
		})
	})
}
