package update_group

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/validator"
)

type Input struct {
	Name string `json:"name" validate:"required,group-name"`
}

func InputFromRequest(r *http.Request) (Input, error) {
	input, err := decoder.JSON[*Input](r.Body)
	if err != nil {
		return Input{}, err
	}

	err = validator.V(r.Context(), input)
	if err != nil {
		return Input{}, err
	}

	return *input, nil
}
