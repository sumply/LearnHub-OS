package usecase

import (
	"context"
	"log/slog"
	"server/internal/domain"

	"github.com/google/uuid"
)

type ctxKey string

const (
	ctxIdentity ctxKey = "identity"
)

type Identity interface {
	ID() uuid.UUID
	Role() domain.UserRole
}

func IdentityWithContext(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, ctxIdentity, identity)
}

func IdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(ctxIdentity).(Identity)
	return identity, ok
}

func IdentityToSlogAttr(identity Identity) slog.Attr {
	return slog.Group("identity",
		slog.String("id", identity.ID().String()),
		slog.String("role", string(identity.Role())),
	)
}
