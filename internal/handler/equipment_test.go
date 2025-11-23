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

func TestEquipmentHandler_Create_Unauthorized(t *testing.T) {
	handler := &EquipmentHandler{}

	body := strings.NewReader(`{"name": "Camera", "category": "Photography"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipment", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

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

func TestEquipmentHandler_Create_InvalidJSON(t *testing.T) {
	handler := &EquipmentHandler{}

	// Add claims to context
	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "owner",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	body := strings.NewReader("invalid json")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipment", body).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "INVALID_JSON" {
		t.Error("expected INVALID_JSON error code")
	}
}

func TestEquipmentHandler_GetByID_InvalidID(t *testing.T) {
	handler := &EquipmentHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipment/invalid-uuid", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "INVALID_ID" {
		t.Error("expected INVALID_ID error code")
	}
}

func TestEquipmentHandler_Update_Unauthorized(t *testing.T) {
	handler := &EquipmentHandler{}

	body := strings.NewReader(`{"name": "Updated Camera"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipment/"+uuid.New().String(), body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEquipmentHandler_Update_InvalidJSON(t *testing.T) {
	handler := &EquipmentHandler{}

	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "owner",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	body := strings.NewReader("invalid json")
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipment/"+uuid.New().String(), body).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestEquipmentHandler_Delete_Unauthorized(t *testing.T) {
	handler := &EquipmentHandler{}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/equipment/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEquipmentHandler_UploadPhoto_Unauthorized(t *testing.T) {
	handler := &EquipmentHandler{}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipment/"+uuid.New().String()+"/photos", nil)
	w := httptest.NewRecorder()

	handler.UploadPhoto(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEquipmentHandler_GetAvailability_InvalidID(t *testing.T) {
	handler := &EquipmentHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipment/invalid-uuid/availability", nil)
	w := httptest.NewRecorder()

	handler.GetAvailability(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
