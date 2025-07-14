package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmaister/gots-template/api"
	"github.com/jmaister/gots-template/db"
	"github.com/jmaister/gots-template/handlers"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	Host string
	Port int
}

// Server represents the main HTTP server with its dependencies
type Server struct {
	HTTPServer     *http.Server
	Mux            *http.ServeMux
	UserRepository db.UserRepository
}

// NewServer creates and configures a new server instance
func NewServer(serverConfig *ServerConfig) (*Server, error) {
	// Initialize the database connection
	db.Init()

	mux := http.NewServeMux()

	// Create server handler
	var handler http.Handler = mux

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}

	userRepository := db.NewDBUserRepository(db.GetConnection())

	server := &Server{
		HTTPServer:     httpServer,
		Mux:            mux,
		UserRepository: userRepository,
	}

	// Configure OpenAPI endpoints
	err := server.configureOpenAPIRoutes()
	if err != nil {
		return nil, fmt.Errorf("error configuring OpenAPI routes: %w", err)
	}

	return server, nil
}

// configureOpenAPIRoutes sets up the OpenAPI endpoints
func (s *Server) configureOpenAPIRoutes() error {
	log.Printf("Registering OpenAPI routes...")

	// Create the strict API server
	strictApiServer := handlers.NewStrictApiServer()

	// Configure strict handler options with simple error handling
	strictHandlerOptions := api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Request error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Response error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		},
	}

	// Create the standard API server without middleware
	standardApiServer := api.NewStrictHandlerWithOptions(
		strictApiServer,
		[]api.StrictMiddlewareFunc{}, // No middleware for now
		strictHandlerOptions,
	)

	// Create the OpenAPI handler
	openApiHandler := api.HandlerWithOptions(standardApiServer, api.StdHTTPServerOptions{
		BaseURL:     "",
		Middlewares: []api.MiddlewareFunc{}, // No additional middlewares
	})

	// Register the API routes under /api prefix
	apiPrefix := "/api"
	if !strings.HasSuffix(apiPrefix, "/") {
		apiPrefix += "/"
	}

	s.Mux.Handle(apiPrefix, http.StripPrefix(strings.TrimSuffix("/api", "/"), openApiHandler))
	log.Printf("Registered OpenAPI Routes under prefix: %s", apiPrefix)

	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.HTTPServer.Addr)
	return s.HTTPServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	log.Printf("Stopping server...")
	return s.HTTPServer.Close()
}
