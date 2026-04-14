package create_subject

import (
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())
		ctx := r.Context()

		input, err := InputFromRequest(r)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendDTOValidateError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(ctx)
		if !ok {
			response.LogTokenError(ctx, log)
			response.SendAuthTokenError(w, nil)
			return
		}

		output, err := uc.CreateSubject(ctx, token, input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, output)
	}
}
