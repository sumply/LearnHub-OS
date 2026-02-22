package param

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IDParam string

const (
	AttemptID IDParam = "attempt_id"
	UserID    IDParam = "user_id"
	QuizID    IDParam = "quiz_id"
	GroupID   IDParam = "group_id"
	SubjectID IDParam = "subject_id"
)

func (param IDParam) Path() string {
	return "{" + string(param) + "}"
}

func ID(r *http.Request, key IDParam) (uuid.UUID, error) {
	param := chi.URLParam(r, string(key))
	return uuid.Parse(param)
}
