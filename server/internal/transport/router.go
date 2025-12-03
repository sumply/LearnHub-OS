package transport

import (
	"net/http"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	uh := NewUserHandler(usecase.NewFakeUser())
	gh := groupsHandler{&usecase.FakeGroup{}}
	sh := subjectshandler{&usecase.FakeSubject{}}

	r.Post("/login", uh.Login)

	r.Group(func(r chi.Router) {
		r.Use(getTokenMiddleware)
		r.Use(validateTokenMiddleware(&FakeTokenParser{}))

		addUserRouting(r, uh)
		addGroupRouting(r, &gh)
		addSubjectsRouting(r, &sh)
	})

	return r
}

func addUserRouting(r chi.Router, h *UserHandler) {
	r.Post("/users", h.Post)
	r.Get("/users", h.Get)
	r.Get("/users/me", h.GetMe)
	r.Put("/users", h.Put)
	r.Delete("/users/{user_id}", h.Delete)
}

func addGroupRouting(r chi.Router, h *groupsHandler) {
	r.Post("/groups", h.Post)
	r.Get("/groups", h.Get)
}

func addSubjectsRouting(r chi.Router, h *subjectshandler) {
	r.Post("/subjects", h.Post)
	r.Get("/subjects", h.Get)
}
