package jwt

import (
	"server/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Generator struct {
	iss        string
	secretKey  []byte
	accessDur  time.Duration
	refreshDur time.Duration
}

func NewGenerator(iss string, secretKey []byte, accessDur, refreshDur time.Duration) *Generator {
	return &Generator{
		iss:        iss,
		secretKey:  secretKey,
		accessDur:  accessDur,
		refreshDur: refreshDur,
	}
}

func (g *Generator) GenerateAccess(id uuid.UUID, role domain.UserRole) (string, error) {
	claims := New(id, role, g.iss, g.accessDur, TypeAccess)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(g.secretKey)
}

func (g *Generator) GenerateRefresh(id uuid.UUID, role domain.UserRole) (string, error) {
	claims := New(id, role, g.iss, g.refreshDur, TypeRefresh)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(g.secretKey)
}
