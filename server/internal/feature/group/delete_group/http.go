package delete_group

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, err := param.ID(r, param.GroupID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		identity, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.SendAuthTokenError(w, nil)
			return
		}

		err = uc.DeleteGroup(r.Context(), identity, groupID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
