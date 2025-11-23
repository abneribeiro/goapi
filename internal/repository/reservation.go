package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
)

var ErrReservationNotFound = errors.New("reservation not found")

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	query := `
		INSERT INTO reservations (id, equipment_id, renter_id, start_date, end_date, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	reservation.ID = uuid.New()
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.EquipmentID,
		reservation.RenterID,
		reservation.StartDate,
		reservation.EndDate,
		reservation.Status,
		reservation.TotalPrice,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)

	return err
}

func (r *ReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Reservation, error) {
	query := `
		SELECT r.id, r.equipment_id, r.renter_id, r.start_date, r.end_date, r.status, r.total_price, r.cancellation_reason, r.created_at, r.updated_at,
		       e.id, e.name, e.category, e.price_per_hour, e.price_per_day, e.price_per_week, e.location, e.owner_id,
		       u.id, u.email, u.name, u.phone
		FROM reservations r
		LEFT JOIN equipment e ON r.equipment_id = e.id
		LEFT JOIN users u ON r.renter_id = u.id
		WHERE r.id = $1
	`

	reservation := &model.Reservation{
		Equipment: &model.Equipment{},
		Renter:    &model.User{},
	}
	var cancellationReason sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.EquipmentID,
		&reservation.RenterID,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.Status,
		&reservation.TotalPrice,
		&cancellationReason,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Equipment.ID,
		&reservation.Equipment.Name,
		&reservation.Equipment.Category,
		&reservation.Equipment.PricePerHour,
		&reservation.Equipment.PricePerDay,
		&reservation.Equipment.PricePerWeek,
		&reservation.Equipment.Location,
		&reservation.Equipment.OwnerID,
		&reservation.Renter.ID,
		&reservation.Renter.Email,
		&reservation.Renter.Name,
		&reservation.Renter.Phone,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if cancellationReason.Valid {
		reservation.CancellationReason = cancellationReason.String
	}

	return reservation, nil
}

func (r *ReservationRepository) List(ctx context.Context, filter *model.ReservationFilter, pag pagination.Params) ([]*model.Reservation, int64, error) {
	baseQuery := `FROM reservations r
		LEFT JOIN equipment e ON r.equipment_id = e.id
		WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter != nil {
		if filter.RenterID != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND r.renter_id = $%d", argCount)
			args = append(args, *filter.RenterID)
		}
		if filter.OwnerID != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND e.owner_id = $%d", argCount)
			args = append(args, *filter.OwnerID)
		}
		if filter.EquipmentID != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND r.equipment_id = $%d", argCount)
			args = append(args, *filter.EquipmentID)
		}
		if filter.Status != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND r.status = $%d", argCount)
			args = append(args, *filter.Status)
		}
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT r.id, r.equipment_id, r.renter_id, r.start_date, r.end_date, r.status, r.total_price, r.cancellation_reason, r.created_at, r.updated_at,
		e.id, e.name, e.category, e.location ` + baseQuery
	selectQuery += " ORDER BY r.created_at DESC"
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pag.PerPage, pag.Offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reservations []*model.Reservation
	for rows.Next() {
		res := &model.Reservation{Equipment: &model.Equipment{}}
		var cancellationReason sql.NullString

		err := rows.Scan(
			&res.ID,
			&res.EquipmentID,
			&res.RenterID,
			&res.StartDate,
			&res.EndDate,
			&res.Status,
			&res.TotalPrice,
			&cancellationReason,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Equipment.ID,
			&res.Equipment.Name,
			&res.Equipment.Category,
			&res.Equipment.Location,
		)
		if err != nil {
			return nil, 0, err
		}

		if cancellationReason.Valid {
			res.CancellationReason = cancellationReason.String
		}

		reservations = append(reservations, res)
	}

	return reservations, total, nil
}

func (r *ReservationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.ReservationStatus, reason string) error {
	query := `
		UPDATE reservations
		SET status = $1, cancellation_reason = $2, updated_at = $3
		WHERE id = $4
	`

	var reasonPtr *string
	if reason != "" {
		reasonPtr = &reason
	}

	result, err := r.db.ExecContext(ctx, query, status, reasonPtr, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrReservationNotFound
	}

	return nil
}

func (r *ReservationRepository) GetEquipmentOwnerID(ctx context.Context, reservationID uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT e.owner_id
		FROM reservations r
		JOIN equipment e ON r.equipment_id = e.id
		WHERE r.id = $1
	`

	var ownerID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, reservationID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrReservationNotFound
		}
		return uuid.Nil, err
	}

	return ownerID, nil
}
