package update_user

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/validator"

	"github.com/google/uuid"
)

type Input struct {
	ID        uuid.UUID `json:"-"`
	FirstName string    `json:"first_name" validate:"required,user-name"`
	LastName  string    `json:"last_name" validate:"required,user-name"`
}

func InputFromRequest(r *http.Request) (Input, error) {
	var input Input

	id, err := param.ID(r, param.UserID)
	if err != nil {
		return Input{}, err
	}
	input.ID = id

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return Input{}, err
	}

	if err := validator.V(r.Context(), input); err != nil {
		return Input{}, err
	}

	return input, nil
}
