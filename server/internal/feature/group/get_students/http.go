package get_students

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, err := param.ID(r, param.GroupID)
		if err != nil {
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
