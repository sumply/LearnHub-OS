package handler

import (
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/rest/dto"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

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
	var req dto.GroupCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	param := usecase.GroupCreateParam{
		Name:         req.Name,
		CuratorID:    domain.UserID(req.CuratorID),
		SpecialityID: domain.SpecialityID(req.SpecialityID),
	}
	if err := h.usecase.Create(r.Context(), identity(auth), param); err != nil {
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

	resp := make([]*dto.GroupResp, len(data))
	for i, d := range data {
		resp[i] = dto.NewGroupResp(d)
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
