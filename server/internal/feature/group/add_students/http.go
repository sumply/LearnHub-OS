package add_students

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

		input, err := InputFromRequest(r)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendDTOValidateError(w, err)
			return
		}
		input.GroupID = groupID

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.LogTokenError(ctx, log)
			response.SendAuthTokenError(w, nil)
			return
		}

		if err := uc.AddStudents(r.Context(), token, input); err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
