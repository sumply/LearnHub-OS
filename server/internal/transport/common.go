package transport

import (
	"context"
	"encoding/json"
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
