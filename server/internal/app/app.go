package app

import (
	"net/http"
	myhttp "server/internal/controller/http"
)

func Run() error {
	return http.ListenAndServe("0.0.0.0:8000", myhttp.Router())
}
