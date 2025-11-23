package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/middleware"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/service"
)

type EquipmentHandler struct {
	equipmentService *service.EquipmentService
}

func NewEquipmentHandler(equipmentService *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{
		equipmentService: equipmentService,
	}
}

func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	var req model.CreateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	equipment, err := h.equipmentService.Create(r.Context(), claims.UserID, &req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondJSON(w, http.StatusBadRequest, model.ErrorResponse("VALIDATION_ERROR", validationErrors.Error()))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to create equipment"))
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse(equipment))
}

func (h *EquipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/equipment/")
	idStr = strings.Split(idStr, "/")[0]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid equipment ID"))
		return
	}

	equipment, err := h.equipmentService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get equipment"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(equipment))
}

func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := pagination.FromRequest(r)
	query := r.URL.Query()

	filter := &model.EquipmentFilter{
		Category: query.Get("category"),
		Location: query.Get("location"),
	}

	if availableStr := query.Get("available"); availableStr != "" {
		available := availableStr == "true"
		filter.Available = &available
	}

	equipment, total, err := h.equipmentService.List(r.Context(), filter, pag)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to list equipment"))
		return
	}

	meta := &model.Meta{
		Page:       pag.Page,
		PerPage:    pag.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total, pag.PerPage),
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(equipment, meta))
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/equipment/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid equipment ID"))
		return
	}

	var req model.UpdateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_JSON", "Invalid request body"))
		return
	}

	equipment, err := h.equipmentService.Update(r.Context(), id, claims.UserID, &req)
	if err != nil {
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not the owner of this equipment"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to update equipment"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(equipment))
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/equipment/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid equipment ID"))
		return
	}

	err = h.equipmentService.Delete(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not the owner of this equipment"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to delete equipment"))
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func (h *EquipmentHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, model.ErrorResponse("UNAUTHORIZED", "User not authenticated"))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/equipment/")
	idStr = strings.TrimSuffix(idStr, "/photos")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid equipment ID"))
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_FORM", "Invalid form data"))
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("MISSING_FILE", "Photo file is required"))
		return
	}
	defer file.Close()

	isPrimary := r.FormValue("is_primary") == "true"

	photo, err := h.equipmentService.AddPhoto(r.Context(), id, claims.UserID, file, header.Filename, isPrimary)
	if err != nil {
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			respondJSON(w, http.StatusForbidden, model.ErrorResponse("FORBIDDEN", "Not the owner of this equipment"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to upload photo"))
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse(photo))
}

func (h *EquipmentHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/equipment/")
	idStr = strings.TrimSuffix(idStr, "/availability")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse("INVALID_ID", "Invalid equipment ID"))
		return
	}

	query := r.URL.Query()
	startDateStr := query.Get("start_date")
	endDateStr := query.Get("end_date")

	startDate := time.Now()
	endDate := time.Now().AddDate(0, 1, 0)

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	availability, err := h.equipmentService.GetAvailability(r.Context(), id, startDate, endDate)
	if err != nil {
		if errors.Is(err, service.ErrEquipmentNotFound) {
			respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "Equipment not found"))
			return
		}
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get availability"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(availability))
}

func (h *EquipmentHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.equipmentService.GetCategories(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to get categories"))
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(categories))
}

func (h *EquipmentHandler) Search(w http.ResponseWriter, r *http.Request) {
	pag := pagination.FromRequest(r)
	query := r.URL.Query().Get("q")

	equipment, total, err := h.equipmentService.Search(r.Context(), query, pag)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse("INTERNAL_ERROR", "Failed to search equipment"))
		return
	}

	meta := &model.Meta{
		Page:       pag.Page,
		PerPage:    pag.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total, pag.PerPage),
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(equipment, meta))
}
