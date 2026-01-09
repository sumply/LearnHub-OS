package handler

import (
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/logger"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type User struct {
	handler
	u usecase.UserInterface
}

func NewUser(u usecase.UserInterface) (*User, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &User{
		u: u,
	}, nil
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	data, err := h.u.Login(r.Context(), &req)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := dto.NewLoginResp(data)

	if err := transport.EncodeJSON(w, resp); err != nil {
		h.sendEncodeError(w)
	}
}

func (h *User) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method get")

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	log.Debug("Calling a usecase method Get.")
	data, err := h.u.Get(r.Context(), identity)
	if err != nil {
		transport.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	log.Debug("Mapping a usecase domain to a response struct")
	var resp []*dto.UserShortResp
	for _, d := range data {
		user := dto.NewUserShortResp(d)
		resp = append(resp, user)
	}

	if err := transport.EncodeJSON(w, resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *User) GetMe(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	data, err := h.u.GetMe(r.Context(), identity)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := dto.NewUserFullResp(data, nil)

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

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	resp, err := h.u.GetByID(r.Context(), identity, userID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *User) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if !ok {
		h.sendDecodeError(w)
		return
	}

	err := h.u.Create(r.Context(), identity, &req)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	userID, err := h.getParamUserID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	err = h.u.Delete(r.Context(), identity, userID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
