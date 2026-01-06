package rest

import (
	"net/http"
	"server/internal/rest/handler"
	"server/internal/rest/middleware"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	uu usecase.UserInterface,
	ug usecase.GroupInterface,
	us usecase.SubjectInterface,
	uq usecase.QuizInterface,
	up usecase.ProgressInterface,
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
	quizHandler, err := handler.NewQuiz(uq)
	if err != nil {
		return nil, err
	}
	progressHandler, err := handler.NewProgress(up)
	if err != nil {
		return nil, err
	}

	r.Post("/login", uh.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.GetToken)
		r.Use(middleware.ValidateToken(p))

		addUserRouting(r, uh)
		addGroupRouting(r, gh)
		addSubjectsRouting(r, sh)
		addQuizRouting(r, quizHandler)
		addProgressRouting(r, progressHandler)
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
	r.Post("/groups/{group_id}/students", h.PostStudents)
}

func addSubjectsRouting(r chi.Router, h *handler.Subject) {
	r.Post("/subjects", h.Post)
	r.Get("/subjects", h.Get)
}

func addQuizRouting(r chi.Router, h *handler.Quiz) {
	r.Post("/quizzes", h.Post)
	r.Get("/quizzes", h.Get)
	r.Delete("/quizzes/{quiz_id}", h.Delete)
}

func addProgressRouting(r chi.Router, h *handler.Progress) {
	r.Get("/progress", h.Get)
	r.Post("/progress/{progress_id}/start", h.PostStart)
	r.Post("/progress/{progress_id}/finish", h.PostFinish)
	r.Patch("/progress/{progress_id}/answer/{answer_id}", h.PatchAnswer)
	r.Post("/progress/{progress_id}/answer/{answer_id}/correct", h.PostAnswerCorrect)
	r.Post("/progress/{progress_id}/answer/{answer_id}/incorrect", h.PostAnswerIncorrect)
}
