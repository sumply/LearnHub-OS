package create_user

import (
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())

		input, err := InputFromRequest(r)
		if err != nil {
			log.WarnContext(r.Context(), "Failed creating request dto")
			response.SendDTOValidateError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			log.WarnContext(r.Context(), "Failed getting identity from context")
			response.SendAuthTokenError(w, nil)
			return
		}

		output, err := uc.CreateUser(r.Context(), token, input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, output)
	}
}
