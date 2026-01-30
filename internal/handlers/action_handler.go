package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"action-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActionHandler struct {
	db *sql.DB
}

func NewActionHandler(db *sql.DB) *ActionHandler {
	return &ActionHandler{db: db}
}

func (h *ActionHandler) CreateAction(c *gin.Context) {
	var req models.CreateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quantity := 1
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	query := `
		INSERT INTO actions (name, description, category, quantity, unit, timestamp, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, name, description, category, quantity, unit, timestamp, created_at, updated_at
	`

	var action models.Action
	err := h.db.QueryRow(
		query,
		req.Name,
		req.Description,
		req.Category,
		quantity,
		req.Unit,
		timestamp,
	).Scan(
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create action"})
		return
	}

	c.JSON(http.StatusCreated, action)
}

func (h *ActionHandler) GetActions(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	limit := 100
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	query := `
		SELECT id, name, description, category, quantity, unit, timestamp, created_at, updated_at
		FROM actions
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if name != "" {
		query += " AND name = $" + strconv.Itoa(argIndex)
		args = append(args, name)
		argIndex++
	}

	if category != "" {
		query += " AND category = $" + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if startDate != "" {
		if parsedStartDate, err := time.Parse(time.RFC3339, startDate); err == nil {
			query += " AND timestamp >= $" + strconv.Itoa(argIndex)
			args = append(args, parsedStartDate)
			argIndex++
		}
	}

	if endDate != "" {
		if parsedEndDate, err := time.Parse(time.RFC3339, endDate); err == nil {
			query += " AND timestamp <= $" + strconv.Itoa(argIndex)
			args = append(args, parsedEndDate)
			argIndex++
		}
	}

	query += " ORDER BY timestamp DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch actions"})
		return
	}
	defer rows.Close()

	var actions []models.Action
	for rows.Next() {
		var action models.Action
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan action"})
			return
		}
		actions = append(actions, action)
	}

	c.JSON(http.StatusOK, actions)
}

func (h *ActionHandler) GetAction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}

	query := `
		SELECT id, name, description, category, quantity, unit, timestamp, created_at, updated_at
		FROM actions
		WHERE id = $1
	`

	var action models.Action
	err = h.db.QueryRow(query, id).Scan(
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

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch action"})
		return
	}

	c.JSON(http.StatusOK, action)
}
