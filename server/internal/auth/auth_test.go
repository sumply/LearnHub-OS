package auth

import (
	"server/internal/domain"
	"testing"
	"time"
)

func TestNewJWT(t *testing.T) {
	aDur := time.Hour
	rDur := time.Hour + 1
	iss := "iss"
	secretKey := []byte("secret")
	jwt := NewJWT(secretKey, aDur, rDur, iss)
	if string(jwt.secretKey) != string(secretKey) {
		t.Errorf("secretKey: %s; except: %s", jwt.secretKey, secretKey)
	}
	if jwt.accessDuration != aDur {
		t.Errorf("access duration: %v; except: %v", jwt.accessDuration, aDur)
	}
	if jwt.refreshDuration != rDur {
		t.Errorf("refresh duration: %v; except: %v", jwt.refreshDuration, rDur)
	}
	if jwt.iss != iss {
		t.Errorf("iss: %v; except: %v", jwt.iss, iss)
	}
}

func TestGenerateJWTToken(t *testing.T) {
	aDur := time.Hour
	rDur := time.Hour + 1
	iss := "iss"
	secretKey := []byte("secret")
	jwt := NewJWT(secretKey, aDur, rDur, iss)
	_, err := jwt.GenerateTokenPair(&domain.User{ID: 1, Role: 1})
	if err != nil {
		t.Errorf("%v", err.Error())
	}
}

func TestParseJWTToken(t *testing.T) {
	aDur := time.Hour
	rDur := time.Hour + 1
	iss := "iss"
	secretKey := []byte("secret")
	jwt := NewJWT(secretKey, aDur, rDur, iss)
	pair, err := jwt.GenerateTokenPair(&domain.User{ID: 1, Role: 1})
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	identity, err := jwt.Parse(pair.Access)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", identity)
}
