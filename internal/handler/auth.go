package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("VALIDATION_ERROR", validationErrors.Error()))
			return
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			respondJSON(w, http.StatusConflict, model.ErrorResponse("EMAIL_EXISTS", "Email already registered"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to register user"))
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse(resp))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("VALIDATION_ERROR", validationErrors.Error()))
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("INVALID_CREDENTIALS", "Invalid email or password"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to login"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(resp))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
