package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
)

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
		sendDecodeError(w)
		return
	}

	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	err := h.usecase.Create(
		r.Context(),
		auth.toIdentity(),
		req.toUCParam(),
	)
	if err != nil {
		sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *quizHandler) get(w http.ResponseWriter, r *http.Request) {
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	data, err := h.usecase.Get(r.Context(), auth.toIdentity())
	if err != nil {
		sendUsecaseError(w, err)
		return
	}

	resp := make([]quizShortResp, len(data))
	for i, d := range data {
		resp[i] = quizShortRespFromDomain(d)
	}

	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *quizHandler) getByID(w http.ResponseWriter, r *http.Request) {
	quizID, err := getParamQuizID(r)
	if err != nil {
		sendParamError(w, err.Error())
		return
	}
	auth, ok := getAuthData(r.Context())
	if !ok {
		sendGetAuthDataError(w)
		return
	}

	data, err := h.usecase.GetByID(r.Context(), auth.toIdentity(), usecase.ID(quizID))
	if err != nil {
		sendUsecaseError(w, err)
		return
	}

	resp := quizFullRespFromDomain(data)
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
