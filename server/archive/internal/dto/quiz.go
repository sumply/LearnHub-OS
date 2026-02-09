package dto

import (
	"server/archive/internal/common"
	"server/archive/internal/domain"
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

type QuizBeforeExecResp struct {
	ID         common.ID                 `json:"id"`
	Title      string                    `json:"title"`
	Summary    string                    `json:"summary"`
	TotalScore int                       `json:"total_score"`
	Owner      *UserShortResp            `json:"owner"`
	Subject    *SubjectResp              `json:"subject"`
	Groups     []*GroupShortResp         `json:"group,omitempty"`
	Questions  []*questionBeforeExecResp `json:"questions"`
}

func NewQuizBeforeExecResp(d *domain.Quiz) *QuizBeforeExecResp {
	if d == nil {
		return nil
	}
	return &QuizBeforeExecResp{
		ID:         d.ID,
		Title:      d.Title,
		Summary:    d.Summary,
		TotalScore: d.TotalScore,
		Owner:      NewUserShortResp(d.Owner),
		Subject:    NewSubjectResp(d.Subject),
		Groups:     NewSliceGroupShortResp(d.Groups),
		Questions:  newSliceQuestionBeforeExecResp(d.Questions),
	}
}

type questionBeforeExecResp struct {
	ID      common.ID               `json:"id"`
	Text    string                  `json:"text"`
	Options []*optionBeforeExecResp `json:"options"`
}

func newQuestionBeforeExecResp(d *domain.Question) *questionBeforeExecResp {
	if d == nil {
		return nil
	}
	return &questionBeforeExecResp{
		ID:      d.ID,
		Text:    d.Text,
		Options: newSliceOptionBeforeExecResp(d.Options),
	}
}

func newSliceQuestionBeforeExecResp(d []*domain.Question) []*questionBeforeExecResp {
	resp := make([]*questionBeforeExecResp, len(d))
	for i := range resp {
		resp[i] = newQuestionBeforeExecResp(d[i])
	}
	return resp
}

type optionBeforeExecResp struct {
	ID   common.ID `json:"id"`
	Text string    `json:"text"`
}

func newOptionBeforeExecResp(d *domain.Option) *optionBeforeExecResp {
	if d == nil {
		return nil
	}
	return &optionBeforeExecResp{
		ID:   d.ID,
		Text: d.Text,
	}
}

func newSliceOptionBeforeExecResp(d []*domain.Option) []*optionBeforeExecResp {
	resp := make([]*optionBeforeExecResp, len(d))
	for i := range resp {
		resp[i] = newOptionBeforeExecResp(d[i])
	}
	return nil
}

func NewQuizShortResp(d *domain.Quiz) *QuizShortResp {
	if d == nil {
		return nil
	}
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
