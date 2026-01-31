package dto

import (
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	ID        uuid.UUID  `json:"id"`
	User      User       `json:"user"`
	Answers   []Answer   `json:"answers"`
	Score     int        `json:"score"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

type Answer struct {
	ID        uuid.UUID `json:"id"`
	Question  Question  `json:"question"`
	Score     int       `json:"score"`
	Answer    string    `json:"answer"`
	IsCorrect bool      `json:"is_correct"`
}
