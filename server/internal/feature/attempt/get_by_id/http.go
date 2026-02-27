package get_by_id

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		attemptID, err := param.ID(r, param.AttemptID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		output, err := usecase.GetByID(r.Context(), attemptID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
