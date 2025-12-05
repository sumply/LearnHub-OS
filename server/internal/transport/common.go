package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
