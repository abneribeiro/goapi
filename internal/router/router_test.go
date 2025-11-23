package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abneribeiro/goapi/internal/handler"
	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/pkg/jwt"
)

func setupTestRouter() *Router {
	jwtManager := jwt.NewManager("test-secret", 24)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Create handlers with nil services for testing routes only
	authHandler := &handler.AuthHandler{}
	userHandler := &handler.UserHandler{}
	equipHandler := &handler.EquipmentHandler{}
	resHandler := &handler.ReservationHandler{}
	notifHandler := &handler.NotificationHandler{}
	docsHandler := handler.NewDocsHandler("../../docs")

	return New(
		authMiddleware,
		authHandler,
		userHandler,
		equipHandler,
		resHandler,
		notifHandler,
		docsHandler,
	)
}

func TestHealthCheck(t *testing.T) {
	r := setupTestRouter()
	handler := r.Setup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"healthy"}`
	if w.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDocsRoute(t *testing.T) {
	r := setupTestRouter()
	handler := r.Setup()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantType   string
	}{
		{
			name:       "docs without trailing slash",
			path:       "/docs",
			wantStatus: http.StatusOK,
			wantType:   "text/html; charset=utf-8",
		},
		{
			name:       "docs with trailing slash",
			path:       "/docs/",
			wantStatus: http.StatusOK,
			wantType:   "text/html; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.wantType {
				t.Errorf("expected Content-Type %s, got %s", tt.wantType, contentType)
			}
		})
	}
}

func TestProtectedRoutesRequireAuth(t *testing.T) {
	r := setupTestRouter()
	handler := r.Setup()

	protectedRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/users/me"},
		{http.MethodPut, "/api/v1/users/me"},
		{http.MethodPost, "/api/v1/equipment"},
		{http.MethodGet, "/api/v1/reservations"},
		{http.MethodGet, "/api/v1/notifications"},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected status %d for %s %s, got %d", http.StatusUnauthorized, route.method, route.path, w.Code)
			}
		})
	}
}

func TestPublicRoutes(t *testing.T) {
	r := setupTestRouter()
	handler := r.Setup()

	publicRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/health"},
		{http.MethodGet, "/docs"},
		{http.MethodGet, "/api/v1/equipment"},
		{http.MethodGet, "/api/v1/equipment/categories"},
		{http.MethodGet, "/api/v1/equipment/search"},
	}

	for _, route := range publicRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Public routes should not return 401
			if w.Code == http.StatusUnauthorized {
				t.Errorf("expected public route %s %s to not require auth", route.method, route.path)
			}
		})
	}
}
