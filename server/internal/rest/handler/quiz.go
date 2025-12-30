package handler

import (
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/rest/transport"
	"server/internal/usecase"
)

type Quiz struct {
	handler
	u *usecase.Quiz
}

func NewQuiz(u *usecase.Quiz) (*Quiz, error) {
	if u == nil {
		return nil, fmt.Errorf("usecase is nil")
	}
	return &Quiz{
		u: u,
	}, nil
}

func (h *Quiz) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.QuizCreateReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}

	if err := h.u.Create(r.Context(), identity, &req); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*
type quizOptionsCreate struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

func (q *quizOptionsCreate) toUCParam() usecase.QuizCreateOption {
	return usecase.QuizCreateOption{
		Text:      q.Text,
		IsCorrect: q.IsCorrect,
	}
}

type quizQuestionCreate struct {
	Title   string `json:"title"`
	Options []quizOptionsCreate
}

func (q *quizQuestionCreate) toUCParam() usecase.QuizCreateQuestion {
	uopts := make([]usecase.QuizCreateOption, len(q.Options))
	for i, opt := range q.Options {
		uopts[i] = opt.toUCParam()
	}
	return usecase.QuizCreateQuestion{
		Title:   q.Title,
		Options: uopts,
	}
}

type quizCreateReq struct {
	Name      string               `json:"name"`
	Summary   string               `json:"summary"`
	SubjectID int64                `json:"subject_id"`
	Questions []quizQuestionCreate `json:"questions"`
}

func (r *quizCreateReq) toUCParam() usecase.QuizCreateParam {
	uq := make([]usecase.QuizCreateQuestion, len(r.Questions))
	for i, q := range r.Questions {
		uq[i] = q.toUCParam()
	}
	p := usecase.QuizCreateParam{
		Name:      r.Name,
		Summary:   r.Summary,
		SubjectID: usecase.ID(r.SubjectID),
		Questions: uq,
	}
	return p
}

type quizOptionsResp struct {
	ID        id     `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

func quizOptionsRespFromDomain(d usecase.QuizOptionsDomain) quizOptionsResp {
	return quizOptionsResp{
		ID:        id(d.ID),
		Text:      d.Text,
		IsCorrect: d.IsCorrect,
	}
}

type quizQuestionResp struct {
	Title   string            `json:"title"`
	Options []quizOptionsResp `json:"options"`
}

func quizQuestionRespFromDomain(d usecase.QuizQuestionDomain) quizQuestionResp {
	opts := make([]quizOptionsResp, len(d.Answers))
	for i, d := range d.Answers {
		opts[i] = quizOptionsRespFromDomain(d)
	}
	return quizQuestionResp{
		Title:   d.Name,
		Options: opts,
	}
}

type quizFullResp struct {
	ID        id                 `json:"id"`
	Title     string             `json:"title"`
	Summary   string             `json:"summary"`
	Subject   subjectResp        `json:"subject"`
	Owner     userShortResp      `json:"owner"`
	Questions []quizQuestionResp `json:"questions"`
}

func quizFullRespFromDomain(domain usecase.QuizDomain) quizFullResp {
	questions := make([]quizQuestionResp, len(domain.Questions))
	for i, d := range domain.Questions {
		questions[i] = quizQuestionRespFromDomain(d)
	}
	return quizFullResp{
		ID:      id(domain.ID),
		Title:   domain.Name,
		Summary: domain.Summary,
	}
}

type quizShortResp struct {
	ID      id            `json:"id"`
	Name    string        `json:"name"`
	Summary string        `json:"summary"`
	Subject subjectResp   `json:"subject"`
	Owner   userShortResp `json:"owner"`
}

func quizShortRespFromDomain(d usecase.QuizDomain) quizShortResp {
	return quizShortResp{
		ID:      id(d.ID),
		Name:    d.Name,
		Summary: d.Summary,
	}
}

type quizHandler struct {
	handler
	usecase usecase.Quiz
}

func newQuizHandler(u usecase.Quiz) (*quizHandler, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &quizHandler{
		usecase: u,
	}, nil
}

func (h *quizHandler) post(w http.ResponseWriter, r *http.Request) {
	var req quizCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}

	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}

	err := h.usecase.Create(
		r.Context(),
		auth.toIdentity(),
		req.toUCParam(),
	)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *quizHandler) get(w http.ResponseWriter, r *http.Request) {
	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}

	data, err := h.usecase.Get(r.Context(), auth.toIdentity())
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := make([]quizShortResp, len(data))
	for i, d := range data {
		resp[i] = quizShortRespFromDomain(d)
	}

	if err := encodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *quizHandler) getByID(w http.ResponseWriter, r *http.Request) {
	quizID, err := h.getParamQuizID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}

	data, err := h.usecase.GetByID(r.Context(), auth.toIdentity(), usecase.ID(quizID))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := quizFullRespFromDomain(data)
	if err := encodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
*/
