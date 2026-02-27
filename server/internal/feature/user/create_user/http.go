package create_user

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.SendJSONDecodeError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
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
