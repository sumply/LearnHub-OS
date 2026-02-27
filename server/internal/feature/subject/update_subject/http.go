package update_subject

import (
	"net/http"
	"server/internal/pkg/http/param"
	"server/internal/pkg/http/response"
)

func HTTP(usecase *UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subjectID, err := param.ID(r, param.SubjectID)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		input, err := InputFromRequest(r)
		if err != nil {
			response.SendDTOValidateError(w, err)
			return
		}

		err = usecase.UpdateSubject(r.Context(), subjectID, &input)
		if err != nil {
			response.SendUseCaseError(w, err)
			return
		}

		response.SendNoContent(w)
	}
}
