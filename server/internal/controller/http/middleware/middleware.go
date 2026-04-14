package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"server/internal/pkg/http/response"
	"server/internal/pkg/jwt"
	"server/internal/pkg/usecase"
	"server/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type Middleware func(next http.Handler) http.Handler

func CORS() Middleware {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriterWrapper) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithAttrs(r.Context(), slog.Group("request",
				slog.String("id", uuid.New().String()),
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.String("host", r.Host),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Time("received_time", time.Now().UTC()),
			))
			logger.Info(ctx, "Received request")

			startTime := time.Now()

			r = r.WithContext(ctx)

			rw := &ResponseWriterWrapper{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(startTime)
			logger.Info(r.Context(), "Sended response", slog.Group("response",
				slog.Int("status", rw.status),
				slog.Int64("duration_ms", duration.Milliseconds()),
			))
		})
	}
}

func Recoverer() Middleware {
	return middleware.Recoverer
}

func Auth(parser *jwt.Parser) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getAuthorizationToken(r)
			if err != nil {
				response.SendAuthTokenError(w, err)
				return
			}

			claims, err := parser.Parse(token)
			if err != nil {
				response.SendAuthTokenError(w, err)
				return
			}

			ctx := usecase.IdentityWithContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var jwtRegex = regexp.MustCompile(`^Bearer\s+(.+)$`)

func getAuthorizationToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	token := jwtRegex.FindStringSubmatch(header)
	if len(token) != 2 {
		return "", fmt.Errorf("token is invalid")
	}
	return token[1], nil
}
