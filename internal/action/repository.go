package action

import (
	"action-tracker/internal/shared/crypto"
	"context"
)

// Repository defines the interface for action data access
type Repository interface {
	Create(ctx context.Context, action *Action) (*Action, error)
	GetByID(ctx context.Context, id crypto.UUID) (*Action, error)
	List(ctx context.Context, filters *ListFilters) ([]*Action, error)
	Update(ctx context.Context, id crypto.UUID, action *Action) (*Action, error)
	Delete(ctx context.Context, id crypto.UUID) error
}
