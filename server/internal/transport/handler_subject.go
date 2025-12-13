package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
)

type subjectResp struct {
	ID   id     `json:"id"`
	Name string `json:"name"`
}

func (r *subjectResp) fromDomain(d usecase.SubjectDomain) {
	r.ID = id(d.ID)
	r.Name = d.Name
}

type subjectDTOPostRequest struct {
	Name string `json:"name"`
}

type subjectHandler struct {
	usecase usecase.Subject
}

func newSubjectHandler(s usecase.Subject) (*subjectHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &subjectHandler{
		usecase: s,
	}, nil
}

func (h *subjectHandler) post(w http.ResponseWriter, r *http.Request) {
	var req subjectDTOPostRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	if err := h.usecase.Create(r.Context(), auth.toIdentity(), req.Name); err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *subjectHandler) get(w http.ResponseWriter, r *http.Request) {
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	data, err := h.usecase.Get(r.Context(), auth.toIdentity())
	if err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	resp := make([]subjectResp, len(data))
	for i := range data {
		resp[i].fromDomain(data[i])
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
