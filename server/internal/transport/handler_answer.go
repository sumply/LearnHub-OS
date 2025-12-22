package transport

import (
	"fmt"
	"net/http"
	"server/internal/usecase"
	"time"
)

type answerReq struct {
	QuestionID id `json:"question_id"`
	OptionID   id `json:"option_id"`
}

func selectedOptionFromAnswerReq(r answerReq) usecase.SelectedOption {
	return usecase.SelectedOption{
		QuestionID: usecase.ID(r.QuestionID),
		OptionID:   usecase.ID(r.OptionID),
	}
}

type answerCreateReq struct {
	Answers []answerReq `json:"answers"`
}

type answerShortResp struct {
	ID            id            `json:"id"`
	Quiz          quizShortResp `json:"quiz"`
	TotalScore    int           `json:"total_score"`
	Score         int           `json:"score"`
	Completed     bool          `json:"completed"`
	User          userShortResp `json:"user"`
	CompletedTime time.Time     `json:"completed_time"`
	Group         groupResp     `json:"group"`
}

func answerShortRespFromDomain(d usecase.AnswerDomain) answerShortResp {
	return answerShortResp{
		ID:         id(d.ID),
		TotalScore: d.TotalScore,
		Score:      d.Score,
		Completed:  d.Completed,
		User:       userShortRespFromDomain(d.User),
		Quiz:       quizShortRespFromDomain(d.Quiz),
	}
}

type answerFullResp struct {
	ID            id                 `json:"id"`
	TotalScore    int                `json:"total_score"`
	Score         int                `json:"score"`
	Completed     bool               `json:"completed"`
	CompletedTime time.Time          `json:"completed_time"`
	User          userShortResp      `json:"user"`
	Quiz          quizShortResp      `json:"quiz"`
	Answers       []answeredQuestion `json:"answers"`
}

func answerFullRespFromDomain(d usecase.AnswerDomain) answerFullResp {
	answers := make([]answeredQuestion, len(d.Answers))
	for i, a := range d.Answers {
		answers[i] = answerQuestionFromDomain(a)
	}
	return answerFullResp{
		ID:         id(d.ID),
		TotalScore: d.TotalScore,
		Score:      d.Score,
		Completed:  d.Completed,
		User:       userShortRespFromDomain(d.User),
		Quiz:       quizShortRespFromDomain(d.Quiz),
		Answers:    answers,
	}
}

type answeredQuestion struct {
	Question quizQuestionResp `json:"question"`
	Answered quizOptionsResp  `json:"answered"`
}

func answerQuestionFromDomain(d usecase.AnsweredQuestion) answeredQuestion {
	return answeredQuestion{
		Question: quizQuestionRespFromDomain(d.Question),
		Answered: quizOptionsRespFromDomain(d.Answered),
	}
}

type answerHandler struct {
	handler
	usecase usecase.Answer
}

func answerCreateParamFromRequest(quizID id, req []answerReq) usecase.AnswerCreateParam {
	answers := make([]usecase.SelectedOption, len(req))
	for i, r := range req {
		answers[i] = selectedOptionFromAnswerReq(r)
	}
	return usecase.AnswerCreateParam{
		QuizID:  usecase.ID(quizID),
		Answers: answers,
	}
}

func newAnswerHandler(u usecase.Answer) (*answerHandler, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &answerHandler{
		usecase: u,
	}, nil
}

func (h *answerHandler) post(w http.ResponseWriter, r *http.Request) {
	quizID, err := h.getParamQuizID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	var req answerCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}
	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}

	err = h.usecase.Create(r.Context(), auth.toIdentity(), answerCreateParamFromRequest(quizID, req.Answers))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *answerHandler) get(w http.ResponseWriter, r *http.Request) {
	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}
	domains, err := h.usecase.Get(r.Context(), auth.toIdentity())
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	resp := make([]answerShortResp, len(domains))
	for i, domain := range domains {
		resp[i] = answerShortRespFromDomain(domain)
	}
	if err := encodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *answerHandler) getByID(w http.ResponseWriter, r *http.Request) {
	answerID, err := h.getParamAnswerID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	auth, ok := h.getAuthData(r.Context())
	if !ok {
		h.sendGetAuthDataError(w)
		return
	}
	domain, err := h.usecase.GetByID(r.Context(), auth.toIdentity(), usecase.ID(answerID))
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	resp := answerFullRespFromDomain(domain)
	if err := encodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}
