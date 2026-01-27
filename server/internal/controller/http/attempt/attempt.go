package attempt

import (
	"server/internal/adapter/postgres"
	"server/internal/adapter/redis"
	"server/internal/attempt/start_attempt"

	"github.com/go-chi/chi/v5"
)

func Route(router chi.Router, p *postgres.Postgres, r *redis.Redis) {
	startUC := start_attempt.New(p, r)

	router.Post("/attempts", start_attempt.HTTP(startUC))
}
