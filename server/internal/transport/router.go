package transport

import (
	"net/http"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	uu usecase.User,
	ug usecase.Group,
	us usecase.Subject,
	t TokenParser,
) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	uh, err := newUserHandler(uu)
	if err != nil {
		return nil, err
	}
	gh, err := newGroupHandler(ug)
	if err != nil {
		return nil, err
	}
	sh, err := newSubjectHandler(us)
	if err != nil {
		return nil, err
	}
	qh := newQuizHandler()

	r.Post("/login", uh.login)

	r.Group(func(r chi.Router) {
		r.Use(getTokenMiddleware)
		r.Use(validateTokenMiddleware(t))

		addUserRouting(r, uh)
		addGroupRouting(r, gh)
		addSubjectsRouting(r, sh)
		addQuizRouting(r, qh)
	})

	return r, nil
}

func addUserRouting(r chi.Router, h *userHandler) {
	r.Post("/users", h.post)
	r.Get("/users", h.get)
	r.Get("/users/me", h.getMe)
	r.Put("/users", h.put)
	r.Delete("/users/{user_id}", h.delete)
}

func addGroupRouting(r chi.Router, h *groupHandler) {
	r.Post("/groups", h.post)
	r.Get("/groups", h.get)
}

func addSubjectsRouting(r chi.Router, h *subjecthandler) {
	r.Post("/subjects", h.post)
	r.Get("/subjects", h.get)
}

func addQuizRouting(r chi.Router, h *quizHandler) {
	r.Post("/quizzes", h.post)
	r.Get("/quizzes", h.get)
	r.Get("/quizzes/{quiz_id}", h.getByID)
	r.Post("/quizzes/{quiz_id}/result", h.postByIDResult)
	r.Get("/quizzes/{quiz_id}/result", h.getByIDResult)
}
