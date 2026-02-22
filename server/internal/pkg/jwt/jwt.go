package jwt

import (
	"server/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TypeAccess  TokenType = "access"
	TypeRefresh TokenType = "refresh"
)

type Claims struct {
	Subject        uuid.UUID       `json:"sub"`
	UserRole       domain.UserRole `json:"role"`
	Issuer         string          `json:"iss"`
	ExpirationTime time.Time       `json:"exp"`
	IssuedAt       time.Time       `json:"iat"`
	JWTID          uuid.UUID       `json:"jti"`
	Type           TokenType       `json:"type"`
}

func New(subject uuid.UUID, role domain.UserRole, issuer string, duration time.Duration, tType TokenType) *Claims {
	now := time.Now().UTC()
	return &Claims{
		Subject:        subject,
		UserRole:       role,
		Issuer:         issuer,
		ExpirationTime: now.Add(duration),
		IssuedAt:       now,
		JWTID:          uuid.New(),
		Type:           tType,
	}
}

func (j *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(j.ExpirationTime), nil
}

func (j *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(j.IssuedAt), nil
}

func (j *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (j *Claims) GetIssuer() (string, error) {
	return j.Issuer, nil
}

func (j *Claims) GetSubject() (string, error) {
	return j.Subject.String(), nil
}

func (j *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (j *Claims) ID() uuid.UUID {
	return j.Subject
}

func (j *Claims) Role() domain.UserRole {
	return j.UserRole
}
