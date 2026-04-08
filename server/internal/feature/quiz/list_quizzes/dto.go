package list_quizzes

import (
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Quizzes []ResponseQuiz `json:"quizzes"`
}

type ResponseQuiz struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	SubjectID   uuid.UUID  `json:"subject_id"`
	GroupIDs    uuid.UUIDs `json:"group_ids"`
	Deadline    *time.Time `json:"deadline"`
	MaxAttempts int        `json:"max_attempts"`
	TotalScore  int        `json:"total_score"`
	CreatedAt   time.Time  `json:"created_at"`
}

func NewResponseQuiz(d *domain.Quiz) ResponseQuiz {
	return ResponseQuiz{
		ID:          d.ID,
		Title:       d.Title,
		Summary:     d.Summary,
		OwnerID:     d.OwnerID,
		SubjectID:   d.SubjectID,
		GroupIDs:    d.GroupIDs,
		Deadline:    d.Deadline,
		MaxAttempts: d.MaxAttempts,
		TotalScore:  d.TotalScore,
		CreatedAt:   d.CreatedAt,
	}
}
