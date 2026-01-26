package quiz

import (
	"server/internal/adapter/postgres"
	"server/internal/quiz/create_quiz"
	"server/internal/quiz/delete_quiz"
	"server/internal/quiz/get_quiz"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_quiz.New(p)
	getUC := get_quiz.New(p)
	deleteUC := delete_quiz.New(p)

	r.Post("/quizzes", create_quiz.HTTP(createUC))
	r.Get("/quizzes", get_quiz.HTTP(getUC))
	r.Delete("/quizzes", delete_quiz.HTTP(deleteUC))
}
