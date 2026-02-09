package quiz

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/quiz/create_quiz"
	"server/internal/feature/quiz/delete_quiz"
	"server/internal/feature/quiz/find_quiz"
	"server/internal/feature/quiz/get_by_id"
	"server/internal/feature/quiz/get_quiz"
	"server/internal/feature/quiz/get_users"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_quiz.New(p)
	getUC := get_quiz.New(p)
	getByIDUC := get_by_id.New(p)
	deleteUC := delete_quiz.New(p)
	findUC := find_quiz.New(p)
	getUsersUC := get_users.New(p)

	r.Post("/quizzes", create_quiz.HTTP(createUC))
	r.Get("/quizzes", get_quiz.HTTP(getUC))
	r.Get("/quizzes/{quiz_id}", get_by_id.HTTP(getByIDUC))
	r.Get("/quizzes/{quiz_id}/users", get_users.HTTP(getUsersUC))
	r.Get("/student/quizzes", find_quiz.HTTPForStudent(findUC))
	r.Delete("/quizzes", delete_quiz.HTTP(deleteUC))
}
