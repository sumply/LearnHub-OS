package dto

import (
	"server/internal/common"
	"server/internal/domain"
	"time"
)

type ProgressCreateReq struct {
	QuizID common.ID `json:"quiz_id"`
}

type AnswerPatchReq struct {
	Text string `json:"text"`
}

type ProgressShortResp struct {
	ID            common.ID             `json:"id"`
	Quiz          *QuizShortResp        `json:"quiz"`
	User          *UserShortResp        `json:"user"`
	Status        domain.ProgressStatus `json:"status"`
	Score         int                   `json:"score"`
	CompletedDate *time.Time            `json:"completed_date,omitempty"`
	StartDate     *time.Time            `json:"start_date,omitempty"`
}

func NewProgressShortResp(d *domain.QuizProgress) *ProgressShortResp {
	if d == nil {
		return nil
	}
	return &ProgressShortResp{
		ID:            d.ID,
		Quiz:          NewQuizShortResp(d.Quiz),
		User:          NewUserShortResp(d.User),
		Status:        d.Status,
		Score:         d.Score,
		CompletedDate: d.CompletedDate,
		StartDate:     d.StartDate,
	}
}
