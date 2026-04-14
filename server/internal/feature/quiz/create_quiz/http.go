package create_quiz

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/http/response"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		req, err := decoder.JSON[Request](r.Body)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendJSONDecodeError(w, err)
			return
		}

		resp, err := uc.CreateQuiz(r.Context(), req)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, resp)
	}
}
