package jwt

import (
	"server/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		subject  = uuid.New()
		role     = domain.RoleAdmin
		issuer   = "issuer"
		duration = time.Hour
		tType    = TypeAccess
	)

	jwt := New(subject, role, issuer, duration, tType)

	assert.Equal(t, subject, jwt.Subject)
	assert.Equal(t, role, jwt.UserRole)
	assert.Equal(t, issuer, jwt.Issuer)
	assert.Equal(t, tType, jwt.Type)
	assert.Equal(t, jwt.IssuedAt.Add(duration), jwt.ExpirationTime)

	assert.Equal(t, subject, jwt.ID())
	assert.Equal(t, role, jwt.Role())
}
