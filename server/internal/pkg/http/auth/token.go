package auth

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

var ctxKey string = "token"

type Token struct {
	id   uuid.UUID
	role string
}

func NewToken(id uuid.UUID, role domain.UserRole) *Token {
	return &Token{
		id:   id,
		role: string(role),
	}
}

func TokenFromContext(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(ctxKey).(*Token)
	return token, ok
}

func (t *Token) ID() uuid.UUID {
	return t.id
}

func (t *Token) Role() domain.UserRole {
	return domain.UserRole(t.role)
}

func (t *Token) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, t)
}
