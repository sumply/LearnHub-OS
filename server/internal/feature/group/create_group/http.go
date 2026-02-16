package create_group

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		output, err := usecase.CreateGroup(r.Context(), &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
