package start_attempt

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/http/response"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decoder.JSON[Request](r.Body)
		if err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		resp, err := uc.StartAttempt(r.Context(), req)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, resp)
	}
}
