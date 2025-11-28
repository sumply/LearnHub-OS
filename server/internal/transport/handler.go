package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"server/internal/usecase"
)

type userDTOLoginRequest struct {
	Email    string `json:"email"`
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
}

type userDTOPutRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type UserHandler struct {
	u usecase.User
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userDTOLoginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		return
	}

	ctx := context.Background()

	re, err := h.u.Login(ctx, usecase.UserLoginParam{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return
	}

	resp := userDTOLoginResponse{
		RefreshToken: re.RefreshToken,
		AccessToken:  re.AccessToken,
	}

	if err := encodeJSON(w, &resp); err != nil {
		return
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	auth, ok := getAuthData(r.Context())
	if !ok {
		return
	}

	err := h.u.Create(
		r.Context(),
		usecase.AuthData{
			Subject: auth.subject,
			Role:    auth.role,
		},
		usecase.UserCreateParam{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		},
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Delete(http.ResponseWriter, *http.Request) {

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
	auth, ok := ctx.Value(auth).(authData)
	return auth, ok
}
