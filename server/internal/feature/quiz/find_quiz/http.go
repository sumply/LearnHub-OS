package find_quiz

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/response"
)

func HTTPForStudent(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		resp, err := usecase.FindForStudent(r.Context(), &request)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}
	}
}
