package middleware

import (
	"bytes"
	"log"
	"net/http"

	"io"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

type loggingResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b) // сохраняем тело
	return lrw.ResponseWriter.Write(b)
}

func Logger() Middleware {
	return middleware.Logger
}

func BodyLogger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestBody []byte
			if r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			log.Printf("Request: %s %s\nBody: %s\n", r.Method, r.URL.Path, string(requestBody))
			// Оборачиваем ResponseWriter
			lrw := &loggingResponseWriter{ResponseWriter: w, body: &bytes.Buffer{}}

			// Вызываем следующий обработчик
			next.ServeHTTP(lrw, r)

			// Логируем тело ответа
			log.Printf("Response: %s\nBody: %s\n", lrw.Header().Get("Status"), lrw.body.String())
		})
	}
}

func Recoverer() Middleware {
	return middleware.Recoverer
}

/*
func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := getAuthorizationToken(r)
			if err != nil {
				response.SendAuthTokenError(w, err)
				return
			}
		})
	}
}

var jwtRegex = regexp.MustCompile(`^Bearer\s+(.+)$`)

func getAuthorizationToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	token := jwtRegex.FindString(header)
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return token, nil
}

*/
