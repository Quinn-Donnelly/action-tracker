package server

import (
	"action-tracker/internal/action"
	"action-tracker/internal/shared/config"
	"database/sql"
	"log"
	"net/http"
)

// Use middleware from helpers
type middlewareFunc func(http.Handler) http.Handler

// Server represents the HTTP server
type Server struct {
	addr    string
	handler http.Handler
	db      *sql.DB
}

// NewServer creates a new HTTP server with all routes and middleware
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	// Create action feature components
	actionRepo := action.NewSQLRepository(db)
	actionService := action.NewService(actionRepo)
	actionHandler := action.NewHandler(actionService)

	// Set up routes
	mux := http.NewServeMux()

	// Action routes
	mux.HandleFunc("POST /api/actions", actionHandler.CreateAction)
	mux.HandleFunc("GET /api/actions", actionHandler.GetActions)
	mux.HandleFunc("GET /api/actions/{id}", actionHandler.GetAction)

	// Apply middleware chain
	handler := Chain(mux,
		recoveryMiddleware,
		loggingMiddleware,
		corsMiddleware,
	)

	return &Server{
		addr:    ":" + cfg.Port,
		handler: handler,
		db:      db,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Server starting on port %s", s.addr)
	return http.ListenAndServe(s.addr, s.handler)
}

// Close shuts down the server and database connection
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
