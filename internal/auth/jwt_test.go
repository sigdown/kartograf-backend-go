package auth

import (
	"testing"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func TestTokenManagerGenerateAndParse(t *testing.T) {
	manager := NewTokenManager("secret", time.Hour)

	token, err := manager.Generate(domain.User{
		ID:       42,
		Username: "alice",
		Role:     domain.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("unexpected user id: %d", claims.UserID)
	}
	if claims.Username != "alice" {
		t.Fatalf("unexpected username: %s", claims.Username)
	}
	if claims.Role != domain.RoleAdmin {
		t.Fatalf("unexpected role: %s", claims.Role)
	}
}

func TestTokenManagerRejectsExpiredToken(t *testing.T) {
	manager := NewTokenManager("secret", -time.Minute)

	token, err := manager.Generate(domain.User{
		ID:       1,
		Username: "bob",
		Role:     domain.RoleUser,
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	if _, err := manager.Parse(token); err == nil {
		t.Fatal("expected expired token error")
	}
}
