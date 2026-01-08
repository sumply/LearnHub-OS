package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"server/internal/dto"
)

type ctxKey string

const identityKey ctxKey = "identity"

func ContextWithIdentity(ctx context.Context, identity *dto.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func NewIdentityFromCtx(ctx context.Context) (*dto.Identity, bool) {
	identity, ok := ctx.Value(identityKey).(*dto.Identity)
	return identity, ok
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
