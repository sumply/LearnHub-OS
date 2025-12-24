package handler

import (
	"fmt"
	"net/http"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type subjectResp struct {
	ID   id     `json:"id"`
	Name string `json:"name"`
}

type subjectCreateReq struct {
	Name string `json:"name"`
}

type Subject struct {
	handler
	usecase usecase.Subject
}

func NewSubject(s usecase.Subject) (*Subject, error) {
	if s == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &Subject{
		usecase: s,
	}, nil
}

func (h *Subject) Post(w http.ResponseWriter, r *http.Request) {
	var req subjectCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if err := h.usecase.Create(r.Context(), identity(auth), req.Name); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Subject) Get(w http.ResponseWriter, r *http.Request) {
	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	data, err := h.usecase.Get(r.Context(), identity(auth))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := make([]subjectResp, len(data))
	for i, d := range data {
		resp[i] = subjectResp{
			ID:   id(d.ID),
			Name: string(d.Name),
		}
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
