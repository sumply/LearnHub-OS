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

func (r *subjectResp) fromUCSubjectDomain(d usecase.SubjectDomain) {
	r.ID = id(d.ID)
	r.Name = d.Name
}

type subjectDTOPostRequest struct {
	Name string `json:"name"`
}

type subjecthandler struct {
	usecase usecase.Subject
}

func newSubjectHandler(s usecase.Subject) (*subjecthandler, error) {
	if s == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &subjecthandler{
		usecase: s,
	}, nil
}

func (h *subjecthandler) post(w http.ResponseWriter, r *http.Request) {
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

func (h *subjecthandler) get(w http.ResponseWriter, r *http.Request) {
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
		resp[i].fromUCSubjectDomain(data[i])
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
