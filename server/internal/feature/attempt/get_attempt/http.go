package get_attempt

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		attemptID, err := param.ID(r, param.AttemptID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		resp, err := uc.GetAttempt(r.Context(), attemptID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, resp)
	}
}
