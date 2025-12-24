package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type ctxKey string

const authKey ctxKey = "AuthData"

type AuthData struct {
	ID   uint64
	Role string
}

func (d *AuthData) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, authKey, *d)
}

func NewAuthDataFromCtx(ctx context.Context) (AuthData, bool) {
	auth, ok := ctx.Value(authKey).(AuthData)
	return auth, ok
}

func DecodeJSON(r io.ReadCloser, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

func EncodeJSON(w io.Writer, v any) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, statusCode int, what string) {
	body := map[string]any{
		"error": what,
	}
	w.WriteHeader(statusCode)
	EncodeJSON(w, body)
}

func SendAuthDataError(w http.ResponseWriter) {
	SendError(
		w,
		http.StatusInternalServerError,
		"Не удалось получить данные токена авторизации.",
	)
}
