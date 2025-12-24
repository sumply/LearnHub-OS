package handler

import (
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

func formatShortName(f string, l string, m string) string {
	if f == "" || l == "" {
		return ""
	}
	name := fmt.Sprintf("%s %v.", f, l[0])
	if m != "" {
		name = fmt.Sprintf("%s %v.", name, m[0])
	}
	return name
}

type id uint64

type userShortResp struct {
	ID        id     `json:"id"`
	ShortName string `json:"short_name"`
}

type userFullResp struct {
	ID         id     `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

func newUserFullResp(d *domain.User) userFullResp {
	return userFullResp{
		ID:         id(d.ID),
		FirstName:  string(d.FirstName),
		LastName:   string(d.LastName),
		MiddleName: string(d.MiddleName),
	}
}

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type userCreateReq struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
}

type userDTOPutRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type User struct {
	handler
	u usecase.User
}

func NewUser(u usecase.User) (*User, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &User{
		u: u,
	}, nil
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	param := usecase.UserLoginParam{
		Login:    req.Login,
		Password: req.Password,
	}

	data, err := h.u.Login(r.Context(), param)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := loginResp{
		RefreshToken: data.Refresh,
		AccessToken:  data.Access,
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
	}
}

func (h *User) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method get")

	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	log.Debug("Calling a usecase method Get.")
	data, err := h.u.Get(r.Context(), identity(auth))
	if err != nil {
		transport.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	log.Debug("Mapping a usecase domain to a response struct")
	var resp []userShortResp
	for _, d := range data {
		user := userShortResp{
			ID: id(d.ID),
			ShortName: formatShortName(
				string(d.FirstName),
				string(d.LastName),
				string(d.MiddleName),
			),
		}
		resp = append(resp, user)
	}

	if err := transport.EncodeJSON(w, resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *User) GetMe(w http.ResponseWriter, r *http.Request) {
	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	data, err := h.u.GetMe(r.Context(), identity(auth))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := newUserFullResp(data)

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *User) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getParamUserID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	data, err := h.u.GetByID(r.Context(), identity(auth), domain.UserID(userID))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := newUserFullResp(data)
	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *User) Post(w http.ResponseWriter, r *http.Request) {
	var req userCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)

		return
	}

	auth, ok := transport.NewAuthDataFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if !ok {
		h.sendDecodeError(w)
		return
	}

	param := usecase.UserCreateParam{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: *req.MiddleName,
		Email:      req.Email,
		Role:       userRole(auth.Role),
	}

	err := h.u.Create(r.Context(), identity(auth), param)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
