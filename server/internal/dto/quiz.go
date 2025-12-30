package dto

import "server/internal/common"

type OptionCreateReq struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionCreateReq struct {
	Text    string             `json:"text"`
	Options []*OptionCreateReq `json:"options,omitempty"`
}

type QuizCreateReq struct {
	Title     string               `json:"title"`
	Summary   string               `json:"summary"`
	Questions []*QuestionCreateReq `json:"questions"`
	SubjectID common.ID            `json:"subject_id"`
	GroupIDs  []common.ID          `json:"group_ids,omitempty"`
}
