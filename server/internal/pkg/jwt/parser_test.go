package jwt

import (
	"server/internal/domain"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var (
		id        = uuid.New()
		role      = domain.RoleAdmin
		issuer    = "issuer"
		secretKey = []byte("secret")
		duration  = time.Hour
	)
	generator := NewGenerator(issuer, secretKey, duration, duration)
	token, err := generator.GenerateAccess(id, role)

	assert.NoError(t, err)

	parser := NewParser(secretKey)
	claims, err := parser.Parse(token)

	assert.NoError(t, err)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, id, claims.Subject)
	assert.Equal(t, role, claims.UserRole)
	assert.Equal(t, claims.IssuedAt.Add(duration), claims.ExpirationTime)
}

func TestInvalidSecret(t *testing.T) {
	var (
		id             = uuid.New()
		role           = domain.RoleAdmin
		issuer         = "issuer"
		genSecretKey   = []byte("secret")
		parseSecretKey = []byte("another")
		duration       = time.Hour
	)

	generator := NewGenerator(issuer, genSecretKey, duration, duration)
	token, err := generator.GenerateAccess(id, role)
	assert.NoError(t, err)

	parser := NewParser(parseSecretKey)
	_, err = parser.Parse(token)
	assert.Error(t, err)
}

func TestExpiredToken(t *testing.T) {
	var (
		id        = uuid.New()
		role      = domain.RoleAdmin
		issuer    = "issuer"
		secretKey = []byte("secret")
		duration  = time.Duration(0)
	)

	generator := NewGenerator(issuer, secretKey, duration, duration)
	token, err := generator.GenerateAccess(id, role)
	assert.NoError(t, err)

	parser := NewParser(secretKey)
	_, err = parser.Parse(token)
	assert.ErrorAs(t, err, &jwt.ErrTokenExpired)
}
