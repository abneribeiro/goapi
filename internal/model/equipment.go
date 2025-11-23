package model

import (
	"time"

	"github.com/google/uuid"
)

type Equipment struct {
	ID           uuid.UUID         `json:"id"`
	OwnerID      uuid.UUID         `json:"owner_id"`
	Owner        *User             `json:"owner,omitempty"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Category     string            `json:"category"`
	PricePerHour *float64          `json:"price_per_hour,omitempty"`
	PricePerDay  *float64          `json:"price_per_day,omitempty"`
	PricePerWeek *float64          `json:"price_per_week,omitempty"`
	Location     string            `json:"location,omitempty"`
	Latitude     *float64          `json:"latitude,omitempty"`
	Longitude    *float64          `json:"longitude,omitempty"`
	Available    bool              `json:"available"`
	AutoApprove  bool              `json:"auto_approve"`
	Photos       []EquipmentPhoto  `json:"photos,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type EquipmentPhoto struct {
	ID          uuid.UUID `json:"id"`
	EquipmentID uuid.UUID `json:"equipment_id"`
	URL         string    `json:"url"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateEquipmentRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	Category     string   `json:"category"`
	PricePerHour *float64 `json:"price_per_hour,omitempty"`
	PricePerDay  *float64 `json:"price_per_day,omitempty"`
	PricePerWeek *float64 `json:"price_per_week,omitempty"`
	Location     string   `json:"location,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	AutoApprove  bool     `json:"auto_approve"`
}

type UpdateEquipmentRequest struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Category     string   `json:"category,omitempty"`
	PricePerHour *float64 `json:"price_per_hour,omitempty"`
	PricePerDay  *float64 `json:"price_per_day,omitempty"`
	PricePerWeek *float64 `json:"price_per_week,omitempty"`
	Location     string   `json:"location,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Available    *bool    `json:"available,omitempty"`
	AutoApprove  *bool    `json:"auto_approve,omitempty"`
}

type EquipmentFilter struct {
	Category  string
	Location  string
	Available *bool
	MinPrice  *float64
	MaxPrice  *float64
	StartDate *time.Time
	EndDate   *time.Time
	OwnerID   *uuid.UUID
}

type EquipmentAvailability struct {
	Date      time.Time `json:"date"`
	Available bool      `json:"available"`
}
