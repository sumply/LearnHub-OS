package list_quizzes

import (
	"net/http"
	"server/internal/pkg/http/response"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := uc.ListQuizzes(r.Context())
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, resp)
	}
}
