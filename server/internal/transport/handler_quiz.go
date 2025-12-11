package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
	"time"
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

func (q *quizOptionsResp) fromUC(d usecase.QuizOptionsDomain) {
	q.ID = id(d.ID)
	q.Text = d.Text
	q.IsCorrect = d.IsCorrect
}

type quizQuestionResp struct {
	Title   string            `json:"title"`
	Options []quizOptionsResp `json:"options"`
}

func (q *quizQuestionResp) fromUC(d usecase.QuizQuestionDomain) {
	q.Title = d.Name
	opts := make([]quizOptionsResp, len(d.Answers))
	for i, o := range d.Answers {
		opts[i].fromUC(o)
	}
	q.Options = opts
}

type quizFullResp struct {
	ID        id                 `json:"id"`
	Title     string             `json:"title"`
	Summary   string             `json:"summary"`
	Subject   subjectResp        `json:"subject"`
	Owner     userShortResp      `json:"owner"`
	Questions []quizQuestionResp `json:"questions"`
}

func (q *quizFullResp) fromUC(domain usecase.QuizDomain) {
	q.ID = id(domain.ID)
	q.Title = domain.Name
	q.Summary = domain.Summary
	questions := make([]quizQuestionResp, len(domain.Quiestions))
	for i, d := range domain.Quiestions {
		questions[i].fromUC(d)
	}
	q.Questions = questions
}

type quizShortResp struct {
	ID      id            `json:"id"`
	Name    string        `json:"name"`
	Summary string        `json:"summary"`
	Subject subjectResp   `json:"subject"`
	Owner   userShortResp `json:"owner"`
}

func (q *quizShortResp) fromUC(d usecase.QuizDomain) {
	q.ID = id(d.ID)
	q.Name = d.Name
	q.Summary = d.Summary
}

type quizResultCreateReq struct {
	Answers []struct {
		QuestionID int64 `json:"question_id"`
		OptionID   int64 `json:"option_id"`
	} `json:"answers"`
}

type quizResultResp struct {
	ID            id            `json:"id"`
	TotalScore    int           `json:"total_score"`
	Score         int           `json:"score"`
	Completed     userShortResp `json:"completed"`
	CompletedTime time.Time     `json:"completed_time"`
	Group         groupResp     `json:"group"`
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
		sendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
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
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	resp := make([]quizShortResp, len(data))
	for i, d := range data {
		resp[i].fromUC(d)
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
		sendError(w, http.StatusInternalServerError, "")
		return
	}

	var resp quizFullResp
	resp.fromUC(data)
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *quizHandler) postByIDResult(w http.ResponseWriter, r *http.Request) {
	var req quizResultCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *quizHandler) getByIDResult(w http.ResponseWriter, r *http.Request) {
	var resp []quizResultResp
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *quizHandler) getByIDResultByID(w http.ResponseWriter, r *http.Request) {

}
