package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"action-tracker/internal/shared/crypto"
	"action-tracker/internal/shared/response"
)

// Handler handles HTTP requests for actions
type Handler struct {
	service *Service
}

// NewHandler creates a new action handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateAction handles POST /api/actions
func (h *Handler) CreateAction(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, map[string]string{
			"body": "Invalid JSON format",
		})
		return
	}

	action, err := h.service.CreateAction(r.Context(), &req)
	if err != nil {
		if validationErr, ok := err.(*ValidationError); ok {
			response.ValidationError(w, map[string]string{
				validationErr.Field: validationErr.Message,
			})
			return
		}
		if IsNotFoundError(err) {
			response.NotFound(w, "Action not found")
			return
		}
		response.InternalServerError(w, "Failed to create action")
		return
	}

	response.JSON(w, http.StatusCreated, action)
}

// GetActions handles GET /api/actions
func (h *Handler) GetActions(w http.ResponseWriter, r *http.Request) {
	filters, err := parseGetActionsQuery(r.URL.Query())
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	actions, err := h.service.ListActions(r.Context(), filters)
	if err != nil {
		if validationErr, ok := err.(*ValidationError); ok {
			response.ValidationError(w, map[string]string{
				validationErr.Field: validationErr.Message,
			})
			return
		}
		response.InternalServerError(w, "Failed to list actions")
		return
	}

	response.JSON(w, http.StatusOK, actions)
}

// GetAction handles GET /api/actions/{id}
func (h *Handler) GetAction(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := crypto.ParseUUID(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid action ID format")
		return
	}

	action, err := h.service.GetAction(r.Context(), *id)
	if err != nil {
		if IsNotFoundError(err) {
			response.NotFound(w, "Action not found")
			return
		}
		response.InternalServerError(w, "Failed to get action")
		return
	}

	response.JSON(w, http.StatusOK, action)
}

// parseGetActionsQuery parses query parameters for list actions
func parseGetActionsQuery(query map[string][]string) (*ListFilters, error) {
	filters := &ListFilters{}

	// Parse name filter
	if nameValues := query["name"]; len(nameValues) > 0 && nameValues[0] != "" {
		name := nameValues[0]
		filters.Name = &name
	}

	// Parse category filter
	if categoryValues := query["category"]; len(categoryValues) > 0 && categoryValues[0] != "" {
		category := categoryValues[0]
		filters.Category = &category
	}

	// Parse limit
	if limitValues := query["limit"]; len(limitValues) > 0 {
		limit, err := strconv.Atoi(limitValues[0])
		if err != nil || limit <= 0 {
			return nil, fmt.Errorf("invalid limit parameter")
		}
		filters.Limit = &limit
	}

	// Parse offset
	if offsetValues := query["offset"]; len(offsetValues) > 0 {
		offset, err := strconv.Atoi(offsetValues[0])
		if err != nil || offset < 0 {
			return nil, fmt.Errorf("invalid offset parameter")
		}
		filters.Offset = &offset
	}

	// Parse date filters (simplified - assuming RFC3339 format)
	if startDateValues := query["start_date"]; len(startDateValues) > 0 {
		// Add date parsing if needed
	}

	if endDateValues := query["end_date"]; len(endDateValues) > 0 {
		// Add date parsing if needed
	}

	return filters, nil
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
