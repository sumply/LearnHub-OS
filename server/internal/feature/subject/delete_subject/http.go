package delete_subject

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

		subjectID, err := param.ID(r, param.SubjectID)
		if err != nil {
			response.LogParamError(ctx, log, param.SubjectID, err)
			response.SendParamError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.LogTokenError(ctx, log)
			response.SendAuthTokenError(w, nil)
			return
		}

		err = uc.DeleteSubject(r.Context(), token, subjectID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
