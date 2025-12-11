package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type id uint64

type userShortResp struct {
	ID        id
	ShortName string `json:"short_name"`
}

func (u *userShortResp) fromUserData(d usecase.UserDomain) {
	u.ID = id(d.ID)
	u.ShortName = formatShortName(
		d.FirstName,
		d.LastName,
		d.MiddleName,
	)
}

type userFullResp struct {
	ID         id
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

func (u *userFullResp) fromUserData(d usecase.UserDomain) {
	u.ID = id(d.ID)
	u.FirstName = d.FirstName
	u.LastName = d.LastName
	u.MiddleName = d.MiddleName
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *loginRequest) toLoginParam() usecase.UserLoginParam {
	return usecase.UserLoginParam{
		Login:    r.Login,
		Password: usecase.Password(r.Password),
	}
}

type loginResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (r *loginResponse) fromJWT(d usecase.JWT) {
	r.AccessToken = d.AccessToken
	r.RefreshToken = d.RefreshToken
}

type userDTOPostRequest struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Role       string  `json:"role"`
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
	var req loginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}

	data, err := h.u.Login(r.Context(), req.toLoginParam())

	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	var resp loginResponse
	resp.fromJWT(data)

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
	}
}

func (h *userHandler) get(w http.ResponseWriter, r *http.Request) {
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	data, err := h.u.Get(r.Context(), auth.toIdentity())
	if err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	var resp []userShortResp
	for _, d := range data {
		var u userShortResp
		u.fromUserData(d)
		resp = append(resp, u)
	}

	if err := encodeJSON(w, resp); err != nil {
		sendEncodeError(w)
		return
	}
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

	data, err := h.u.GetMe(r.Context(), auth.toIdentity())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	var resp userFullResp
	resp.fromUserData(data)

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
		sendGetAuthDataError(w)
		return
	}

	param := usecase.UserCreateParam{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: *req.MiddleName,
	}

	err := h.u.Create(r.Context(), auth.toIdentity(), param)
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
		sendGetAuthDataError(w)
		return
	}

	err = h.u.Delete(r.Context(), auth.toIdentity(), usecase.ID(userID))
	if err != nil {
		sendError(w, http.StatusInternalServerError, "")
		return
	}
}
