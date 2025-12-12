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

type answerCreateReq struct {
	Answers []answerReq `json:"answers"`
}

type answerShortResp struct {
	ID            id            `json:"id"`
	Quiz          quizShortResp `json:"quiz"`
	TotalScore    int           `json:"total_score"`
	Score         int           `json:"score"`
	Completed     userShortResp `json:"completed"`
	CompletedTime time.Time     `json:"completed_time"`
	Group         groupResp     `json:"group"`
}

type answerFullResp struct {
	ID            id            `json:"id"`
	TotalScore    int           `json:"total_score"`
	Score         int           `json:"score"`
	Completed     bool          `json:"completed"`
	CompletedTime time.Time     `json:"completed_time"`
	User          userShortResp `json:"user"`
	Quiz          quizShortResp `json:"quiz"`
	Answers       []answerResp  `json:"answers"`
}

type answerResp struct {
	Question quizQuestionResp `json:"question"`
	Answer   quizOptionsResp  `json:"answer"`
}

type answerHandler struct {
	usecase usecase.QuizResult
}

func newAnswerHandler(u usecase.QuizResult) (*answerHandler, error) {
	if u == nil {
		return nil, fmt.Errorf("не передана реализация интерфейса")
	}
	return &answerHandler{
		usecase: u,
	}, nil
}

func (h *answerHandler) post(w http.ResponseWriter, r *http.Request) {
	var req answerCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *answerHandler) get(w http.ResponseWriter, r *http.Request) {
	var resp []answerShortResp
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *answerHandler) getByID(w http.ResponseWriter, r *http.Request) {
	_, err := getParamQuizID(r)
	if err != nil {
		sendParamError(w, err.Error())
		return
	}
	_, err = getParamResultID(r)
	if err != nil {
		sendParamError(w, err.Error())
		return
	}
	var resp []answerFullResp
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}
