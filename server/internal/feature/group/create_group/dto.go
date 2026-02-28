package create_group

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/validator"

	"github.com/google/uuid"
)

type Input struct {
	Name string `json:"name" validate:"required,group-name"`

	CuratorID  uuid.UUID  `json:"curator_id"`
	StudentIDs uuid.UUIDs `json:"student_ids"`
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
