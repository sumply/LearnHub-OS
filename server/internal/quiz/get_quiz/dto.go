package get_quiz

import (
	"time"

	"github.com/google/uuid"
)

type Output struct {
	Quizzes []OutputQuiz `json:"quizzes"`
}

type OutputQuiz struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	OwnerID    uuid.UUID `json:"owner_id"`
	SubjectID  uuid.UUID `json:"subject_id"`
	Content    any       `json:"content"`
	TotalScore int       `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
}
