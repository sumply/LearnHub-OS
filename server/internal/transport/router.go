package transport

import (
	"net/http"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	uu usecase.User,
	ug usecase.Group,
	us usecase.Subject,
	uq usecase.Quiz,
	uqr usecase.Answer,
	t TokenParser,
) (http.Handler, error) {
	r := chi.NewRouter()

	mw := middlewareBuilder{}

	r.Use(mw.buildLoggingRequest)
	r.Use(mw.buildLoggingResponse)

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
	qh, err := newQuizHandler(uq)
	if err != nil {
		return nil, err
	}
	qrh, err := newAnswerHandler(uqr)
	if err != nil {
		return nil, err
	}

	r.Post("/login", uh.login)

	r.Group(func(r chi.Router) {
		r.Use(mw.buildGetToken)
		r.Use(mw.buildValidateToken(t))

		addUserRouting(r, uh)
		addGroupRouting(r, gh)
		addSubjectsRouting(r, sh)
		addQuizRouting(r, qh)
		addAnswerRouting(r, qrh)
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

func addSubjectsRouting(r chi.Router, h *subjectHandler) {
	r.Post("/subjects", h.post)
	r.Get("/subjects", h.get)
}

func addQuizRouting(r chi.Router, h *quizHandler) {
	r.Post("/quizzes", h.post)
	r.Get("/quizzes", h.get)
	r.Get("/quizzes/{quiz_id}", h.getByID)
}

func addAnswerRouting(r chi.Router, h *answerHandler) {
	r.Post("/quizzes/{quiz_id}/answers", h.post)
	r.Get("/answers", h.get)
	r.Get("/answers/{answer_id}", h.getByID)
}
