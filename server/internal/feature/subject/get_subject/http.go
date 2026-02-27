package get_subject

import (
	"net/http"
	"server/internal/pkg/http/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := usecase.GetSubject(r.Context())
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
