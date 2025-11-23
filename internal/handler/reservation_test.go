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

func TestReservationHandler_Create_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	body := strings.NewReader(`{"equipment_id": "` + uuid.New().String() + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != "UNAUTHORIZED" {
		t.Error("expected UNAUTHORIZED error code")
	}
}

func TestReservationHandler_Create_InvalidJSON(t *testing.T) {
	handler := &ReservationHandler{}

	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	body := strings.NewReader("invalid json")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", body).WithContext(ctx)
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

func TestReservationHandler_GetByID_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_GetByID_InvalidID(t *testing.T) {
	handler := &ReservationHandler{}

	claims := &jwt.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "renter",
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/invalid-uuid", nil).WithContext(ctx)
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

func TestReservationHandler_ListMyReservations_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	w := httptest.NewRecorder()

	handler.ListMyReservations(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_ListOwnerReservations_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/owner", nil)
	w := httptest.NewRecorder()

	handler.ListOwnerReservations(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_Approve_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/"+uuid.New().String()+"/approve", nil)
	w := httptest.NewRecorder()

	handler.Approve(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_Reject_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/"+uuid.New().String()+"/reject", nil)
	w := httptest.NewRecorder()

	handler.Reject(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_Cancel_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/"+uuid.New().String()+"/cancel", nil)
	w := httptest.NewRecorder()

	handler.Cancel(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestReservationHandler_Complete_Unauthorized(t *testing.T) {
	handler := &ReservationHandler{}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/"+uuid.New().String()+"/complete", nil)
	w := httptest.NewRecorder()

	handler.Complete(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
