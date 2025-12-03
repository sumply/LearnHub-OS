package transport

import (
	"net/http"
	"server/internal/usecase"
)

type subjectDTOPostRequest struct {
	Name string `json:"name"`
}

type subjectshandler struct {
	usecase usecase.Subject
}

func (h *subjectshandler) Post(w http.ResponseWriter, r *http.Request) {
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

func (h *subjectshandler) Get(w http.ResponseWriter, r *http.Request) {
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
