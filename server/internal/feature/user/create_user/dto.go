package create_user

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/validator"

	"github.com/google/uuid"
)

type Input struct {
	FirstName string `json:"first_name" validate:"required,user-name"`
	LastName  string `json:"last_name" validate:"required,user-name"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role"`

	GroupID    uuid.UUID  `json:"group_id"`
	GroupIDs   uuid.UUIDs `json:"group_ids"`
	SubjectIDs uuid.UUIDs `json:"subject_ids"`
}

func InputFromRequest(r *http.Request) (Input, error) {
	input, err := decoder.JSON[Input](r.Body)
	if err != nil {
		return Input{}, err
	}

	if err := validator.V(r.Context(), input); err != nil {
		return Input{}, err
	}

	return input, nil
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
