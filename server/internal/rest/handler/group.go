package handler

import (
	"fmt"
	"net/http"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type groupResp struct {
	ID   id     `json:"id"`
	Name string `json:"name"`
}

type groupCreateReq struct {
	Name string `json:"name"`
}

type Group struct {
	handler
	usecase usecase.Group
}

func NewGroup(g usecase.Group) (*Group, error) {
	if g == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &Group{
		usecase: g,
	}, nil
}

func (h *Group) Post(w http.ResponseWriter, r *http.Request) {
	var req groupCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	if err := h.usecase.Create(r.Context(), req.Name); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Group) Get(w http.ResponseWriter, r *http.Request) {
	data, err := h.usecase.Get(r.Context())
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := make([]groupResp, len(data))
	for i, d := range data {
		resp[i] = groupResp{
			ID:   id(d.ID),
			Name: string(d.Name),
		}
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
