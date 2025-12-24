package handler

import (
	"net/http"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type specialCreateReq struct {
	Name string `json:"name"`
}

type Speciality struct {
	handler
	u usecase.Speciality
}

func NewSpeciality(u usecase.Speciality) *Speciality {
	return &Speciality{u: u}
}

func (h *Speciality) Post(w http.ResponseWriter, r *http.Request) {
	var req specialCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

}
