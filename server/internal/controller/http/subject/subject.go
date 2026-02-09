package subject

import (
	"server/internal/adapter/postgres"
	"server/internal/feature/subject/create_subject"
	"server/internal/feature/subject/delete_subject"
	"server/internal/feature/subject/get_subject"
	"server/internal/feature/subject/update_subject"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, p *postgres.Postgres) {
	createUC := create_subject.New(p)
	updateUC := update_subject.New(p)
	deleteUC := delete_subject.New(p)
	getUC := get_subject.New(p)

	r.Post("/subjects", create_subject.HTTP(createUC))
	r.Get("/subjects", get_subject.HTTP(getUC))
	r.Put("/subjects/{subject_id}", update_subject.HTTP(updateUC))
	r.Delete("/subjects/{subject_id}", delete_subject.HTTP(deleteUC))
}
