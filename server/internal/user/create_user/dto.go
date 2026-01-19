package create_user

import "github.com/google/uuid"

type Input struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
