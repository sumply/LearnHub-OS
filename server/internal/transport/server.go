package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(h HandlerInterface) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/authorization", h.Authz)
	r.Post("/registration", h.Registration)
	r.Get("/materials/card", h.GetMaterialCard)
	return r
}
