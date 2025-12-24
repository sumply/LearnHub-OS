package handler

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/rest/transport"
	"server/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct{}

func (h *handler) sendDecodeError(w http.ResponseWriter) (int, string) {
	const errMsg = "Не удалось распарсить тело запроса."
	transport.SendError(
		w,
		http.StatusBadRequest,
		errMsg,
	)
	return http.StatusBadRequest, errMsg
}

func (h *handler) sendEncodeError(w http.ResponseWriter) (int, string) {
	const errMsg = "Ошибка при маршалинге ответа."
	transport.SendError(
		w,
		http.StatusInternalServerError,
		errMsg,
	)
	return http.StatusInternalServerError, errMsg
}

func (h *handler) getParam(r *http.Request, key string) (string, error) {
	p := chi.URLParam(r, key)
	if p == "" {
		return "", fmt.Errorf("параметр {%s} пустой", key)
	}
	return p, nil
}

func (h *handler) getParamInt(r *http.Request, key string) (int, error) {
	p, err := h.getParam(r, key)
	if err != nil {
		return 0, nil
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (h *handler) getParamQuizID(r *http.Request) (id, error) {
	i, err := h.getParamInt(r, "quiz_id")
	return id(i), err
}

func (h *handler) getParamAnswerID(r *http.Request) (id, error) {
	i, err := h.getParamInt(r, "answer_id")
	return id(i), err
}

func (h *handler) getParamUserID(r *http.Request) (id, error) {
	i, err := h.getParamInt(r, "user_id")
	return id(i), err
}

func (h *handler) sendParamError(w http.ResponseWriter, what string) {
	transport.SendError(
		w,
		http.StatusBadRequest,
		what,
	)
}

func (h *handler) sendUsecaseError(w http.ResponseWriter, err error) (int, string) {
	var status int
	var msg string
	switch {
	case errors.Is(err, usecase.ErrAccess):
		status = http.StatusForbidden
		msg = "forbidden"
	case errors.Is(err, usecase.ErrCollision):
		status = http.StatusConflict
		msg = "conflict"
	case errors.Is(err, usecase.ErrNotFound):
		status = http.StatusNotFound
		msg = "Not found"
	default:
		status = http.StatusInternalServerError
		msg = "Internal server error"
	}
	transport.SendError(w, status, msg)
	return status, msg
}

var roleMap = map[string]domain.UserRole{
	"root":    domain.UserRoot,
	"admin":   domain.UserAdmin,
	"teacher": domain.UserTeacher,
	"student": domain.UserStudent,
}

func userRole(s string) domain.UserRole {
	r, ok := roleMap[s]
	if !ok {
		r = domain.UserInvalid
	}
	return r
}

func identity(a transport.AuthData) usecase.Identity {
	role, ok := roleMap[a.Role]
	if !ok {
		role = domain.UserInvalid
	}
	return usecase.Identity{
		ID:   domain.UserID(a.ID),
		Role: role,
	}
}
