package get_subject

import "github.com/google/uuid"

type OutputSubject struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Output struct {
	Subjects []OutputSubject `json:"subjects"`
}
