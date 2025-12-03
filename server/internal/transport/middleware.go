package transport

import (
	"context"
	"errors"
	"net/http"
	"regexp"
)

type TokenParser interface {
	Parse(token string) (authData, bool)
}

type FakeTokenParser struct{}

func (p *FakeTokenParser) Parse(token string) (authData, bool) {
	return authData{subject: 1, role: "admin"}, false
}

type ctxKey string

type authData struct {
	subject int64
	role    string
}

var authBearer = regexp.MustCompile(`^Bearer\s+(.+)$`)

const (
	tokenKey ctxKey = "token"
	authKey  ctxKey = "authData"
)

var (
	errUnauthorized = errors.New("Authorization token is invalid")
)

func getTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func makeValidateTokenFunc(p TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			ctx := context.WithValue(r.Context(), authKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendError(w http.ResponseWriter, statusCode int, what string) {
	body := map[string]any{
		"error": what,
	}
	w.WriteHeader(statusCode)
	encodeJSON(w, body)
}
