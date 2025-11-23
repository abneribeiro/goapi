package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/service"
)

type ReservationHandler struct {
	reservationService *service.ReservationService
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	var req model.CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	reservation, err := h.reservationService.Create(r.Context(), claims.UserID, &req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("VALIDATION_ERROR", validationErrors.Error()))
			return
		}
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		if errors.Is(err, service.ErrEquipmentUnavailable) {
			respondJSON(w, http.StatusConflict, model.ErrorResponse("UNAVAILABLE", "Equipment not available for selected dates"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to create reservation"))
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse(reservation))
}

func (h *ReservationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/reservations/")
	idStr = strings.Split(idStr, "/")[0]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid reservation ID"))
		return
	}

	reservation, err := h.reservationService.GetByID(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Reservation not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized to view this reservation"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get reservation"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(reservation))
}

func (h *ReservationHandler) ListMyReservations(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	pag := pagination.FromRequest(r)

	reservations, total, err := h.reservationService.ListMyReservations(r.Context(), claims.UserID, pag)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to list reservations"))
		return
	}

	meta := &model.Meta{
		Page:       pag.Page,
		PerPage:    pag.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total, pag.PerPage),
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(reservations, meta))
}

func (h *ReservationHandler) ListOwnerReservations(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	pag := pagination.FromRequest(r)

	reservations, total, err := h.reservationService.ListOwnerReservations(r.Context(), claims.UserID, pag)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to list reservations"))
		return
	}

	meta := &model.Meta{
		Page:       pag.Page,
		PerPage:    pag.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total, pag.PerPage),
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(reservations, meta))
}

func (h *ReservationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/reservations/")
	idStr = strings.TrimSuffix(idStr, "/approve")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid reservation ID"))
		return
	}

	reservation, err := h.reservationService.Approve(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Reservation not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized to approve this reservation"))
			return
		}
		if errors.Is(err, service.ErrReservationNotPending) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_STATUS", "Reservation is not pending"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to approve reservation"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(reservation))
}

func (h *ReservationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/reservations/")
	idStr = strings.TrimSuffix(idStr, "/reject")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid reservation ID"))
		return
	}

	var req model.CancelReservationRequest
	json.NewDecoder(r.Body).Decode(&req)

	reservation, err := h.reservationService.Reject(r.Context(), id, claims.UserID, req.Reason)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Reservation not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized to reject this reservation"))
			return
		}
		if errors.Is(err, service.ErrReservationNotPending) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_STATUS", "Reservation is not pending"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to reject reservation"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(reservation))
}

func (h *ReservationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/reservations/")
	idStr = strings.TrimSuffix(idStr, "/cancel")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid reservation ID"))
		return
	}

	var req model.CancelReservationRequest
	json.NewDecoder(r.Body).Decode(&req)

	reservation, err := h.reservationService.Cancel(r.Context(), id, claims.UserID, req.Reason)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Reservation not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized to cancel this reservation"))
			return
		}
		if errors.Is(err, service.ErrCannotCancel) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("CANNOT_CANCEL", "Cannot cancel this reservation"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to cancel reservation"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(reservation))
}

func (h *ReservationHandler) Complete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/reservations/")
	idStr = strings.TrimSuffix(idStr, "/complete")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid reservation ID"))
		return
	}

	reservation, err := h.reservationService.Complete(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Reservation not found"))
			return
		}
		if errors.Is(err, service.ErrNotAuthorized) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not authorized to complete this reservation"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to complete reservation"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(reservation))
}
