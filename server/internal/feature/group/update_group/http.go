package update_group

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/pkg/logger"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		groupID, err := param.ID(r, param.GroupID)
		if err != nil {
			response.LogParamError(ctx, log, param.GroupID, err)
			response.SendParamError(w, err)
			return
		}

		input, err := InputFromRequest(r)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendDTOValidateError(w, err)
			return
		}

		err = usecase.UpdateGroup(r.Context(), groupID, &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
