package transport

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"server/internal/logger"
	"server/internal/usecase"

	"github.com/google/uuid"
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

type middlewareBuilder struct{}

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

func (b *middlewareBuilder) buildGetToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := b.getAuthorization(r.Header)
		matches := authBearer.FindStringSubmatch(auth)

		if len(matches) != 2 {
			const errMsg = `Authorization header must be "Bearer <token>"`
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

func (b *middlewareBuilder) buildValidateToken(p TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromCtx(r.Context())

			token, ok := r.Context().Value(tokenKey).(string)
			if !ok {
				const errMsg = "Authorization token is invalid"
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

func (b *middlewareBuilder) buildLoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())
		log = b.loggerWithRequest(log, r)
		ctx := logger.WithLoggerCtx(r.Context(), log)
		log.Info("Request received")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *middlewareBuilder) buildLoggingResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriterWrapper(w)

		next.ServeHTTP(rw, r)

		log := logger.FromCtx(r.Context())
		b.loggerWithResponse(log, rw.statusCode, rw.size).
			Info("Response sent")
	})
}

func (b *middlewareBuilder) loggerWithRequest(log logger.Logger, r *http.Request) logger.Logger {
	reqID := logger.TraceField{
		Key:   "Request_id",
		Value: uuid.New().String(),
	}
	field := logger.TraceField{
		Key: "Request",
		Value: map[string]any{
			"URI":    r.RequestURI,
			"Method": r.Method,
		},
	}
	logger.OnDebug(func() {
		header := r.Header.Clone()
		auth := b.getAuthorization(header)
		matches := authBearer.FindStringSubmatch(auth)
		if len(matches) != 2 {
			auth = "[INVALID_TOKEN_FORMAT]"
		} else {
			auth = "Bearer [MASKED_TOKEN]"
		}
		header.Set("Authorization", auth)
		field.Value.(map[string]any)["Headers"] = header
	})
	return log.With(reqID, field)
}

func (b *middlewareBuilder) loggerWithResponse(log logger.Logger, status int, size int) logger.Logger {
	fields := logger.TraceField{
		Key: "Response",
		Value: map[string]any{
			"Status": status,
			"Size":   size,
		},
	}
	return log.With(fields)
}

func (b *middlewareBuilder) getAuthorization(h http.Header) string {
	return h.Get("Authorization")
}

type responseWriterWrapper struct {
	rw         http.ResponseWriter
	statusCode int
	size       int
	body       bytes.Buffer
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		rw:         w,
		statusCode: http.StatusOK,
	}
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.rw.Header()
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.rw.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	logger.OnDebug(func() {
		w.body.Write(data)
	})
	size, err := w.rw.Write(data)
	w.size += size
	return size, err
}
