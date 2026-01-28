package dto

import (
	"time"

	"github.com/google/uuid"
)

type QuizItem struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Owner       User       `json:"owner"`
	Subject     Subject    `json:"subject"`
	TotalScore  int        `json:"total_score"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	MaxAttempts int        `json:"max_attempts"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Quiz struct {
	QuizItem
	Content []QuizQuestion `json:"content"`
}

func (q *Quiz) DeleteAnswers() {
	for i := range q.Content {
		q.Content[i].DeleteAnswer()
	}
}

type QuizQuestion struct {
	ID      uuid.UUID      `json:"id"`
	Text    string         `json:"text"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (q *QuizQuestion) DeleteAnswer() {
	delete(q.Payload, "correct")
}
