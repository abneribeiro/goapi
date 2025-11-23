package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/pagination"
)

var ErrEquipmentNotFound = errors.New("equipment not found")

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(ctx context.Context, equipment *model.Equipment) error {
	query := `
		INSERT INTO equipment (id, owner_id, name, description, category, price_per_hour, price_per_day, price_per_week, location, latitude, longitude, available, auto_approve, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	equipment.ID = uuid.New()
	equipment.CreatedAt = time.Now()
	equipment.UpdatedAt = time.Now()
	equipment.Available = true

	_, err := r.db.ExecContext(ctx, query,
		equipment.ID,
		equipment.OwnerID,
		equipment.Name,
		equipment.Description,
		equipment.Category,
		equipment.PricePerHour,
		equipment.PricePerDay,
		equipment.PricePerWeek,
		equipment.Location,
		equipment.Latitude,
		equipment.Longitude,
		equipment.Available,
		equipment.AutoApprove,
		equipment.CreatedAt,
		equipment.UpdatedAt,
	)

	return err
}

func (r *EquipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Equipment, error) {
	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.category, e.price_per_hour, e.price_per_day, e.price_per_week, e.location, e.latitude, e.longitude, e.available, e.auto_approve, e.created_at, e.updated_at,
		       u.id, u.email, u.name, u.phone, u.role, u.verified, u.created_at, u.updated_at
		FROM equipment e
		LEFT JOIN users u ON e.owner_id = u.id
		WHERE e.id = $1
	`

	equipment := &model.Equipment{Owner: &model.User{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&equipment.ID,
		&equipment.OwnerID,
		&equipment.Name,
		&equipment.Description,
		&equipment.Category,
		&equipment.PricePerHour,
		&equipment.PricePerDay,
		&equipment.PricePerWeek,
		&equipment.Location,
		&equipment.Latitude,
		&equipment.Longitude,
		&equipment.Available,
		&equipment.AutoApprove,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
		&equipment.Owner.ID,
		&equipment.Owner.Email,
		&equipment.Owner.Name,
		&equipment.Owner.Phone,
		&equipment.Owner.Role,
		&equipment.Owner.Verified,
		&equipment.Owner.CreatedAt,
		&equipment.Owner.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	photos, err := r.GetPhotos(ctx, equipment.ID)
	if err != nil {
		return nil, err
	}
	equipment.Photos = photos

	return equipment, nil
}

func (r *EquipmentRepository) List(ctx context.Context, filter *model.EquipmentFilter, pag pagination.Params) ([]*model.Equipment, int64, error) {
	baseQuery := `FROM equipment e WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter != nil {
		if filter.Category != "" {
			argCount++
			baseQuery += fmt.Sprintf(" AND e.category = $%d", argCount)
			args = append(args, filter.Category)
		}
		if filter.Location != "" {
			argCount++
			baseQuery += fmt.Sprintf(" AND e.location ILIKE $%d", argCount)
			args = append(args, "%"+filter.Location+"%")
		}
		if filter.Available != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND e.available = $%d", argCount)
			args = append(args, *filter.Available)
		}
		if filter.OwnerID != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND e.owner_id = $%d", argCount)
			args = append(args, *filter.OwnerID)
		}
		if filter.MinPrice != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND (e.price_per_day >= $%d OR e.price_per_hour >= $%d OR e.price_per_week >= $%d)", argCount, argCount, argCount)
			args = append(args, *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND (e.price_per_day <= $%d OR e.price_per_hour <= $%d OR e.price_per_week <= $%d)", argCount, argCount, argCount)
			args = append(args, *filter.MaxPrice)
		}
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT e.id, e.owner_id, e.name, e.description, e.category, e.price_per_hour, e.price_per_day, e.price_per_week, e.location, e.latitude, e.longitude, e.available, e.auto_approve, e.created_at, e.updated_at ` + baseQuery
	selectQuery += " ORDER BY e.created_at DESC"
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pag.PerPage, pag.Offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var equipment []*model.Equipment
	for rows.Next() {
		e := &model.Equipment{}
		err := rows.Scan(
			&e.ID,
			&e.OwnerID,
			&e.Name,
			&e.Description,
			&e.Category,
			&e.PricePerHour,
			&e.PricePerDay,
			&e.PricePerWeek,
			&e.Location,
			&e.Latitude,
			&e.Longitude,
			&e.Available,
			&e.AutoApprove,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		equipment = append(equipment, e)
	}

	return equipment, total, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	query := `
		UPDATE equipment
		SET name = $1, description = $2, category = $3, price_per_hour = $4, price_per_day = $5, price_per_week = $6, location = $7, latitude = $8, longitude = $9, available = $10, auto_approve = $11, updated_at = $12
		WHERE id = $13
	`

	equipment.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		equipment.Name,
		equipment.Description,
		equipment.Category,
		equipment.PricePerHour,
		equipment.PricePerDay,
		equipment.PricePerWeek,
		equipment.Location,
		equipment.Latitude,
		equipment.Longitude,
		equipment.Available,
		equipment.AutoApprove,
		equipment.UpdatedAt,
		equipment.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrEquipmentNotFound
	}

	return nil
}

func (r *EquipmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM equipment WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrEquipmentNotFound
	}

	return nil
}

func (r *EquipmentRepository) AddPhoto(ctx context.Context, photo *model.EquipmentPhoto) error {
	query := `
		INSERT INTO equipment_photos (id, equipment_id, url, is_primary, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	photo.ID = uuid.New()
	photo.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		photo.ID,
		photo.EquipmentID,
		photo.URL,
		photo.IsPrimary,
		photo.CreatedAt,
	)

	return err
}

func (r *EquipmentRepository) GetPhotos(ctx context.Context, equipmentID uuid.UUID) ([]model.EquipmentPhoto, error) {
	query := `
		SELECT id, equipment_id, url, is_primary, created_at
		FROM equipment_photos
		WHERE equipment_id = $1
		ORDER BY is_primary DESC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []model.EquipmentPhoto
	for rows.Next() {
		var photo model.EquipmentPhoto
		err := rows.Scan(
			&photo.ID,
			&photo.EquipmentID,
			&photo.URL,
			&photo.IsPrimary,
			&photo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (r *EquipmentRepository) DeletePhoto(ctx context.Context, photoID uuid.UUID) error {
	query := `DELETE FROM equipment_photos WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, photoID)
	return err
}

func (r *EquipmentRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM equipment ORDER BY category`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *EquipmentRepository) CheckAvailability(ctx context.Context, equipmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM reservations
		WHERE equipment_id = $1
		AND status IN ('pending', 'approved')
		AND (
			(start_date <= $2 AND end_date >= $2) OR
			(start_date <= $3 AND end_date >= $3) OR
			(start_date >= $2 AND end_date <= $3)
		)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, equipmentID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *EquipmentRepository) GetAvailabilityCalendar(ctx context.Context, equipmentID uuid.UUID, startDate, endDate time.Time) ([]model.EquipmentAvailability, error) {
	query := `
		SELECT start_date, end_date
		FROM reservations
		WHERE equipment_id = $1
		AND status IN ('pending', 'approved')
		AND start_date <= $3
		AND end_date >= $2
	`

	rows, err := r.db.QueryContext(ctx, query, equipmentID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservedDates := make(map[string]bool)
	for rows.Next() {
		var resStart, resEnd time.Time
		if err := rows.Scan(&resStart, &resEnd); err != nil {
			return nil, err
		}
		for d := resStart; !d.After(resEnd); d = d.AddDate(0, 0, 1) {
			reservedDates[d.Format("2006-01-02")] = true
		}
	}

	var availability []model.EquipmentAvailability
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		availability = append(availability, model.EquipmentAvailability{
			Date:      d,
			Available: !reservedDates[d.Format("2006-01-02")],
		})
	}

	return availability, nil
}

func (r *EquipmentRepository) Search(ctx context.Context, query string, pag pagination.Params) ([]*model.Equipment, int64, error) {
	searchTerms := strings.Fields(query)
	if len(searchTerms) == 0 {
		return r.List(ctx, nil, pag)
	}

	baseQuery := `FROM equipment e WHERE (`
	args := []interface{}{}
	conditions := []string{}

	for i, term := range searchTerms {
		conditions = append(conditions, fmt.Sprintf("(e.name ILIKE $%d OR e.description ILIKE $%d OR e.category ILIKE $%d)", i+1, i+1, i+1))
		args = append(args, "%"+term+"%")
	}
	baseQuery += strings.Join(conditions, " AND ") + ")"

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	argCount := len(args)
	selectQuery := `SELECT e.id, e.owner_id, e.name, e.description, e.category, e.price_per_hour, e.price_per_day, e.price_per_week, e.location, e.latitude, e.longitude, e.available, e.auto_approve, e.created_at, e.updated_at ` + baseQuery
	selectQuery += " ORDER BY e.created_at DESC"
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pag.PerPage, pag.Offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var equipment []*model.Equipment
	for rows.Next() {
		e := &model.Equipment{}
		err := rows.Scan(
			&e.ID,
			&e.OwnerID,
			&e.Name,
			&e.Description,
			&e.Category,
			&e.PricePerHour,
			&e.PricePerDay,
			&e.PricePerWeek,
			&e.Location,
			&e.Latitude,
			&e.Longitude,
			&e.Available,
			&e.AutoApprove,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		equipment = append(equipment, e)
	}

	return equipment, total, nil
}
