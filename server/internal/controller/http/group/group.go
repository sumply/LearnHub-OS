package group

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/group/add_students"
	"server/internal/feature/group/create_group"
	"server/internal/feature/group/delete_group"
	"server/internal/feature/group/get_group"
	"server/internal/feature/group/get_students"
	"server/internal/feature/group/update_group"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_group.New(
		postgres.NewGroup(p),
		postgres.NewUser(p),
	)
	getUC := get_group.New(p.MakeQuery())
	updateUC := update_group.New(
		postgres.NewGroup(p),
	)
	deleteUC := delete_group.New(
		postgres.NewGroup(p),
	)
	getStudentsUC := get_students.New(p)
	addStudentsUC := add_students.New(
		postgres.NewGroup(p),
	)

	r.Route("/groups", func(r chi.Router) {
		r.Post("/", create_group.HTTP(createUC))
		r.Get("/", get_group.HTTP(getUC))
		r.Route("/"+param.GroupID.Path(), func(r chi.Router) {
			r.Put("/", update_group.HTTP(updateUC))
			r.Delete("/", delete_group.HTTP(deleteUC))
			r.Route("/students", func(r chi.Router) {
				r.Post("/", add_students.HTTP(addStudentsUC))
				r.Get("/", get_students.HTTP(getStudentsUC))
			})
		})
	})
}
