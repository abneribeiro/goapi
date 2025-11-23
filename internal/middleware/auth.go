package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/jwt"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
	jwtManager *jwt.Manager
}

func NewAuthMiddleware(jwtManager *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.respondUnauthorized(w, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.respondUnauthorized(w, "invalid authorization header format")
			return
		}

		claims, err := m.jwtManager.Validate(parts[1])
		if err != nil {
			if err == jwt.ErrExpiredToken {
				m.respondUnauthorized(w, "token has expired")
				return
			}
			m.respondUnauthorized(w, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.ErrorResponse("UNAUTHORIZED", message))
}

func GetUserFromContext(ctx context.Context) *jwt.Claims {
	claims, ok := ctx.Value(UserContextKey).(*jwt.Claims)
	if !ok {
		return nil
	}
	return claims
}
