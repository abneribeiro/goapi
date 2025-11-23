package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/jwt"
	"github.com/google/uuid"
)

func TestUserHandler_GetMe_Unauthorized(t *testing.T) {
	handler := &UserHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Error == nil || response.Error.Code != "UNAUTHORIZED" {
		t.Error("expected UNAUTHORIZED error code")
	}
}

func TestUserHandler_UpdateMe_Unauthorized(t *testing.T) {
	handler := &UserHandler{}

	body := strings.NewReader(`{"name": "New Name"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserHandler_UpdateMe_InvalidJSON(t *testing.T) {
	handler := &UserHandler{}

	// Add claims to context
	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	body := strings.NewReader("invalid json")
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", body).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "INVALID_JSON" {
		t.Error("expected INVALID_JSON error code")
	}
}
