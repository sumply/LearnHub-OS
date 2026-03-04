package subject

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/subject/create_subject"
	"server/internal/feature/subject/delete_subject"
	"server/internal/feature/subject/get_subject"
	"server/internal/feature/subject/update_subject"
	"server/internal/pkg/http/param"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_subject.New(
		postgres.NewSubject(p),
	)
	updateUC := update_subject.New(
		postgres.NewSubject(p),
	)
	deleteUC := delete_subject.New(
		postgres.NewGroup(p),
	)
	getUC := get_subject.New(p)

	r.Route("/subjects", func(r chi.Router) {
		r.Post("/", create_subject.HTTP(createUC))
		r.Get("/", get_subject.HTTP(getUC))
		r.Route("/"+param.SubjectID.Path(), func(r chi.Router) {
			r.Put("/", update_subject.HTTP(updateUC))
			r.Delete("/", delete_subject.HTTP(deleteUC))
		})
	})
}
