package transport

import (
	"net/http"
)

type quizCreateReq struct {
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	SubjectID int64  `json:"subject_id"`
	Questions []struct {
		Text    string `json:"text"`
		Options []struct {
			Text      string `json:"text"`
			IsCorrect bool   `json:"is_correct"`
		}
	}
}

type quizShortResp struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Subject struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"subject"`
	Owner struct {
		ID        int64  `json:"id"`
		ShortName string `json:"short_name"`
	} `json:"owner"`
}

type quizFullResp struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Subject struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"subject"`
	Owner     userShortResp
	Questions []struct {
		ID      int64  `json:"id"`
		Text    string `json:"text"`
		Options []struct {
			ID   int64  `json:"id"`
			Text string `json:"text"`
		}
	}
}

type quizResultCreateReq struct {
	Answers []struct {
		QuestionID int64 `json:"question_id"`
		OptionID   int64 `json:"option_id"`
	} `json:"answers"`
}

type quizHandler struct {
}

func newQuizHandler() *quizHandler {
	return &quizHandler{}
}

func (h *quizHandler) post(w http.ResponseWriter, r *http.Request) {
	var req quizCreateReq
	if err := decodeJSON(r.Body, &req); err != nil {
		sendDecodeError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *quizHandler) get(w http.ResponseWriter, r *http.Request) {
	var resp quizShortResp
	if err := encodeJSON(w, &resp); err != nil {
		sendEncodeError(w)
		return
	}
}

func (h *quizHandler) getByID(w http.ResponseWriter, r *http.Request) {
	var resp quizFullResp
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

}
