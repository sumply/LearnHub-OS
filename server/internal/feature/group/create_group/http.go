package create_group

import (
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		input, err := InputFromRequest(r)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendDTOValidateError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.LogTokenError(ctx, log)
			response.SendAuthTokenError(w, nil)
			return
		}

		output, err := uc.CreateGroup(r.Context(), token, &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, output)
	}
}
