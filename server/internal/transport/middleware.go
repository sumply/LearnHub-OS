package transport

import (
	"bytes"
	"context"
	"io"
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
	subject id
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
		log := logger.FromCtx(r.Context())

		authHeader := r.Header.Get("Authorization")
		matches := authBearer.FindStringSubmatch(authHeader)

		if len(matches) != 2 {
			const errMsg = `Authorization header must be "Bearer <token>"`
			loggerWithResponse(log, http.StatusBadRequest, errMsg).
				Warn("authorization header missing or invalid")
			sendError(
				w,
				http.StatusBadRequest,
				errMsg,
			)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, tokenKey, matches[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateTokenMiddleware(p TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromCtx(r.Context())

			token, ok := r.Context().Value(tokenKey).(string)
			if !ok {
				const errMsg = "Authorization token is invalid"
				loggerWithResponse(log, http.StatusUnauthorized, errMsg).
					Warn("authorization token from context is missing")
				sendError(
					w,
					http.StatusUnauthorized,
					errMsg,
				)
				return
			}

			auth, ok := p.Parse(token)
			if !ok {
				const errMsg = "Authorization token is expired"
				loggerWithResponse(log, http.StatusUnauthorized, errMsg).
					Warn("authorization token failed verification")
				sendError(
					w,
					http.StatusUnauthorized,
					errMsg,
				)
				return
			}

			log = loggerWithAuthData(log, auth)

			ctx := context.WithValue(r.Context(), authKey, auth)
			ctx = logger.WithLoggerCtx(ctx, log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func loggingRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())
		log = loggerWithRequest(log, r)
		ctx := logger.WithLoggerCtx(r.Context(), log)
		log.Info("Request received")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerWithRequest(log logger.Logger, r *http.Request) logger.Logger {
	fields := logger.TraceField{
		Key: "Request",
		Value: map[string]any{
			"URI":    r.RequestURI,
			"Method": r.Method,
		},
	}
	logger.OnDebug(func() {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		fields.Value.(map[string]any)["Body"] = string(data)
	})
	return log.With(fields)
}

func loggerWithResponse(log logger.Logger, status int, what string) logger.Logger {
	fields := logger.TraceField{
		Key: "Response",
		Value: map[string]any{
			"Status": status,
			"What":   what,
		},
	}
	return log.With(fields)
}
