package user

import (
	"server/internal/adapter/postgres"
	"server/internal/user/create_user"
	"server/internal/user/delete_user"
	"server/internal/user/get_user"
	"server/internal/user/update_user"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_user.New(p)
	getUC := get_user.New(p)
	deleteUC := delete_user.New(p)
	updateUC := update_user.New(p)

	r.Post("/users", create_user.HTTP(createUC))
	r.Get("/users", get_user.HTTP(getUC))
	r.Delete("/users/{user_id}", delete_user.HTTP(deleteUC))
	r.Patch("/users/{user_id}", update_user.HTTP(updateUC))
}
