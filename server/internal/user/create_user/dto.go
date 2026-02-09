package create_user

import "github.com/google/uuid"

type Input struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`

	GroupID    uuid.UUID  `json:"group_id"`
	GroupIDs   uuid.UUIDs `json:"group_ids"`
	SubjectIDs uuid.UUIDs `json:"subject_ids"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
