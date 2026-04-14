package delete_user

import (
	"log/slog"
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())

		userID, err := param.ID(r, param.UserID)
		if err != nil {
			log.WarnContext(r.Context(), "Failed getting user_id param", slog.String("error", err.Error()))
			response.SendParamError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			log.WarnContext(r.Context(), "Failed getting identity from context")
			response.SendAuthTokenError(w, nil)
			return
		}

		err = uc.DeleteUser(r.Context(), token, userID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
