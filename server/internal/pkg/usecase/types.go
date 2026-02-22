package usecase

import (
	"context"
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
