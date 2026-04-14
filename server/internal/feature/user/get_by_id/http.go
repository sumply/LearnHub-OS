package get_by_id

import (
	"log/slog"
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/pkg/logger"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())

		userID, err := param.ID(r, param.UserID)
		if err != nil {
			log.WarnContext(r.Context(), "Failed getting user_id param", slog.String("error", err.Error()))
			response.SendParamError(w, err)
			return
		}

		output, err := usecase.GetByID(r.Context(), userID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
