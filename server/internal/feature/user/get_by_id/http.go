package get_by_id

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/param"
	"server/internal/pkg/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := param.ID(r, param.UserID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		output, err := usecase.GetByID(r.Context(), userID)
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
