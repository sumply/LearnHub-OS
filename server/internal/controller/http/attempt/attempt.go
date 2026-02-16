package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/attempt/finish_attempt"
	"server/internal/feature/attempt/get_by_id"
	"server/internal/pkg/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	finishUC := finish_attempt.New(p)
	getByIDUC := get_by_id.New(p)

	r.Route("/attempts", func(r chi.Router) {
		r.Route("/"+param.AttemptID.Path(), func(r chi.Router) {
			r.Post("/finish", finish_attempt.HTTP(finishUC))
			r.Get("/", get_by_id.HTTP(getByIDUC))
		})
	})
}
