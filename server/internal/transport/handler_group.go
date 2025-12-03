package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
)

type groupDTOPostRequest struct {
	Name string `json:"name"`
}

type groupHandler struct {
	usecase usecase.Group
}

func newGroupHandler(g usecase.Group) (*groupHandler, error) {
	if g == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &groupHandler{
		usecase: g,
	}, nil
}

func (h *groupHandler) post(w http.ResponseWriter, r *http.Request) {
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

func (h *groupHandler) get(w http.ResponseWriter, r *http.Request) {
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
