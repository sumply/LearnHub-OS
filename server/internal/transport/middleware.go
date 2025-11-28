package transport

import (
	"context"
	"net/http"
	"regexp"
)

type ctxKey string

type authData struct {
	subject int64
	role    string
}

var authBearer = regexp.MustCompile(`^Bearer\s+(.+)$`)

const (
	token ctxKey = "token"
	auth  ctxKey = "authData"
)

func getTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		matches := authBearer.FindStringSubmatch(authHeader)

		if len(matches) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, token, matches[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(token).(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// *Логика обработки токена
		ctx := r.Context()
		ctx = context.WithValue(ctx, auth, authData{subject: 10, role: "admin"}) // пример результата

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
