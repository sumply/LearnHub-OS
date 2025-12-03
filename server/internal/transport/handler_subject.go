package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
)

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

	if err := h.usecase.Create(r.Context(), req.Name); err != nil {
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
	resp, err := h.usecase.Get(r.Context())
	if err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
