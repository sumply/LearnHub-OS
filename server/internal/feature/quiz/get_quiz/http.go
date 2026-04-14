package get_quiz

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
	"server/pkg/logger"
)

func HTTP(uc *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		quizID, err := param.ID(r, param.QuizID)
		if err != nil {
			response.LogParamError(ctx, log, param.QuizID, err)
			response.SendParamError(w, err)
			return
		}

		resp, err := uc.GetQuiz(r.Context(), quizID)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, resp)
	}
}
