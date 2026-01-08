package rest

import (
	"net/http"
	"server/internal/rest/handler"
	"server/internal/rest/middleware"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	r.Delete("/users/{user_id}", h.DeleteByID)
}

func addGroupRouting(r chi.Router, h *handler.Group) {
	r.Post("/groups", h.Post)
	r.Get("/groups", h.Get)
	r.Get("/groups/{group_id}", h.GetByID)
	r.Post("/groups/{group_id}/students", h.PostStudents)
	r.Delete("/groups/{group_id}/students/{user_id}", h.DeleteStudentByID)
	r.Delete("/groups/{group_id}", h.DeleteByID)
}

func addSubjectsRouting(r chi.Router, h *handler.Subject) {
	r.Post("/subjects", h.Post)
	r.Get("/subjects", h.Get)
	r.Delete("/subjects/{subject_id}", h.DeleteByID)
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
