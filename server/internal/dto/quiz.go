package dto

import (
	"server/internal/common"
	"server/internal/domain"
)

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

type QuizShortResp struct {
	ID         common.ID         `json:"id"`
	Title      string            `json:"title"`
	Summary    string            `json:"summary"`
	TotalScore int               `json:"total_score"`
	Owner      *UserShortResp    `json:"owner"`
	Subject    *SubjectResp      `json:"subject"`
	Groups     []*GroupShortResp `json:"group,omitempty"`
}

func NewQuizShortResp(d *domain.Quiz) *QuizShortResp {
	return &QuizShortResp{
		ID:         d.ID,
		Title:      d.Title,
		Summary:    d.Summary,
		TotalScore: d.TotalScore,
		Owner:      NewUserShortResp(d.Owner),
		Subject:    NewSubjectResp(d.Subject),
		Groups:     NewSliceGroupShortResp(d.Groups),
	}
}

func NewSliceQuizShortResp(domains []*domain.Quiz) []*QuizShortResp {
	resp := make([]*QuizShortResp, len(domains))
	for i, d := range domains {
		resp[i] = NewQuizShortResp(d)
	}
	return resp
}
