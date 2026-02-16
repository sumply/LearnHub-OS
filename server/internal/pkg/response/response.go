package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/pkg/param"
	"server/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type ErrorMessage struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func (e ErrorMessage) Bytes() []byte {
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
		ErrorMessage{
			Error: err.Error(),
		}.Bytes(),
	)
}

func SendDTOValidateError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	var target validator.ValidationErrors
	var msg ErrorMessage
	if errors.As(err, &target) {
		msg.Error = "Validation"
		details := make(map[string]string)
		for _, e := range target {
			details[e.Field()] = e.ActualTag()
		}
		msg.Details = details
	} else {
		msg = ErrorMessage{
			Error: err.Error(),
		}
	}
	w.Write(msg.Bytes())
}

func SendParamError(w http.ResponseWriter, err error) {
	var queryErr *param.QueryError
	var msg ErrorMessage
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
	var msg ErrorMessage

	var e *usecase.ValidationError
	if errors.As(err, &e) {
		w.WriteHeader(http.StatusBadRequest)
		msg = handleValidationError(e)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		msg = ErrorMessage{
			Error: err.Error(),
		}
	}

	w.Write(msg.Bytes())
}

func handleValidationError(err *usecase.ValidationError) ErrorMessage {
	msg := ErrorMessage{
		Error: err.Error(),
	}
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		msg.Details = domainErr.ToMap()
	} else {
		msg.Details = err.Error()
	}
	return msg
}
