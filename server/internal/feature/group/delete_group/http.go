package delete_group

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		groupID, err := param.ID(r, param.GroupID)
		if err != nil {
			response.LogParamError(ctx, log, param.GroupID, err)
			response.SendParamError(w, err)
			return
		}

		identity, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.LogTokenError(ctx, log)
			response.SendAuthTokenError(w, nil)
			return
		}

		err = uc.DeleteGroup(r.Context(), identity, groupID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
