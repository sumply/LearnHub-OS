package handler

import (
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type Subject struct {
	handler
	usecase usecase.SubjectInterface
}

func NewSubject(s usecase.SubjectInterface) (*Subject, error) {
	if s == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &Subject{
		usecase: s,
	}, nil
}

func (h *Subject) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.SubjectCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if err := h.usecase.Create(r.Context(), identity, &req); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Subject) Get(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	data, err := h.usecase.Get(r.Context(), identity)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := make([]*dto.SubjectResp, len(data))
	for i, d := range data {
		resp[i] = dto.NewSubjectResp(d)
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *Subject) DeleteByID(w http.ResponseWriter, r *http.Request) {}
