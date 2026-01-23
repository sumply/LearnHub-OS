package create_quiz

import "github.com/google/uuid"

type Input struct {
	Title     string         `json:"title"`
	Summary   string         `json:"summary"`
	SubjectID uuid.UUID      `json:"subject_id"`
	Content   []InputContent `json:"content"`
}

type InputContent struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Payload any    `json:"payload"`
	Score   int    `json:"score"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
