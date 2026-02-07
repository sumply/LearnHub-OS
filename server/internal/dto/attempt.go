package dto

import (
	"time"

	"github.com/google/uuid"
)

type FinishedAttempt struct {
	Attempt Attempt `json:"attempt"`
	Quiz    Quiz    `json:"quiz"`
}

type Attempt struct {
	ID        uuid.UUID `json:"id"`
	User      User      `json:"user"`
	Answers   []Answer  `json:"answers"`
	Score     int       `json:"score"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

type AttemptItem struct {
	ID        uuid.UUID  `json:"id"`
	Score     int        `json:"score"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}

type Answer struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Answer     any       `json:"answer"`
	Score      int       `json:"score"`
	IsCorrect  bool      `json:"is_correct"`
}
