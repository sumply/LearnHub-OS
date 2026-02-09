package get_students

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/param"
	"server/internal/pkg/response"
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

		if err := json.NewEncoder(w).Encode(output); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}
	}
}
