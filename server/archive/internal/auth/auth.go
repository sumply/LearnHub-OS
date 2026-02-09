package auth

import (
	"fmt"
	"server/archive/internal/common"
	"server/archive/internal/domain"
	"server/archive/internal/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT struct {
	secretKey       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
	iss             string
}

func NewJWT(secretKey []byte, aDur, rDur time.Duration, iss string) *JWT {
	return &JWT{
		secretKey:       secretKey,
		accessDuration:  aDur,
		refreshDuration: rDur,
		iss:             iss,
	}
}

func (a *JWT) GenerateTokenPair(u *domain.User) (*domain.TokenPair, error) {
	access, err := a.createClaims(u.ID, u.Role, a.accessDuration)
	if err != nil {
		return nil, err
	}
	refresh, err := a.createClaims(u.ID, u.Role, a.refreshDuration)
	if err != nil {
		return nil, err
	}
	return &domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (a *JWT) Parse(s string) (*dto.Identity, error) {
	token, err := jwt.Parse(s, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted signing method: %s", t.Method.Alg())
		}
		return a.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if err := a.validateToken(token); err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims is not map")
	}

	identity, err := a.getIdentity(claims)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (a *JWT) validateToken(t *jwt.Token) error {
	if !t.Valid {
		return fmt.Errorf("token is invalid")
	}
	exp, err := t.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	fmt.Printf("exp: %d\nnow: %d", exp.Time.Unix(), time.Now().Unix())
	if exp.Time.Before(time.Now().UTC()) {
		return jwt.ErrTokenExpired
	}
	iss, err := t.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != a.iss {
		return fmt.Errorf("invalid iss: %s", iss)
	}
	return nil
}

func (a *JWT) getIdentity(c jwt.MapClaims) (*dto.Identity, error) {
	sub, ok := c["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("sub is missing")
	}
	role, ok := c["role"].(float64)
	if !ok {
		return nil, fmt.Errorf("role is missing")
	}
	return &dto.Identity{
		ID:   common.ID(sub),
		Role: domain.UserRole(role),
	}, nil
}

func (a *JWT) createClaims(id common.ID, role domain.UserRole, duration time.Duration) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(duration)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  a.iss,
		"sub":  id,
		"role": role,
		"exp":  exp.Unix(),
		"iat":  now,
		"jti":  uuid.NewString(),
	})
	return claims.SignedString(a.secretKey)
}
