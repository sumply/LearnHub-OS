package http

import (
	"net/http"
	"server/internal/adapter/postgres"
	"server/internal/config"
	"server/internal/controller/http/attempt"
	"server/internal/controller/http/auth"
	"server/internal/controller/http/group"
	"server/internal/controller/http/middleware"
	"server/internal/controller/http/quiz"
	"server/internal/controller/http/subject"
	"server/internal/controller/http/user"

	"github.com/go-chi/chi/v5"
)

func Router(creator config.Creator) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.CORS(),
		middleware.Logger(),
		middleware.BodyLogger(),
		middleware.Recoverer(),
	)

	cfgPostgres := creator.CreatePostgresConnection()

	p, err := postgres.New(cfgPostgres.CreateOptions())
	if err != nil {
		panic(err)
	}

	user.Route(r, p)
	group.Route(r, p)
	subject.Route(r, p)
	quiz.Route(r, p)
	attempt.Route(r, p)
	auth.Route(r, p)

	return r
}
