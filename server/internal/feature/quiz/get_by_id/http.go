package get_by_id

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quizID, err := param.ID(r, param.QuizID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		output, err := usecase.GetByID(r.Context(), quizID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
