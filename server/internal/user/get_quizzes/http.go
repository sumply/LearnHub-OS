package get_quizzes

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := usecase.GetQuizzes(r.Context())
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(output); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}
	}
}
