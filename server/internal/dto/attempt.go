package dto

import (
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	ID        uuid.UUID       `json:"id"`
	Score     int             `json:"score"`
	User      User            `json:"user"`
	Answers   []AttemptAnswer `json:"answers"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
}

type AttemptAnswer struct {
	ID        uuid.UUID    `json:"id"`
	Question  QuizQuestion `json:"question"`
	Score     int          `json:"score"`
	Answer    string       `json:"answer"`
	IsCorrect bool         `json:"is_correct"`
}
