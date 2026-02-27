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
	"server/internal/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

func Router(creator config.Creator) http.Handler {
	r := chi.NewMux()

	prepareCommonMiddleware(r)

	postgresCfg := creator.CreatePostgresConnection()

	p, err := postgres.New(postgresCfg.CreateOptions())
	if err != nil {
		panic(err)
	}

	jwtCfg := creator.CreateJWT()

	registerPublicRoutes(r, jwtCfg, p)
	registerPrivateRoutes(r, jwtCfg, p)

	return r
}

func prepareCommonMiddleware(r chi.Router) {
	r.Use(
		middleware.CORS(),
		middleware.Logger(),
		middleware.BodyLogger(),
		middleware.Recoverer(),
	)
}

func registerPublicRoutes(r chi.Router, cfg config.JWT, postgres *postgres.Postgres) {
	jwtGenerator := jwt.NewGenerator(cfg.Issuer, cfg.Secret, cfg.AccessDur, cfg.RefreshDur)

	auth.Route(r, jwtGenerator, postgres)
}

func registerPrivateRoutes(r chi.Router, cfg config.JWT, postgres *postgres.Postgres) {
	jwtParser := jwt.NewParser(cfg.Secret)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtParser))

		user.Route(r, postgres)
		group.Route(r, postgres)
		subject.Route(r, postgres)
		quiz.Route(r, postgres)
		attempt.Route(r, postgres)
	})
}
