package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/jwt"
	"github.com/google/uuid"
)

func TestNotificationHandler_List_Unauthorized(t *testing.T) {
	handler := &NotificationHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "UNAUTHORIZED" {
		t.Error("expected UNAUTHORIZED error code")
	}
}

func TestNotificationHandler_GetUnreadCount_Unauthorized(t *testing.T) {
	handler := &NotificationHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count", nil)
	w := httptest.NewRecorder()

	handler.GetUnreadCount(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNotificationHandler_MarkAsRead_Unauthorized(t *testing.T) {
	handler := &NotificationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/"+uuid.New().String()+"/read", nil)
	w := httptest.NewRecorder()

	handler.MarkAsRead(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNotificationHandler_MarkAsRead_InvalidID(t *testing.T) {
	handler := &NotificationHandler{}

	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/invalid-uuid/read", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.MarkAsRead(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "INVALID_ID" {
		t.Error("expected INVALID_ID error code")
	}
}

func TestNotificationHandler_MarkAllAsRead_Unauthorized(t *testing.T) {
	handler := &NotificationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/read-all", nil)
	w := httptest.NewRecorder()

	handler.MarkAllAsRead(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNotificationHandler_Delete_Unauthorized(t *testing.T) {
	handler := &NotificationHandler{}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notifications/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNotificationHandler_Delete_InvalidID(t *testing.T) {
	handler := &NotificationHandler{}

	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notifications/invalid-uuid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
