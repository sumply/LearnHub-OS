package handler

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/common"
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

func (h *handler) getParamID(r *http.Request, key string) (common.ID, error) {
	p, err := h.getParam(r, key)
	if err != nil {
		return 0, nil
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return common.ID(i), nil
}

func (h *handler) getParamQuizID(r *http.Request) (common.ID, error) {
	i, err := h.getParamID(r, "quiz_id")
	return i, err
}

func (h *handler) getParamAnswerID(r *http.Request) (common.ID, error) {
	i, err := h.getParamID(r, "answer_id")
	return i, err
}

func (h *handler) getParamUserID(r *http.Request) (common.ID, error) {
	i, err := h.getParamID(r, "user_id")
	return i, err
}

func (h *handler) getParamGroupID(r *http.Request) (common.ID, error) {
	return h.getParamID(r, "group_id")
}

func (h *handler) getParamProgressID(r *http.Request) (common.ID, error) {
	return h.getParamID(r, "progress_id")
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
	transport.SendError(w, status, err.Error())
	return status, msg
}
