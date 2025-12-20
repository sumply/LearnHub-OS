package transport

import (
	"fmt"
	"net/http"
	"server/internal/logger"
	"server/internal/usecase"
)

type id uint64

type userShortResp struct {
	ID        id     `json:"id"`
	ShortName string `json:"short_name"`
}

func userShortRespFromDomain(d usecase.UserDomain) userShortResp {
	return userShortResp{
		ID: id(d.ID),
		ShortName: formatShortName(
			d.FirstName,
			d.LastName,
			d.MiddleName,
		),
	}
}

type userFullResp struct {
	ID         id     `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

func userFullRespFromDomain(d usecase.UserDomain) userFullResp {
	return userFullResp{
		ID:         id(d.ID),
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		MiddleName: d.MiddleName,
	}
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
	log := logger.FromCtx(r.Context())

	var req loginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		status, msg := sendDecodeError(w)
		loggerWithResponse(log, status, msg).
			Warn("failed decoded request")
		return
	}

	data, err := h.u.Login(r.Context(), req.toLoginParam())
	if err != nil {
		status, msg := sendUsecaseError(w, err)
		loggerWithResponse(log, status, msg).
			Warn("failed authorization")
		return
	}

	var resp loginResponse
	resp.fromJWT(data)

	if err := encodeJSON(w, &resp); err != nil {
		status, msg := sendEncodeError(w)
		loggerWithResponse(log, status, msg).
			Error("failed encoding response")
	}
}

func (h *userHandler) get(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method get")

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	log.Debug("Calling a usecase method Get.")
	data, err := h.u.Get(r.Context(), auth.toIdentity())
	if err != nil {
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	log.Debug("Mapping a usecase domain to a response struct")
	var resp []userShortResp
	for _, d := range data {
		resp = append(resp, userShortRespFromDomain(d))
	}

	if err := encodeJSON(w, resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *userHandler) getMe(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method getMe")

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	log.Debug("Calling a usecase method GetMe")
	data, err := h.u.GetMe(r.Context(), auth.toIdentity())
	if err != nil {
		sendUsecaseError(w, err)
		return
	}

	log.Debug("Mapping from a usecase domain to a response struct")
	resp := userFullRespFromDomain(data)

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}

	log.Debug("Ending a handler method getMe")
}

func (h *userHandler) put(w http.ResponseWriter, r *http.Request) {
	var req userDTOPutRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		return
	}
}

func (h *userHandler) post(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method post()")

	log.Debug("Decoding a request")
	var req userDTOPostRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}

	log.Debug("Calling the function getAuthData()")
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

	log.Debug("Calling a usecase method create()")
	err := h.u.Create(r.Context(), auth.toIdentity(), param)
	if err != nil {
		sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) delete(w http.ResponseWriter, r *http.Request) {
	log := logger.FromCtx(r.Context())
	log.Debug("Called a handler method delete()")

	log.Debug("Calling the function getParamUserID()")
	userID, err := getParamUserID(r)
	if err != nil {
		sendParamError(w, err.Error())
		return
	}

	log.Debug("Calling the function getAuthData()")
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	log.Debug("Calling a usecase method Delete()")
	err = h.u.Delete(r.Context(), auth.toIdentity(), usecase.ID(userID))
	if err != nil {
		sendUsecaseError(w, err)
		return
	}
}
