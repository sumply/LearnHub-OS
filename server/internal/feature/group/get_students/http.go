package get_students

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

		output, err := usecase.GetStudents(r.Context(), groupID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
