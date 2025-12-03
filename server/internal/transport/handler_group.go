package transport

import (
	"net/http"
	"server/internal/usecase"
)

type groupDTOPostRequest struct {
	Name string `json:"name"`
}

type groupsHandler struct {
	usecase usecase.Group
}

func (h *groupsHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req groupDTOPostRequest
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

func (h *groupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.Get(r.Context())
	if err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
