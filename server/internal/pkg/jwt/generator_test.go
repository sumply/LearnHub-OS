package jwt

import (
	"server/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createGenerator() *Generator {
	var (
		issuer    = "issuer"
		secretKey = []byte("secret")
		aDur      = time.Minute
		rDur      = time.Hour
	)
	return NewGenerator(issuer, secretKey, aDur, rDur)
}

func TestGenerateAccess(t *testing.T) {
	var (
		id   = uuid.New()
		role = domain.RoleAdmin
	)

	generator := createGenerator()
	token, err := generator.GenerateAccess(id, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefresh(t *testing.T) {
	var (
		id   = uuid.New()
		role = domain.RoleAdmin
	)

	generator := createGenerator()
	token, err := generator.GenerateRefresh(id, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
