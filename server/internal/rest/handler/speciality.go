package handler

import (
	"net/http"
	"server/internal/rest/dto"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type Speciality struct {
	handler
	u usecase.Speciality
}

func NewSpeciality(u usecase.Speciality) *Speciality {
	return &Speciality{u: u}
}

func (h *Speciality) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.SpecialityCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}
	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	err := h.u.Create(r.Context(), identity(auth), req.Name)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Speciality) Get(w http.ResponseWriter, r *http.Request) {
	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	data, err := h.u.Get(r.Context(), identity(auth))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	resp := make([]*dto.SpecialityResp, len(data))
	for i, d := range data {
		resp[i] = dto.NewSpecialityResp(d)
	}
	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
