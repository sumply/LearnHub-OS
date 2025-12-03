package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"server/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	u usecase.User
}

func NewUserHandler(u usecase.User) *UserHandler {
	return &UserHandler{
		u: u,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userDTOLoginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendError(
			w,
			http.StatusBadRequest,
			"Не удалось распарсить тело запроса.",
		)
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
		sendError(
			w,
			http.StatusInternalServerError,
			"Ошибка при маршалинге ответа.",
		)
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
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
		sendError(
			w,
			http.StatusInternalServerError,
			"Ошибка при маршалинге ответа.",
		)
	}
}

func (h *UserHandler) Put(w http.ResponseWriter, r *http.Request) {
	var req userDTOPutRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		return
	}
}

func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req userDTOPostRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendError(
			w,
			http.StatusBadRequest,
			"Не удалось распарсить тело запроса.",
		)
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

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func decodeJSON(r io.ReadCloser, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

func encodeJSON(w io.Writer, v any) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func getAuthData(ctx context.Context) (authData, bool) {
	auth, ok := ctx.Value(authKey).(authData)
	return auth, ok
}
