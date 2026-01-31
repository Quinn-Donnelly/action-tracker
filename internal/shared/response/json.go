package response

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// JSON writes a JSON response with the given status code and data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Error writes an error response with detailed information
func Error(w http.ResponseWriter, statusCode int, message string, details map[string]interface{}) {
	errorResp := map[string]interface{}{
		"error": map[string]interface{}{
			"code":      statusCode,
			"message":   message,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Include details if provided (only in development or explicitly allowed)
	if details != nil {
		errorResp["error"].(map[string]interface{})["details"] = details
	}

	JSON(w, statusCode, errorResp)
}

// ValidationError writes a validation error response with field-level details
func ValidationError(w http.ResponseWriter, fieldErrors map[string]string) {
	details := make(map[string]interface{})
	for field, error := range fieldErrors {
		details[field] = error
	}

	Error(w, http.StatusBadRequest, "Validation failed", details)
}

// NotFound writes a standard 404 error response
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, nil)
}

// BadRequest writes a standard 400 error response
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message, nil)
}

// InternalServerError writes a standard 500 error response
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message, nil)
}
