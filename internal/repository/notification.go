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

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, read, reference_id, reference_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	notification.ID = uuid.New()
	notification.Read = false
	notification.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Read,
		notification.ReferenceID,
		notification.ReferenceType,
		notification.CreatedAt,
	)

	return err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, read, reference_id, reference_type, created_at
		FROM notifications
		WHERE id = $1
	`

	notification := &model.Notification{}
	var refID sql.NullString
	var refType sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.Read,
		&refID,
		&refType,
		&notification.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}

	if refID.Valid {
		id, _ := uuid.Parse(refID.String)
		notification.ReferenceID = &id
	}
	if refType.Valid {
		notification.ReferenceType = refType.String
	}

	return notification, nil
}

func (r *NotificationRepository) List(ctx context.Context, filter *model.NotificationFilter, pag pagination.Params) ([]*model.Notification, int64, error) {
	baseQuery := `FROM notifications WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter != nil {
		if filter.UserID != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, *filter.UserID)
		}
		if filter.Read != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND read = $%d", argCount)
			args = append(args, *filter.Read)
		}
		if filter.Type != nil {
			argCount++
			baseQuery += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, *filter.Type)
		}
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT id, user_id, type, title, message, read, reference_id, reference_type, created_at ` + baseQuery
	selectQuery += " ORDER BY created_at DESC"
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pag.PerPage, pag.Offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		var refID sql.NullString
		var refType sql.NullString

		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Read,
			&refID,
			&refType,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if refID.Valid {
			id, _ := uuid.Parse(refID.String)
			n.ReferenceID = &id
		}
		if refType.Valid {
			n.ReferenceType = refType.String
		}

		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET read = true WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotificationNotFound
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotificationNotFound
	}

	return nil
}
