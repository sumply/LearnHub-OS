package query

import (
	"time"

	"github.com/google/uuid"
)

type QuizItem struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Owner      User      `json:"owner"`
	Subject    Subject   `json:"subject"`
	TotalScore int       `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
}

type Quiz struct {
	QuizItem
	Content []QuizQuestion `json:"content"`
}

type QuizQuestion struct {
	ID      uuid.UUID      `json:"id"`
	Text    string         `json:"text"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}
