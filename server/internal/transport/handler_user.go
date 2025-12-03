package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userDTOLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userDTOLoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type userDTOPostRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Role       string `json:"role"`
}

type userDTOPutRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type userHandler struct {
	u usecase.User
}

func newUserHandler(u usecase.User) (*userHandler, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &userHandler{
		u: u,
	}, nil
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var req userDTOLoginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}

	re, err := h.u.Login(r.Context(), usecase.UserLoginParam{
		Email:    req.Login,
		Password: req.Password,
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	resp := userDTOLoginResponse{
		RefreshToken: re.RefreshToken,
		AccessToken:  re.AccessToken,
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
	}
}

func (h *userHandler) get(w http.ResponseWriter, r *http.Request) {
}

func (h *userHandler) getMe(w http.ResponseWriter, r *http.Request) {
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendError(
			w,
			http.StatusInternalServerError,
			"Не удалось получить данные токена авторизации.",
		)
		return
	}

	uauth := usecase.AuthData{
		Subject: auth.subject,
		Role:    auth.role,
	}

	resp, err := h.u.GetMe(r.Context(), uauth)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *userHandler) put(w http.ResponseWriter, r *http.Request) {
	var req userDTOPutRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		return
	}
}

func (h *userHandler) post(w http.ResponseWriter, r *http.Request) {
	var req userDTOPostRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendError(
			w,
			http.StatusInternalServerError,
			"Не удалось получить данные токена авторизации.",
		)
		return
	}

	uath := usecase.AuthData{
		Subject: auth.subject,
		Role:    auth.role,
	}
	param := usecase.UserCreateParam{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: req.MiddleName,
	}

	err := h.u.Create(r.Context(), uath, param)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) delete(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		sendError(
			w,
			http.StatusBadRequest,
			"Не удалось преобразовать параметр {user_id} в число.",
		)
	}

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendError(
			w,
			http.StatusInternalServerError,
			"Не удалось получить данные токена авторизации.",
		)
		return
	}

	uath := usecase.AuthData{Subject: auth.subject, Role: auth.role}
	err = h.u.Delete(r.Context(), uath, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}
}
