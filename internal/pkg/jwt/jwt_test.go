package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestManager_GenerateAndValidate(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)
	userID := uuid.New()
	email := "test@example.com"
	role := "renter"

	token, err := manager.Generate(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("token should not be empty")
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestManager_ValidateInvalidToken(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)

	_, err := manager.Validate("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestManager_ValidateExpiredToken(t *testing.T) {
	manager := NewManager("test-secret", -time.Hour)
	userID := uuid.New()

	token, err := manager.Generate(userID, "test@example.com", "renter")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager.Validate(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}

	if err != ErrExpiredToken {
		t.Errorf("expected ErrExpiredToken, got %v", err)
	}
}

func TestManager_ValidateWrongSecret(t *testing.T) {
	manager1 := NewManager("secret-1", time.Hour)
	manager2 := NewManager("secret-2", time.Hour)

	token, err := manager1.Generate(uuid.New(), "test@example.com", "renter")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager2.Validate(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}

	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
