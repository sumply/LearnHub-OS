package get_quizzes

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

		opts, err := parseOptions(r)
		if err != nil {
			response.SendParamError(w, err)
			return
		}

		output, err := usecase.GetQuizzes(r.Context(), userID, opts)
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

func parseOptions(r *http.Request) (Options, error) {
	queryErr := param.NewQueryError(2)

	scope := OptionScope(r.URL.Query().Get("scope"))

	if !scope.Validate() {
		queryErr.Add(scope)
	}

	include := OptionInclude(r.URL.Query().Get("include"))

	if !include.Validate() {
		queryErr.Add(include)
	}

	if !queryErr.Empty() {
		return Options{}, queryErr
	}

	return Options{
		Scope:   scope,
		Include: include,
	}, nil
}
