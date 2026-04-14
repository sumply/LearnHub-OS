package login

import (
	"net/http"
	"server/internal/pkg/decoder"
	"server/internal/pkg/http/response"
	"server/internal/pkg/validator"
	"server/pkg/logger"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		input, err := decoder.JSON[Input](r.Body)
		if err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendJSONDecodeError(w, err)
			return
		}

		if err := validator.V(r.Context(), input); err != nil {
			response.LogDTOValidateError(ctx, log)
			response.SendDTOValidateError(w, err)
			return
		}

		output, err := usecase.Login(r.Context(), input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendOK(w, output)
	}
}
