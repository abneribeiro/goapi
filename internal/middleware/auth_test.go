package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/pkg/jwt"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret", time.Hour)
	middleware := NewAuthMiddleware(jwtManager)

	userID := uuid.New()
	token, _ := jwtManager.Generate(userID, "test@example.com", "renter")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		if claims == nil {
			t.Error("expected claims in context")
			return
		}
		if claims.UserID != userID {
			t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret", time.Hour)
	middleware := NewAuthMiddleware(jwtManager)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret", time.Hour)
	middleware := NewAuthMiddleware(jwtManager)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret", time.Hour)
	middleware := NewAuthMiddleware(jwtManager)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	w := httptest.NewRecorder()
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret", -time.Hour)
	middleware := NewAuthMiddleware(jwt.NewManager("test-secret", time.Hour))

	token, _ := jwtManager.Generate(uuid.New(), "test@example.com", "renter")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetUserFromContext_NoUser(t *testing.T) {
	ctx := context.Background()
	claims := GetUserFromContext(ctx)

	if claims != nil {
		t.Error("expected nil claims for context without user")
	}
}

func TestGetUserFromContext_WithUser(t *testing.T) {
	expectedClaims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}

	ctx := context.WithValue(context.Background(), UserContextKey, expectedClaims)
	claims := GetUserFromContext(ctx)

	if claims == nil {
		t.Fatal("expected claims in context")
	}

	if claims.UserID != expectedClaims.UserID {
		t.Errorf("expected user ID %s, got %s", expectedClaims.UserID, claims.UserID)
	}
}
