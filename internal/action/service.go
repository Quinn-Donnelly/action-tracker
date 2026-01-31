package action

import (
	"action-tracker/internal/shared/crypto"
	"context"
	"fmt"
	"time"
)

// Service provides business logic for actions
type Service struct {
	repo Repository
}

// NewService creates a new action service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateAction creates a new action with business validation
func (s *Service) CreateAction(ctx context.Context, req *CreateRequest) (*Action, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Set defaults
	quantity := 1
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	// Create action
	id, err := crypto.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	action := &Action{
		ID:          *id,
		Name:        *req.Name,
		Description: req.Description,
		Category:    req.Category,
		Quantity:    quantity,
		Unit:        req.Unit,
		Timestamp:   timestamp,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.Create(ctx, action)
}

// GetAction retrieves an action by ID
func (s *Service) GetAction(ctx context.Context, id crypto.UUID) (*Action, error) {
	return s.repo.GetByID(ctx, id)
}

// ListActions retrieves actions based on filters
func (s *Service) ListActions(ctx context.Context, filters *ListFilters) ([]*Action, error) {
	// Apply business logic filters and defaults
	limit := 100
	if filters.Limit != nil {
		limit = *filters.Limit
		if limit <= 0 {
			return nil, &ValidationError{Field: "limit", Message: "limit must be greater than 0"}
		}
		if limit > 1000 {
			return nil, &ValidationError{Field: "limit", Message: "limit cannot exceed 1000"}
		}
	}

	offset := 0
	if filters.Offset != nil {
		offset = *filters.Offset
		if offset < 0 {
			return nil, &ValidationError{Field: "offset", Message: "offset cannot be negative"}
		}
	}

	// Copy filters with applied business rules
	businessFilters := &ListFilters{
		Name:      filters.Name,
		Category:  filters.Category,
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
		Limit:     &limit,
		Offset:    &offset,
	}

	return s.repo.List(ctx, businessFilters)
}

// UpdateAction updates an existing action
func (s *Service) UpdateAction(ctx context.Context, id crypto.UUID, req *CreateRequest) (*Action, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if action exists first
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	updated := *existing // Copy existing
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Description != nil {
		updated.Description = req.Description
	}
	if req.Category != nil {
		updated.Category = req.Category
	}
	if req.Quantity != nil {
		updated.Quantity = *req.Quantity
	}
	if req.Unit != nil {
		updated.Unit = req.Unit
	}
	if req.Timestamp != nil {
		updated.Timestamp = *req.Timestamp
	}
	updated.UpdatedAt = time.Now()

	return s.repo.Update(ctx, id, &updated)
}

// DeleteAction removes an action
func (s *Service) DeleteAction(ctx context.Context, id crypto.UUID) error {
	return s.repo.Delete(ctx, id)
}
