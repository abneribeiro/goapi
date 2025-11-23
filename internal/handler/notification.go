package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/service"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	pag := pagination.FromRequest(r)

	notifications, total, err := h.notificationService.List(r.Context(), claims.UserID, pag)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to list notifications"))
		return
	}

	meta := &model.Meta{
		Page:       pag.Page,
		PerPage:    pag.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total, pag.PerPage),
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(notifications, meta))
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	count, err := h.notificationService.GetUnreadCount(r.Context(), claims.UserID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get unread count"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(map[string]int64{"unread_count": count}))
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")
	idStr = strings.TrimSuffix(idStr, "/read")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid notification ID"))
		return
	}

	err = h.notificationService.MarkAsRead(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Notification not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to mark notification as read"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(map[string]string{"message": "Notification marked as read"}))
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	err := h.notificationService.MarkAllAsRead(r.Context(), claims.UserID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to mark all notifications as read"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(map[string]string{"message": "All notifications marked as read"}))
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid notification ID"))
		return
	}

	err = h.notificationService.Delete(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Notification not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to delete notification"))
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
