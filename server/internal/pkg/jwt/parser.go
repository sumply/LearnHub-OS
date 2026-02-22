package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Parser struct {
	secretKey []byte
}

func NewParser(secret []byte) *Parser {
	return &Parser{
		secretKey: secret,
	}
}

func (p *Parser) Parse(tokenString string) (*Claims, error) {
	claims := new(Claims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("algorithm is invalid")
		}
		if t.Header["typ"] != "JWT" {
			return nil, fmt.Errorf("typ is invalid")
		}
		return p.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
