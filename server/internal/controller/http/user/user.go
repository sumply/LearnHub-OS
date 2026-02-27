package user

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/user/create_user"
	"server/internal/feature/user/delete_user"
	"server/internal/feature/user/get_by_id"
	"server/internal/feature/user/get_quizzes"
	"server/internal/feature/user/get_user"
	"server/internal/feature/user/update_user"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_user.New(p)
	getUC := get_user.New(p)
	deleteUC := delete_user.New(p)
	updateUC := update_user.New(p)
	getQuizzesUC := get_quizzes.New(get_quizzes.Repository{
		Group:    p.MakeGroup(),
		QuizItem: p.MakeQuizItem(),
	})
	getByIDUC := get_by_id.New(p)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", create_user.HTTP(createUC))
		r.Get("/", get_user.HTTP(getUC))
		r.Route("/"+param.UserID.Path(), func(r chi.Router) {
			r.Delete("/", delete_user.HTTP(deleteUC))
			r.Put("/", update_user.HTTP(updateUC))
			r.Get("/", get_by_id.HTTP(getByIDUC))
			r.Get("/quizzes", get_quizzes.HTTP(getQuizzesUC))
		})
	})
}
