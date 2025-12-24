package rest

import (
	"net/http"
	"server/internal/domain"
	"server/internal/rest/handler"
	"server/internal/rest/middleware"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	uu usecase.User,
	ug usecase.Group,
	us usecase.Subject,
	usp usecase.Speciality,
	g domain.TokenGenerator,
	p middleware.TokenParser,
) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.LogRequest)
	r.Use(middleware.LogResponse)

	uh, err := handler.NewUser(uu)
	if err != nil {
		return nil, err
	}
	gh, err := handler.NewGroup(ug)
	if err != nil {
		return nil, err
	}
	sh, err := handler.NewSubject(us)
	if err != nil {
		return nil, err
	}
	sph := handler.NewSpeciality(usp)

	r.Post("/login", uh.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.GetToken)
		r.Use(middleware.ValidateToken(p))

		addUserRouting(r, uh)
		addGroupRouting(r, gh)
		addSubjectsRouting(r, sh)
		addSpecialityRouting(r, sph)
	})

	return r, nil
}

func addUserRouting(r chi.Router, h *handler.User) {
	r.Post("/users", h.Post)
	r.Get("/users", h.Get)
	r.Get("/users/me", h.GetMe)
	r.Get("/users/{user_id}", h.GetByID)
}

func addGroupRouting(r chi.Router, h *handler.Group) {
	r.Post("/groups", h.Post)
	r.Get("/groups", h.Get)
}

func addSubjectsRouting(r chi.Router, h *handler.Subject) {
	//r.Post("/subjects", h.post)
	//r.Get("/subjects", h.Get)
}

func addSpecialityRouting(r chi.Router, h *handler.Speciality) {
	//r.Post("/speciality", h.Post)
}

/*
func addQuizRouting(r chi.Router, h *handler.) {
	r.Post("/quizzes", h.post)
	r.Get("/quizzes", h.get)
	r.Get("/quizzes/{quiz_id}", h.getByID)
}

func addAnswerRouting(r chi.Router, h *answerHandler) {
	r.Post("/quizzes/{quiz_id}/answers", h.post)
	r.Get("/answers", h.get)
	r.Get("/answers/{answer_id}", h.getByID)
}
*/
