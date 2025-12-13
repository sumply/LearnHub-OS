package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func decodeJSON(r io.ReadCloser, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

func encodeJSON(w io.Writer, v any) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func getAuthData(ctx context.Context) (authData, bool) {
	auth, ok := ctx.Value(authKey).(authData)
	return auth, ok
}

func sendError(w http.ResponseWriter, statusCode int, what string) {
	body := map[string]any{
		"error": what,
	}
	w.WriteHeader(statusCode)
	encodeJSON(w, body)
}

func sendDecodeError(w http.ResponseWriter) {
	sendError(
		w,
		http.StatusBadRequest,
		"Не удалось распарсить тело запроса.",
	)
}

func sendEncodeError(w http.ResponseWriter) {
	sendError(
		w,
		http.StatusInternalServerError,
		"Ошибка при маршалинге ответа.",
	)
}

func sendGetAuthDataError(w http.ResponseWriter) {
	sendError(
		w,
		http.StatusInternalServerError,
		"Не удалось получить данные токена авторизации.",
	)
}

func formatShortName(f string, l string, m string) string {
	if f == "" || l == "" {
		return ""
	}
	name := fmt.Sprintf("%s %v.", f, l[0])
	if m != "" {
		name = fmt.Sprintf("%s %v.", name, m[0])
	}
	return name
}

func toUserRole(s string) (usecase.UserRole, bool) {
	switch s {
	case "root":
		return usecase.Root, true
	case "admin":
		return usecase.Admin, true
	case "teacher":
		return usecase.Teacher, true
	case "student":
		return usecase.Student, true
	default:
		return usecase.Student, false
	}
}

func getParam(r *http.Request, key string) (string, error) {
	p := chi.URLParam(r, key)
	if p == "" {
		return "", fmt.Errorf("параметр {%s} пустой", key)
	}
	return p, nil
}

func getParamInt(r *http.Request, key string) (int, error) {
	p, err := getParam(r, key)
	if err != nil {
		return 0, nil
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func getParamResultID(r *http.Request) (id, error) {
	i, err := getParamInt(r, "result_id")
	return id(i), err
}

func getParamQuizID(r *http.Request) (id, error) {
	i, err := getParamInt(r, "quiz_id")
	return id(i), err
}

func getParamAnswerID(r *http.Request) (id, error) {
	i, err := getParamInt(r, "answer_id")
	return id(i), err
}

func sendParamError(w http.ResponseWriter, what string) {
	sendError(
		w,
		http.StatusBadRequest,
		what,
	)
}
