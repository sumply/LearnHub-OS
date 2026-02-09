package create_subject

import "github.com/google/uuid"

type Input struct {
	Name string `json:"name"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
