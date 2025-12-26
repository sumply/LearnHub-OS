package middleware

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"server/internal/logger"
	"server/internal/rest/transport"

	"github.com/google/uuid"
)

type ctxKey string

const tokenKey ctxKey = "token"

type TokenParser interface {
	Parse(token string) (transport.AuthData, error)
}

type StubTokenParser struct{}

func (p *StubTokenParser) Parse(token string) (transport.AuthData, error) {
	return transport.AuthData{ID: 1, Role: "root"}, nil
}

func loggerWithAuthData(log logger.Logger, auth transport.AuthData) logger.Logger {
	field := logger.TraceField{
		Key: "Auth",
		Value: map[string]any{
			"ID":   auth.ID,
			"Role": auth.Role,
		},
	}
	return log.With(field)
}

var authBearer = regexp.MustCompile(`^Bearer\s+(.+)$`)

func GetToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := getAuthorization(r.Header)
		matches := authBearer.FindStringSubmatch(auth)

		if len(matches) != 2 {
			const errMsg = `Authorization header must be "Bearer <token>"`
			transport.SendError(
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

func ValidateToken(p TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromCtx(r.Context())

			token, ok := r.Context().Value(tokenKey).(string)
			if !ok {
				const errMsg = "Authorization token is invalid"
				transport.SendError(
					w,
					http.StatusUnauthorized,
					errMsg,
				)
				return
			}

			auth, err := p.Parse(token)
			if err != nil {
				transport.SendError(
					w,
					http.StatusUnauthorized,
					err.Error(),
				)
				return
			}

			log = loggerWithAuthData(log, auth)

			ctx := auth.WithCtx(r.Context())
			ctx = logger.WithLoggerCtx(ctx, log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())
		log = loggerWithRequest(log, r)
		ctx := logger.WithLoggerCtx(r.Context(), log)
		log.Info("Request received")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriterWrapper(w)

		next.ServeHTTP(rw, r)

		log := logger.FromCtx(r.Context())
		loggerWithResponse(log, rw.statusCode, rw.size).
			Info("Response sent")
	})
}

func loggerWithRequest(log logger.Logger, r *http.Request) logger.Logger {
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
		auth := getAuthorization(header)
		if len(auth) != 0 {
			matches := authBearer.FindStringSubmatch(auth)
			if len(matches) != 2 {
				auth = "[INVALID_TOKEN_FORMAT]"
			} else {
				auth = "Bearer [MASKED_TOKEN]"
			}
			header.Set("Authorization", auth)
		}
		field.Value.(map[string]any)["Headers"] = header
	})
	return log.With(reqID, field)
}

func loggerWithResponse(log logger.Logger, status int, size int) logger.Logger {
	fields := logger.TraceField{
		Key: "Response",
		Value: map[string]any{
			"Status": status,
			"Size":   size,
		},
	}
	return log.With(fields)
}

func getAuthorization(h http.Header) string {
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
