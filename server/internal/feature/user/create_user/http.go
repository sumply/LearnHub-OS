package create_user

import (
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := InputFromRequest(r)
		if err != nil {
			logger.Warn(r.Context(), "Failed creating request dto")
			response.SendDTOValidateError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.SendAuthTokenError(w, nil)
			return
		}

		output, err := uc.CreateUser(r.Context(), token, input)
		if err != nil {
			logger.Error(r.Context(), "Failed sending response dto")
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, output)
	}
}
