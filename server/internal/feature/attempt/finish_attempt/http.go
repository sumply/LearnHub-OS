package finish_attempt

import (
	"net/http"
	"server/internal/pkg/decoder"
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

		req, err := decoder.JSON[Request](r.Body)
		if err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		resp, err := uc.FinishAttempt(r.Context(), attemptID, req)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, resp)
	}
}
