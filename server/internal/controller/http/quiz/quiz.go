package quiz

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/quiz/create_quiz"
	"server/internal/feature/quiz/get_quiz"
	"server/internal/feature/quiz/list_quizzes"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_quiz.New(
		postgres.NewQuiz(p),
	)
	getUC := get_quiz.New(
		postgres.NewQuiz(p),
	)
	listUC := list_quizzes.New(
		postgres.NewQuiz(p),
	)
	/*
		getUC := get_quizzes.New(p)
		getByIDUC := get_by_id.New(p)
		deleteUC := delete_quiz.New(
			postgres.NewQuiz(p),
		)
		getUsersUC := get_users.New(p)
		startAttemptUC := start_attempt.New(
			p,
			postgres.NewQuiz(p),
			postgres.NewAttemt(p),
		)
	*/

	r.Route("/quizzes", func(r chi.Router) {
		r.Post("/", create_quiz.HTTP(createUC))
		r.Get("/", list_quizzes.HTTP(listUC))
		r.Route("/"+param.QuizID.Path(), func(r chi.Router) {
			r.Get("/", get_quiz.HTTP(getUC))
		})
		/*
			r.Get("/", get_quizzes.HTTP(getUC))
			r.Route("/"+param.QuizID.Path(), func(r chi.Router) {
				r.Get("/", get_by_id.HTTP(getByIDUC))
				r.Delete("/", delete_quiz.HTTP(deleteUC))
				r.Get("/users", get_users.HTTP(getUsersUC))
				r.Post("/attempt", start_attempt.HTTP(startAttemptUC))
			})
		*/
	})
}
