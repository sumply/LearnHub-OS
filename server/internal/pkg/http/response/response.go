package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/pkg/http/param"
	"server/internal/pkg/usecase"

	"github.com/go-playground/validator/v10"
)

type errorMessage struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func (e errorMessage) Bytes() []byte {
	data, _ := json.Marshal(e)
	return data
}

func SendJSONEncodeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Printf("ENCODING ERROR: %v", err)
}

func SendJSONDecodeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(
		errorMessage{
			Error: err.Error(),
		}.Bytes(),
	)
}

func SendDTOValidateError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	var target validator.ValidationErrors
	var msg errorMessage
	if errors.As(err, &target) {
		msg.Error = "Validation"
		details := make(map[string]string)
		for _, e := range target {
			details[e.Field()] = e.ActualTag()
		}
		msg.Details = details
	} else {
		msg = errorMessage{
			Error: err.Error(),
		}
	}
	w.Write(msg.Bytes())
}

func SendParamError(w http.ResponseWriter, err error) {
	var queryErr *param.QueryError
	var msg errorMessage
	if errors.As(err, &queryErr) {
		msg.Error = queryErr.Error()
		msg.Details = queryErr.ToMap()
	} else {
		msg.Error = err.Error()
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(msg.Bytes())
}

func SendUseCaseError(w http.ResponseWriter, err error) {
	var msg errorMessage

	if target, ok := errors.AsType[*usecase.ValidationError](err); ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg = handleValidationError(target)
	} else if target, ok := errors.AsType[*usecase.AuthError](err); ok {
		w.WriteHeader(http.StatusForbidden)
		msg = errorMessage{
			Error:   target.Error(),
			Details: target.Unwrap().Error(),
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		msg = errorMessage{
			Error: err.Error(),
		}
	}

	w.Write(msg.Bytes())
}

func handleValidationError(err *usecase.ValidationError) errorMessage {
	msg := errorMessage{
		Error: err.Error(),
	}
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		msg.Details = domainErr.ToMap()
	} else {
		msg.Details = errors.Unwrap(err).Error()
	}
	return msg
}

func SendAuthTokenError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(errorMessage{
		Error:   "Unauthorized",
		Details: err.Error(),
	}.Bytes())
}
