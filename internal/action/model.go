package action

import (
	"action-tracker/internal/shared/crypto"
	"embed"
	"time"
)

//go:embed migrations/*.sql
var MigrationFS embed.FS

// GetMigrations returns the embedded migration files
func GetMigrations() embed.FS {
	return MigrationFS
}

// Action represents an action in the system
type Action struct {
	ID          crypto.UUID `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description *string     `json:"description" db:"description"`
	Category    *string     `json:"category" db:"category"`
	Quantity    int         `json:"quantity" db:"quantity"`
	Unit        *string     `json:"unit" db:"unit"`
	Timestamp   time.Time   `json:"timestamp" db:"timestamp"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateRequest represents a request to create an action
type CreateRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
	Quantity    *int       `json:"quantity"`
	Unit        *string    `json:"unit"`
	Timestamp   *time.Time `json:"timestamp"`
}

// ListFilters represents filters for listing actions
type ListFilters struct {
	Name      *string    `json:"name"`
	Category  *string    `json:"category"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Limit     *int       `json:"limit"`
	Offset    *int       `json:"offset"`
}

// Validate validates the CreateRequest
func (r *CreateRequest) Validate() error {
	if r.Name == nil || *r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}

	quantity := 1
	if r.Quantity != nil {
		quantity = *r.Quantity
		if quantity <= 0 {
			return &ValidationError{Field: "quantity", Message: "quantity must be greater than 0"}
		}
	}

	if r.Timestamp != nil && r.Timestamp.IsZero() {
		return &ValidationError{Field: "timestamp", Message: "timestamp cannot be zero"}
	}

	return nil
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
