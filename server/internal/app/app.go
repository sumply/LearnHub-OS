package app

import (
	"fmt"
	"net/http"
	"server/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	s, err := config.NewServer()
	if err != nil {
		return err
	}
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Addr, s.Port), r)
	return nil
}
