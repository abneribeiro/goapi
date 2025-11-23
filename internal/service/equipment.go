package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/repository"
)

var (
	ErrEquipmentNotFound = errors.New("equipment not found")
	ErrNotOwner          = errors.New("not the owner of this equipment")
	ErrNoAvailability    = errors.New("equipment not available for selected dates")
)

type EquipmentService struct {
	equipmentRepo *repository.EquipmentRepository
	uploadPath    string
}

func NewEquipmentService(equipmentRepo *repository.EquipmentRepository, uploadPath string) *EquipmentService {
	return &EquipmentService{
		equipmentRepo: equipmentRepo,
		uploadPath:    uploadPath,
	}
}

func (s *EquipmentService) Create(ctx context.Context, ownerID uuid.UUID, req *model.CreateEquipmentRequest) (*model.Equipment, error) {
	v := validator.New()
	v.Required("name", req.Name)
	v.Required("category", req.Category)

	if req.PricePerHour == nil && req.PricePerDay == nil && req.PricePerWeek == nil {
		v.AddError("price", "at least one price must be set")
	}

	if v.Errors().HasErrors() {
		return nil, v.Errors()
	}

	equipment := &model.Equipment{
		OwnerID:      ownerID,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		PricePerHour: req.PricePerHour,
		PricePerDay:  req.PricePerDay,
		PricePerWeek: req.PricePerWeek,
		Location:     req.Location,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		AutoApprove:  req.AutoApprove,
	}

	if err := s.equipmentRepo.Create(ctx, equipment); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *EquipmentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Equipment, error) {
	equipment, err := s.equipmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}
	return equipment, nil
}

func (s *EquipmentService) List(ctx context.Context, filter *model.EquipmentFilter, pag pagination.Params) ([]*model.Equipment, int64, error) {
	return s.equipmentRepo.List(ctx, filter, pag)
}

func (s *EquipmentService) Update(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, req *model.UpdateEquipmentRequest) (*model.Equipment, error) {
	equipment, err := s.equipmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	if equipment.OwnerID != ownerID {
		return nil, ErrNotOwner
	}

	if req.Name != "" {
		equipment.Name = req.Name
	}
	if req.Description != "" {
		equipment.Description = req.Description
	}
	if req.Category != "" {
		equipment.Category = req.Category
	}
	if req.PricePerHour != nil {
		equipment.PricePerHour = req.PricePerHour
	}
	if req.PricePerDay != nil {
		equipment.PricePerDay = req.PricePerDay
	}
	if req.PricePerWeek != nil {
		equipment.PricePerWeek = req.PricePerWeek
	}
	if req.Location != "" {
		equipment.Location = req.Location
	}
	if req.Latitude != nil {
		equipment.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		equipment.Longitude = req.Longitude
	}
	if req.Available != nil {
		equipment.Available = *req.Available
	}
	if req.AutoApprove != nil {
		equipment.AutoApprove = *req.AutoApprove
	}

	if err := s.equipmentRepo.Update(ctx, equipment); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *EquipmentService) Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	equipment, err := s.equipmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return ErrEquipmentNotFound
		}
		return err
	}

	if equipment.OwnerID != ownerID {
		return ErrNotOwner
	}

	return s.equipmentRepo.Delete(ctx, id)
}

func (s *EquipmentService) AddPhoto(ctx context.Context, equipmentID uuid.UUID, ownerID uuid.UUID, file io.Reader, filename string, isPrimary bool) (*model.EquipmentPhoto, error) {
	equipment, err := s.equipmentRepo.GetByID(ctx, equipmentID)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	if equipment.OwnerID != ownerID {
		return nil, ErrNotOwner
	}

	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	photo := &model.EquipmentPhoto{
		EquipmentID: equipmentID,
		URL:         "/uploads/" + newFilename,
		IsPrimary:   isPrimary,
	}

	if err := s.equipmentRepo.AddPhoto(ctx, photo); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return photo, nil
}

func (s *EquipmentService) GetAvailability(ctx context.Context, equipmentID uuid.UUID, startDate, endDate time.Time) ([]model.EquipmentAvailability, error) {
	_, err := s.equipmentRepo.GetByID(ctx, equipmentID)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	return s.equipmentRepo.GetAvailabilityCalendar(ctx, equipmentID, startDate, endDate)
}

func (s *EquipmentService) GetCategories(ctx context.Context) ([]string, error) {
	return s.equipmentRepo.GetCategories(ctx)
}

func (s *EquipmentService) Search(ctx context.Context, query string, pag pagination.Params) ([]*model.Equipment, int64, error) {
	return s.equipmentRepo.Search(ctx, query, pag)
}
