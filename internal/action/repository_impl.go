package action

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"action-tracker/internal/shared/crypto"
)

// SQLRepository implements Repository interface using SQL database
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a new SQL repository
func NewSQLRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

// Create saves a new action to the database
func (r *SQLRepository) Create(ctx context.Context, action *Action) (*Action, error) {
	query := `
		INSERT INTO actions (id, name, description, category, quantity, unit, timestamp, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, description, category, quantity, unit, timestamp, created_at, updated_at
	`

	var createdAction Action
	err := r.db.QueryRowContext(ctx, query,
		action.ID.String(),
		action.Name,
		action.Description,
		action.Category,
		action.Quantity,
		action.Unit,
		action.Timestamp,
		action.CreatedAt,
		action.UpdatedAt,
	).Scan(
		&createdAction.ID,
		&createdAction.Name,
		&createdAction.Description,
		&createdAction.Category,
		&createdAction.Quantity,
		&createdAction.Unit,
		&createdAction.Timestamp,
		&createdAction.CreatedAt,
		&createdAction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return &createdAction, nil
}

// GetByID retrieves an action by its ID
func (r *SQLRepository) GetByID(ctx context.Context, id crypto.UUID) (*Action, error) {
	query := `
		SELECT id, name, description, category, quantity, unit, timestamp, created_at, updated_at
		FROM actions
		WHERE id = $1
	`

	var action Action
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&action.ID,
		&action.Name,
		&action.Description,
		&action.Category,
		&action.Quantity,
		&action.Unit,
		&action.Timestamp,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Type: "action", ID: id.String()}
		}
		return nil, fmt.Errorf("failed to get action by ID: %w", err)
	}

	return &action, nil
}

// List retrieves actions based on the provided filters
func (r *SQLRepository) List(ctx context.Context, filters *ListFilters) ([]*Action, error) {
	query := `
		SELECT id, name, description, category, quantity, unit, timestamp, created_at, updated_at
		FROM actions
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filters.Name != nil && *filters.Name != "" {
		query += " AND name = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.Name)
		argIndex++
	}

	if filters.Category != nil && *filters.Category != "" {
		query += " AND category = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.Category)
		argIndex++
	}

	if filters.StartDate != nil {
		query += " AND timestamp >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		query += " AND timestamp <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.EndDate)
		argIndex++
	}

	query += " ORDER BY timestamp DESC"

	limit := 100
	offset := 0
	if filters.Limit != nil {
		limit = *filters.Limit
	}
	if filters.Offset != nil {
		offset = *filters.Offset
	}

	query += " LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}
	defer rows.Close()

	var actions []*Action
	for rows.Next() {
		var action Action
		err := rows.Scan(
			&action.ID,
			&action.Name,
			&action.Description,
			&action.Category,
			&action.Quantity,
			&action.Unit,
			&action.Timestamp,
			&action.CreatedAt,
			&action.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action: %w", err)
		}
		actions = append(actions, &action)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating actions: %w", err)
	}

	return actions, nil
}

// Update updates an existing action
func (r *SQLRepository) Update(ctx context.Context, id crypto.UUID, action *Action) (*Action, error) {
	query := `
		UPDATE actions
		SET name = $2, description = $3, category = $4, quantity = $5, unit = $6, timestamp = $7, updated_at = $8
		WHERE id = $1
		RETURNING id, name, description, category, quantity, unit, timestamp, created_at, updated_at
	`

	var updatedAction Action
	err := r.db.QueryRowContext(ctx, query,
		id.String(),
		action.Name,
		action.Description,
		action.Category,
		action.Quantity,
		action.Unit,
		action.Timestamp,
		time.Now(),
	).Scan(
		&updatedAction.ID,
		&updatedAction.Name,
		&updatedAction.Description,
		&updatedAction.Category,
		&updatedAction.Quantity,
		&updatedAction.Unit,
		&updatedAction.Timestamp,
		&updatedAction.CreatedAt,
		&updatedAction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Type: "action", ID: id.String()}
		}
		return nil, fmt.Errorf("failed to update action: %w", err)
	}

	return &updatedAction, nil
}

// Delete removes an action by its ID
func (r *SQLRepository) Delete(ctx context.Context, id crypto.UUID) error {
	query := `DELETE FROM actions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return &NotFoundError{Type: "action", ID: id.String()}
	}

	return nil
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Type string
	ID   string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Type, e.ID)
}
