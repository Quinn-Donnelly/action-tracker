package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Action struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Category    *string   `json:"category" db:"category"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Unit        *string   `json:"unit" db:"unit"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Habit struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	Category       *string   `json:"category" db:"category"`
	TargetQuantity int       `json:"target_quantity" db:"target_quantity"`
	Unit           *string   `json:"unit" db:"unit"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateActionRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
	Quantity    *int       `json:"quantity"`
	Unit        *string    `json:"unit"`
	Timestamp   *time.Time `json:"timestamp"`
}

type GetActionsRequest struct {
	Name      *string    `json:"name"`
	Category  *string    `json:"category"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Limit     *int       `json:"limit"`
	Offset    *int       `json:"offset"`
}

// UUIDValue implements the driver.Valuer interface for UUID
func UUIDValue(u uuid.UUID) (driver.Value, error) {
	return u.String(), nil
}

// UUIDScan implements the sql.Scanner interface for UUID
func UUIDScan(u *uuid.UUID, src interface{}) error {
	if src == nil {
		*u = uuid.Nil
		return nil
	}

	switch src := src.(type) {
	case string:
		parsed, err := uuid.Parse(src)
		if err != nil {
			return err
		}
		*u = parsed
	case []byte:
		parsed, err := uuid.Parse(string(src))
		if err != nil {
			return err
		}
		*u = parsed
	default:
		return fmt.Errorf("cannot scan %T into UUID", src)
	}

	return nil
}
