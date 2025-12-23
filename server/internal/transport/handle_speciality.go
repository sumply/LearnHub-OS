package transport

import (
	"net/http"
	"server/internal/usecase"
)

type specialCreateReq struct {
	Name string `json:"name"`
}

type specialityHandler struct {
	handler
	u usecase.Speciality
}

func newSpecialityHandler(u usecase.Speciality) *specialityHandler {
	return &specialityHandler{u: u}
}

func (h *specialityHandler) post(w http.ResponseWriter, r *http.Request) {
	var req specialCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

}
