package http

import (
	"net/http"
	"server/internal/adapter/postgres"
	"server/internal/adapter/redis"
	"server/internal/controller/http/attempt"
	"server/internal/controller/http/group"
	"server/internal/controller/http/quiz"
	"server/internal/controller/http/subject"
	"server/internal/controller/http/user"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewMux()

	p, err := postgres.New(postgres.Options{
		User:     "postgres",
		Password: "2121",
		DB:       "test",
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}
	redis := redis.Redis{}

	user.Route(r, p)
	group.Route(r, p)
	subject.Route(r, p)
	quiz.Route(r, p)
	attempt.Route(r, p, &redis)

	return r
}
