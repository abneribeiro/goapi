package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
	"github.com/abneribeiro/goapi/internal/repository"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) List(ctx context.Context, userID uuid.UUID, pag pagination.Params) ([]*model.Notification, int64, error) {
	filter := &model.NotificationFilter{
		UserID: &userID,
	}
	return s.notificationRepo.List(ctx, filter, pag)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.notificationRepo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	if notification.UserID != userID {
		return ErrNotAuthorized
	}

	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	if notification.UserID != userID {
		return ErrNotAuthorized
	}

	return s.notificationRepo.Delete(ctx, id)
}
