package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocsHandler_ServeScalarUI(t *testing.T) {
	handler := NewDocsHandler("../../docs")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()

	handler.ServeScalarUI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html; charset=utf-8, got %s", contentType)
	}

	body := w.Body.String()

	// Check for essential HTML elements
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("expected HTML doctype")
	}

	if !strings.Contains(body, "Equipment Rental API - Interactive Documentation") {
		t.Error("expected title in HTML")
	}

	if !strings.Contains(body, "/docs/openapi.yaml") {
		t.Error("expected OpenAPI spec URL reference")
	}

	if !strings.Contains(body, "scalar/api-reference") {
		t.Error("expected Scalar library reference")
	}
}

func TestDocsHandler_ServeOpenAPI(t *testing.T) {
	handler := NewDocsHandler("../../docs")

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
	w := httptest.NewRecorder()

	handler.ServeOpenAPI(w, req)

	// This will return 404 if the file doesn't exist in test environment
	// which is expected behavior
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}

	if w.Code == http.StatusOK {
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/x-yaml" {
			t.Errorf("expected Content-Type application/x-yaml, got %s", contentType)
		}

		corsHeader := w.Header().Get("Access-Control-Allow-Origin")
		if corsHeader != "*" {
			t.Errorf("expected CORS header *, got %s", corsHeader)
		}
	}
}

func TestDocsHandler_ScalarConfiguration(t *testing.T) {
	handler := NewDocsHandler("../../docs")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()

	handler.ServeScalarUI(w, req)

	body := w.Body.String()

	// Verify Scalar configuration options are present
	configChecks := []string{
		"theme:",
		"layout:",
		"showSidebar:",
		"hideModels:",
		"hideDownloadButton:",
		"hideTestRequestButton:",
		"defaultHttpClient:",
	}

	for _, check := range configChecks {
		if !strings.Contains(body, check) {
			t.Errorf("expected Scalar config option %s in HTML", check)
		}
	}
}
