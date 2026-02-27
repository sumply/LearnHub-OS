package delete_user

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := param.ID(r, param.UserID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
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
