package transport

import (
	"context"
	"net/http"
	"regexp"
	"server/internal/logger"
	"server/internal/usecase"
)

type TokenParser interface {
	Parse(token string) (authData, bool)
}

type StubTokenParser struct{}

func (p *StubTokenParser) Parse(token string) (authData, bool) {
	return authData{subject: 1, role: "admin"}, true
}

type ctxKey string

type authData struct {
	subject int64
	role    string
}

func loggerWithAuthData(log logger.Logger, auth authData) logger.Logger {
	field := logger.TraceField{
		Key: "Auth",
		Value: map[string]any{
			"ID":   auth.subject,
			"Role": auth.role,
		},
	}
	return log.With(field)
}

func (a *authData) toIdentity() usecase.Identity {
	r, ok := toUserRole(a.role)
	if !ok {
		return usecase.Identity{}
	}
	return usecase.Identity{
		ID:   usecase.ID(a.subject),
		Role: r,
	}
}

var authBearer = regexp.MustCompile(`^Bearer\s+(.+)$`)

const (
	tokenKey ctxKey = "token"
	authKey  ctxKey = "authData"
)

func getTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromCtx(r.Context())
		logger.Debug("Получение токена из головы запроса.")

		authHeader := r.Header.Get("Authorization")
		matches := authBearer.FindStringSubmatch(authHeader)

		if len(matches) != 2 {
			sendError(
				w,
				http.StatusBadRequest,
				`В запросе отсутствует заголовок типа "Authorization: Bearer <token>"`,
			)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, tokenKey, matches[1])

		logger.Debug("Токен успешно получен.")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateTokenMiddleware(p TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromCtx(r.Context())
			log.Debug("Валидация токена.")

			token, ok := r.Context().Value(tokenKey).(string)
			if !ok {
				sendError(
					w,
					http.StatusUnauthorized,
					"Не удалось получить токен.",
				)
				return
			}

			auth, ok := p.Parse(token)
			if !ok {
				sendError(
					w,
					http.StatusUnauthorized,
					"Невалидный токен.",
				)
				return
			}

			log = loggerWithAuthData(log, auth)
			log.Debug("Токен верный.")

			ctx := context.WithValue(r.Context(), authKey, auth)
			ctx = logger.WithLoggerCtx(ctx, log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
