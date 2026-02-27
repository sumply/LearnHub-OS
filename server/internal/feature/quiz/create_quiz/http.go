package create_quiz

import (
	"net/http"
	"server/internal/pkg/http/response"
	"server/internal/pkg/usecase"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := InputFromRequest(r)
		if err != nil {
			response.SendDTOValidateError(w, err)
			return
		}

		token, ok := usecase.IdentityFromContext(r.Context())
		if !ok {
			response.SendAuthTokenError(w, nil)
			return
		}

		output, err := uc.CreateQuiz(r.Context(), token, &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendCreated(w, output)
	}
}
