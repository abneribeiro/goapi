package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abneribeiro/goapi/internal/model"
)

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		data           interface{}
		expectedStatus int
	}{
		{
			name:           "success response",
			status:         http.StatusOK,
			data:           model.SuccessResponse(map[string]string{"message": "hello"}),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error response",
			status:         http.StatusBadRequest,
			data:           model.ErrorResponse("ERROR", "something went wrong"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "no content",
			status:         http.StatusNoContent,
			data:           nil,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondJSON(w, tt.status, tt.data)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.data != nil {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}
			}
		})
	}
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	handler := &AuthHandler{}

	body := bytes.NewBufferString("invalid json")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Error == nil || response.Error.Code != "INVALID_JSON" {
		t.Error("expected INVALID_JSON error code")
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	handler := &AuthHandler{}

	body := bytes.NewBufferString("invalid json")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response model.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Error == nil || response.Error.Code != "INVALID_JSON" {
		t.Error("expected INVALID_JSON error code")
	}
}
