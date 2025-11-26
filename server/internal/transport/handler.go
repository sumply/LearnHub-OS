package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"server/internal/dto"
	"server/internal/service"
	"time"
)

type HandlerInterface interface {
	Authz(http.ResponseWriter, *http.Request)
	Registration(http.ResponseWriter, *http.Request)
	GetMaterialCard(http.ResponseWriter, *http.Request)
}

type Handler struct {
	service service.Interface
}

func NewHandler(s service.Interface) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Authz(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthzRequest
	if err := decode(r.Body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, err := h.service.Authz(&req)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	encode(w, &resp)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Registration(http.ResponseWriter, *http.Request) {
}

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) Authz(w http.ResponseWriter, r *http.Request) {
	token := map[string]any{
		"jwt_refresh": 10,
		"jwt_access":  10,
	}
	user := map[string]any{
		"user_id":     10,
		"first_name":  "First",
		"last_name":   "Last",
		"middle_name": "Middle",
		"icon_ref":    "ref",
		"created_at":  time.Now().UTC(),
	}
	resp := map[string]any{
		"jwt":  token,
		"user": user,
	}
	json.NewEncoder(w).Encode(&resp)
}

func (h *MockHandler) Registration(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *MockHandler) GetMaterialCard(w http.ResponseWriter, r *http.Request) {
	m := dto.MaterialCardResponse{
		ID:        10,
		Name:      "Mock",
		Type:      "pdf",
		Size:      1024,
		Summary:   "It's the mock material",
		Subject:   "Something",
		Class:     "11A",
		CreatedAt: time.Now(),
		Tags:      []string{"mock"},
	}
	encode(w, &m)
}

func decode(r io.ReadCloser, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

func encode(w io.Writer, v any) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
