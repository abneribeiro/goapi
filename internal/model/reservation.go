package model

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusApproved  ReservationStatus = "approved"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
	StatusRejected  ReservationStatus = "rejected"
)

type Reservation struct {
	ID                 uuid.UUID         `json:"id"`
	EquipmentID        uuid.UUID         `json:"equipment_id"`
	Equipment          *Equipment        `json:"equipment,omitempty"`
	RenterID           uuid.UUID         `json:"renter_id"`
	Renter             *User             `json:"renter,omitempty"`
	StartDate          time.Time         `json:"start_date"`
	EndDate            time.Time         `json:"end_date"`
	Status             ReservationStatus `json:"status"`
	TotalPrice         float64           `json:"total_price"`
	CancellationReason string            `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type CreateReservationRequest struct {
	EquipmentID uuid.UUID `json:"equipment_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type ReservationFilter struct {
	EquipmentID *uuid.UUID
	RenterID    *uuid.UUID
	OwnerID     *uuid.UUID
	Status      *ReservationStatus
	StartDate   *time.Time
	EndDate     *time.Time
}

type CancelReservationRequest struct {
	Reason string `json:"reason,omitempty"`
}
