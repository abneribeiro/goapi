package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	user, err := h.userService.GetByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "User not found"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get user"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(user))
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	user, err := h.userService.Update(r.Context(), claims.UserID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "User not found"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to update user"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(user))
}
