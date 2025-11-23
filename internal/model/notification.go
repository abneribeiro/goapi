package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationReservationCreated   NotificationType = "reservation_created"
	NotificationReservationApproved  NotificationType = "reservation_approved"
	NotificationReservationRejected  NotificationType = "reservation_rejected"
	NotificationReservationCancelled NotificationType = "reservation_cancelled"
	NotificationReservationCompleted NotificationType = "reservation_completed"
	NotificationReservationReminder  NotificationType = "reservation_reminder"
	NotificationEquipmentReturned    NotificationType = "equipment_returned"
	NotificationPaymentReceived      NotificationType = "payment_received"
)

type Notification struct {
	ID            uuid.UUID        `json:"id"`
	UserID        uuid.UUID        `json:"user_id"`
	Type          NotificationType `json:"type"`
	Title         string           `json:"title"`
	Message       string           `json:"message"`
	Read          bool             `json:"read"`
	ReferenceID   *uuid.UUID       `json:"reference_id,omitempty"`
	ReferenceType string           `json:"reference_type,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
}

type NotificationFilter struct {
	UserID *uuid.UUID
	Read   *bool
	Type   *NotificationType
}

type CreateNotificationRequest struct {
	UserID        uuid.UUID        `json:"user_id"`
	Type          NotificationType `json:"type"`
	Title         string           `json:"title"`
	Message       string           `json:"message"`
	ReferenceID   *uuid.UUID       `json:"reference_id,omitempty"`
	ReferenceType string           `json:"reference_type,omitempty"`
}
