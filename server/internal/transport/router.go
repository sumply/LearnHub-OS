package transport

import (
	"net/http"
	"server/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	uh := NewUserHandler(usecase.NewFakeUser())

	r.Post("/login", uh.Login)

	r.Group(func(r chi.Router) {
		r.Use(getTokenFromHeader)
		r.Use(makeValidateTokenFunc(&FakeTokenParser{}))

		addUserRouting(r, uh)
	})

	return r
}

func addUserRouting(r chi.Router, h *UserHandler) {
	r.Post("/users", h.Post)
	r.Get("/users", h.Get)
	r.Get("/users/me", h.GetMe)
	r.Put("/users", h.Put)
	r.Delete("/users/{user_id}", h.Delete)
}

/*
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
*/
