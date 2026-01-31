package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorMessage struct {
	Error error `json:"error"`
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
			Error: err,
		}.Bytes(),
	)
}

func SendParamError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(
		ErrorMessage{
			Error: err,
		}.Bytes(),
	)
}

func SendUseCaseError(w http.ResponseWriter, err error) {

}
