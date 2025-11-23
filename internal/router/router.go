package router

import (
	"net/http"

	"github.com/abneribeiro/goapi/internal/handler"
	"github.com/abneribeiro/goapi/internal/middleware"
)

type Router struct {
	mux            *http.ServeMux
	authMiddleware *middleware.AuthMiddleware
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	equipHandler   *handler.EquipmentHandler
	resHandler     *handler.ReservationHandler
	notifHandler   *handler.NotificationHandler
	docsHandler    *handler.DocsHandler
}

func New(
	authMiddleware *middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	equipHandler *handler.EquipmentHandler,
	resHandler *handler.ReservationHandler,
	notifHandler *handler.NotificationHandler,
	docsHandler *handler.DocsHandler,
) *Router {
	return &Router{
		mux:            http.NewServeMux(),
		authMiddleware: authMiddleware,
		authHandler:    authHandler,
		userHandler:    userHandler,
		equipHandler:   equipHandler,
		resHandler:     resHandler,
		notifHandler:   notifHandler,
		docsHandler:    docsHandler,
	}
}

func (r *Router) Setup() http.Handler {
	r.mux.HandleFunc("GET /health", r.healthCheck)

	// Documentation routes - support both with and without trailing slash
	r.mux.HandleFunc("GET /docs", r.docsHandler.ServeScalarUI)
	r.mux.HandleFunc("GET /docs/", r.docsHandler.ServeScalarUI)
	r.mux.HandleFunc("GET /docs/openapi.yaml", r.docsHandler.ServeOpenAPI)

	r.mux.HandleFunc("POST /api/v1/auth/register", r.authHandler.Register)
	r.mux.HandleFunc("POST /api/v1/auth/login", r.authHandler.Login)

	r.mux.Handle("GET /api/v1/users/me", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.GetMe)))
	r.mux.Handle("PUT /api/v1/users/me", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.UpdateMe)))

	r.mux.HandleFunc("GET /api/v1/equipment", r.equipHandler.List)
	r.mux.HandleFunc("GET /api/v1/equipment/search", r.equipHandler.Search)
	r.mux.HandleFunc("GET /api/v1/equipment/categories", r.equipHandler.GetCategories)
	r.mux.HandleFunc("GET /api/v1/equipment/{id}", r.equipHandler.GetByID)
	r.mux.HandleFunc("GET /api/v1/equipment/{id}/availability", r.equipHandler.GetAvailability)
	r.mux.Handle("POST /api/v1/equipment", r.authMiddleware.Authenticate(http.HandlerFunc(r.equipHandler.Create)))
	r.mux.Handle("PUT /api/v1/equipment/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.equipHandler.Update)))
	r.mux.Handle("DELETE /api/v1/equipment/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.equipHandler.Delete)))
	r.mux.Handle("POST /api/v1/equipment/{id}/photos", r.authMiddleware.Authenticate(http.HandlerFunc(r.equipHandler.UploadPhoto)))

	r.mux.Handle("GET /api/v1/reservations", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.ListMyReservations)))
	r.mux.Handle("GET /api/v1/reservations/owner", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.ListOwnerReservations)))
	r.mux.Handle("GET /api/v1/reservations/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.GetByID)))
	r.mux.Handle("POST /api/v1/reservations", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.Create)))
	r.mux.Handle("PUT /api/v1/reservations/{id}/approve", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.Approve)))
	r.mux.Handle("PUT /api/v1/reservations/{id}/reject", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.Reject)))
	r.mux.Handle("PUT /api/v1/reservations/{id}/cancel", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.Cancel)))
	r.mux.Handle("PUT /api/v1/reservations/{id}/complete", r.authMiddleware.Authenticate(http.HandlerFunc(r.resHandler.Complete)))

	r.mux.Handle("GET /api/v1/notifications", r.authMiddleware.Authenticate(http.HandlerFunc(r.notifHandler.List)))
	r.mux.Handle("GET /api/v1/notifications/unread-count", r.authMiddleware.Authenticate(http.HandlerFunc(r.notifHandler.GetUnreadCount)))
	r.mux.Handle("PUT /api/v1/notifications/{id}/read", r.authMiddleware.Authenticate(http.HandlerFunc(r.notifHandler.MarkAsRead)))
	r.mux.Handle("PUT /api/v1/notifications/read-all", r.authMiddleware.Authenticate(http.HandlerFunc(r.notifHandler.MarkAllAsRead)))
	r.mux.Handle("DELETE /api/v1/notifications/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.notifHandler.Delete)))

	fs := http.FileServer(http.Dir("./uploads"))
	r.mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fs))

	return middleware.CORS(middleware.Logger(middleware.Recovery(r.mux)))
}

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
