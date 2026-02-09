package group

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/group/create_group"
	"server/internal/feature/group/delete_group"
	"server/internal/feature/group/get_group"
	"server/internal/feature/group/update_group"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_group.New(p)
	getUC := get_group.New(p)
	updateUC := update_group.New(p)
	deleteUC := delete_group.New(p)

	r.Post("/groups", create_group.HTTP(createUC))
	r.Get("/groups", get_group.HTTP(getUC))
	r.Put("/groups/{group_id}", update_group.HTTP(updateUC))
	r.Delete("/groups/{group_id}", delete_group.HTTP(deleteUC))
}
