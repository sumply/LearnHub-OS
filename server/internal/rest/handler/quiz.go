package handler

import (
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type Quiz struct {
	handler
	u usecase.QuizInterface
}

func NewQuiz(u usecase.QuizInterface) (*Quiz, error) {
	if u == nil {
		return nil, fmt.Errorf("usecase is nil")
	}
	return &Quiz{
		u: u,
	}, nil
}

func (h *Quiz) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.QuizCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if err := h.u.Create(r.Context(), identity, &req); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Quiz) Get(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	data, err := h.u.Get(r.Context(), identity)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	resp := dto.NewSliceQuizShortResp(data)
	if err := transport.EncodeJSON(w, resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *Quiz) Delete(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	quizID, err := h.getParamQuizID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	err = h.u.Delete(r.Context(), identity, quizID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
