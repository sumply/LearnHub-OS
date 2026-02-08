package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/domain"
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

func SendParamError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(
		ErrorMessage{
			Error: err.Error(),
		}.Bytes(),
	)
}

func SendUseCaseError(w http.ResponseWriter, err error) {
	var e *domain.Error
	var msg ErrorMessage
	if errors.As(err, &e) {
		w.WriteHeader(http.StatusBadRequest)
		msg = ErrorMessage{
			Error:   e.Domain() + " validation",
			Details: e.ToMap(),
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		msg = ErrorMessage{
			Error: err.Error(),
		}
	}
	w.Write(msg.Bytes())
}
