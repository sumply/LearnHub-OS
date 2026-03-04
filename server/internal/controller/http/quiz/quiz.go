package quiz

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/attempt/start_attempt"
	"server/internal/feature/quiz/create_quiz"
	"server/internal/feature/quiz/delete_quiz"
	"server/internal/feature/quiz/get_by_id"
	"server/internal/feature/quiz/get_quizzes"
	"server/internal/feature/quiz/get_users"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_quiz.New(
		postgres.NewTeacher(p),
		nil,
	)
	getUC := get_quizzes.New(p)
	getByIDUC := get_by_id.New(p)
	deleteUC := delete_quiz.New(p)
	getUsersUC := get_users.New(p)
	startAttemptUC := start_attempt.New(p)

	r.Route("/quizzes", func(r chi.Router) {
		r.Post("/", create_quiz.HTTP(createUC))
		r.Get("/", get_quizzes.HTTP(getUC))
		r.Route("/"+param.QuizID.Path(), func(r chi.Router) {
			r.Get("/", get_by_id.HTTP(getByIDUC))
			r.Delete("/", delete_quiz.HTTP(deleteUC))
			r.Get("/users", get_users.HTTP(getUsersUC))
			r.Post("/attempt", start_attempt.HTTP(startAttemptUC))
		})
	})
}
