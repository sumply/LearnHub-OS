package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/attempt/finish_attempt"
	"server/internal/feature/attempt/get_attempt"
	"server/internal/feature/attempt/list_attempts"
	"server/internal/feature/attempt/start_attempt"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	startUC := start_attempt.New(
		postgres.NewQuiz(p),
		postgres.NewAttempt(p),
	)
	getUC := get_attempt.New(
		postgres.NewAttempt(p),
	)
	listUC := list_attempts.New(
		postgres.NewAttempt(p),
	)
	finishUC := finish_attempt.New(
		postgres.NewAttempt(p),
		postgres.NewQuiz(p),
	)

	r.Route("/attempts", func(r chi.Router) {
		r.Post("/", start_attempt.HTTP(startUC))
		r.Get("/", list_attempts.HTTP(listUC))
		r.Route("/"+param.AttemptID.Path(), func(r chi.Router) {
			r.Get("/", get_attempt.HTTP(getUC))
			r.Route("/finish", func(r chi.Router) {
				r.Post("/", finish_attempt.HTTP(finishUC))
			})
		})
	})
}
