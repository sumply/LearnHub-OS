package create_quiz

import "github.com/google/uuid"

type Input struct {
	Title     string           `json:"title"`
	Summary   string           `json:"summary"`
	SubjectID uuid.UUID        `json:"subject_id"`
	Content   []map[string]any `json:"content"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
