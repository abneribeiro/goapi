package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/repository"
)

var (
	ErrReservationNotFound   = errors.New("reservation not found")
	ErrCannotCancel          = errors.New("cannot cancel this reservation")
	ErrNotAuthorized         = errors.New("not authorized to perform this action")
	ErrEquipmentUnavailable  = errors.New("equipment not available for selected dates")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrReservationNotPending = errors.New("reservation is not in pending status")
)

type ReservationService struct {
	reservationRepo   *repository.ReservationRepository
	equipmentRepo     *repository.EquipmentRepository
	notificationRepo  *repository.NotificationRepository
}

func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	equipmentRepo *repository.EquipmentRepository,
	notificationRepo *repository.NotificationRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo:   reservationRepo,
		equipmentRepo:     equipmentRepo,
		notificationRepo:  notificationRepo,
	}
}

func (s *ReservationService) Create(ctx context.Context, renterID uuid.UUID, req *model.CreateReservationRequest) (*model.Reservation, error) {
	v := validator.New()

	if req.StartDate.IsZero() {
		v.AddError("start_date", "is required")
	}
	if req.EndDate.IsZero() {
		v.AddError("end_date", "is required")
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		if !req.EndDate.After(req.StartDate) {
			v.AddError("end_date", "must be after start_date")
		}
		if req.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
			v.AddError("start_date", "must be in the future")
		}
	}

	if v.Errors().HasErrors() {
		return nil, v.Errors()
	}

	equipment, err := s.equipmentRepo.GetByID(ctx, req.EquipmentID)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	if !equipment.Available {
		return nil, ErrEquipmentUnavailable
	}

	available, err := s.equipmentRepo.CheckAvailability(ctx, req.EquipmentID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrEquipmentUnavailable
	}

	totalPrice := s.calculatePrice(equipment, req.StartDate, req.EndDate)

	status := model.StatusPending
	if equipment.AutoApprove {
		status = model.StatusApproved
	}

	reservation := &model.Reservation{
		EquipmentID: req.EquipmentID,
		RenterID:    renterID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      status,
		TotalPrice:  totalPrice,
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	s.createNotification(ctx, equipment.OwnerID, model.NotificationReservationCreated,
		"New Reservation Request",
		"You have a new reservation request for "+equipment.Name,
		&reservation.ID, "reservation")

	if status == model.StatusApproved {
		s.createNotification(ctx, renterID, model.NotificationReservationApproved,
			"Reservation Approved",
			"Your reservation for "+equipment.Name+" has been automatically approved",
			&reservation.ID, "reservation")
	}

	return reservation, nil
}

func (s *ReservationService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReservationNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if reservation.RenterID != userID && reservation.Equipment.OwnerID != userID {
		return nil, ErrNotAuthorized
	}

	return reservation, nil
}

func (s *ReservationService) ListMyReservations(ctx context.Context, userID uuid.UUID, pag pagination.Params) ([]*model.Reservation, int64, error) {
	filter := &model.ReservationFilter{
		RenterID: &userID,
	}
	return s.reservationRepo.List(ctx, filter, pag)
}

func (s *ReservationService) ListOwnerReservations(ctx context.Context, ownerID uuid.UUID, pag pagination.Params) ([]*model.Reservation, int64, error) {
	filter := &model.ReservationFilter{
		OwnerID: &ownerID,
	}
	return s.reservationRepo.List(ctx, filter, pag)
}

func (s *ReservationService) Approve(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReservationNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if reservation.Equipment.OwnerID != ownerID {
		return nil, ErrNotAuthorized
	}

	if reservation.Status != model.StatusPending {
		return nil, ErrReservationNotPending
	}

	if err := s.reservationRepo.UpdateStatus(ctx, id, model.StatusApproved, ""); err != nil {
		return nil, err
	}

	s.createNotification(ctx, reservation.RenterID, model.NotificationReservationApproved,
		"Reservation Approved",
		"Your reservation for "+reservation.Equipment.Name+" has been approved",
		&id, "reservation")

	reservation.Status = model.StatusApproved
	return reservation, nil
}

