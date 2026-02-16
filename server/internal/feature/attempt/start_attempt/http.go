package start_attempt

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/param"
	"server/internal/pkg/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quizID, err := param.ID(r, param.QuizID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		var input Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		output, err := usecase.StartAttempt(r.Context(), quizID, &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}
	}
}
