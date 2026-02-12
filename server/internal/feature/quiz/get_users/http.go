package get_users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/pkg/param"
	"server/internal/pkg/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("callind http")
		quizID, err := param.ID(r, param.QuizID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}
		fmt.Printf("quiz_id = %v\n", quizID)

		output, err := usecase.GetUsers(r.Context(), quizID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		fmt.Println(output)

		if err := json.NewEncoder(w).Encode(output); err != nil {
			response.SendJSONEncodeError(w, err)
			return
		}
	}
}