func (s *ReservationService) Reject(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, reason string) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReservationNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if reservation.Equipment.OwnerID != ownerID {
		return nil, ErrNotAuthorized
	}

	if reservation.Status != model.StatusPending {
		return nil, ErrReservationNotPending
	}

	if err := s.reservationRepo.UpdateStatus(ctx, id, model.StatusRejected, reason); err != nil {
		return nil, err
	}

	s.createNotification(ctx, reservation.RenterID, model.NotificationReservationRejected,
		"Reservation Rejected",
		"Your reservation for "+reservation.Equipment.Name+" has been rejected",
		&id, "reservation")

	reservation.Status = model.StatusRejected
	reservation.CancellationReason = reason
	return reservation, nil
}

func (s *ReservationService) Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReservationNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if reservation.RenterID != userID && reservation.Equipment.OwnerID != userID {
		return nil, ErrNotAuthorized
	}

	if reservation.Status != model.StatusPending && reservation.Status != model.StatusApproved {
		return nil, ErrCannotCancel
	}

	if reservation.Status == model.StatusApproved && reservation.StartDate.Before(time.Now().Add(24*time.Hour)) {
		return nil, ErrCannotCancel
	}

	if err := s.reservationRepo.UpdateStatus(ctx, id, model.StatusCancelled, reason); err != nil {
		return nil, err
	}

	notifyUserID := reservation.Equipment.OwnerID
	if userID == reservation.Equipment.OwnerID {
		notifyUserID = reservation.RenterID
	}

	s.createNotification(ctx, notifyUserID, model.NotificationReservationCancelled,
		"Reservation Cancelled",
		"A reservation for "+reservation.Equipment.Name+" has been cancelled",
		&id, "reservation")

	reservation.Status = model.StatusCancelled
	reservation.CancellationReason = reason
	return reservation, nil
}

func (s *ReservationService) Complete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReservationNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if reservation.Equipment.OwnerID != ownerID {
		return nil, ErrNotAuthorized
	}

	if reservation.Status != model.StatusApproved {
		return nil, errors.New("can only complete approved reservations")
	}

	if err := s.reservationRepo.UpdateStatus(ctx, id, model.StatusCompleted, ""); err != nil {
		return nil, err
	}

	s.createNotification(ctx, reservation.RenterID, model.NotificationReservationCompleted,
		"Reservation Completed",
		"Your reservation for "+reservation.Equipment.Name+" has been marked as completed",
		&id, "reservation")

	reservation.Status = model.StatusCompleted
	return reservation, nil
}

func (s *ReservationService) calculatePrice(equipment *model.Equipment, startDate, endDate time.Time) float64 {
	duration := endDate.Sub(startDate)
	days := int(math.Ceil(duration.Hours() / 24))

	if equipment.PricePerWeek != nil && days >= 7 {
		weeks := days / 7
		remainingDays := days % 7
		weekPrice := float64(weeks) * *equipment.PricePerWeek

		if remainingDays > 0 && equipment.PricePerDay != nil {
			return weekPrice + float64(remainingDays)**equipment.PricePerDay
		}
		return weekPrice
	}

	if equipment.PricePerDay != nil {
		return float64(days) * *equipment.PricePerDay
	}

	if equipment.PricePerHour != nil {
		hours := int(math.Ceil(duration.Hours()))
		return float64(hours) * *equipment.PricePerHour
	}

	return 0
}

func (s *ReservationService) createNotification(ctx context.Context, userID uuid.UUID, notifType model.NotificationType, title, message string, refID *uuid.UUID, refType string) {
	notification := &model.Notification{
		UserID:        userID,
		Type:          notifType,
		Title:         title,
		Message:       message,
		ReferenceID:   refID,
		ReferenceType: refType,
	}
	s.notificationRepo.Create(ctx, notification)
}
