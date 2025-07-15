package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmaister/gots-template/api"
	"github.com/jmaister/gots-template/db"
	"github.com/jmaister/gots-template/handlers"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	Host       string
	Port       int
	WebappFS   embed.FS
	WebappPath string
}

// Server represents the main HTTP server with its dependencies
type Server struct {
	HTTPServer     *http.Server
	Mux            *http.ServeMux
	UserRepository db.UserRepository
	WebappFS       embed.FS
	WebappPath     string
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
		WebappFS:       serverConfig.WebappFS,
		WebappPath:     serverConfig.WebappPath,
	}

	// Configure webapp serving first (will be overridden by more specific routes)
	err := server.configureWebappRoutes()
	if err != nil {
		return nil, fmt.Errorf("error configuring webapp routes: %w", err)
	}

	// Configure OpenAPI endpoints (more specific, will take precedence)
	err = server.configureOpenAPIRoutes()
	if err != nil {
		return nil, fmt.Errorf("error configuring OpenAPI routes: %w", err)
	}

	return server, nil
}

// NewServerFromEmbedFS creates and configures a new server instance with embedded webapp files
func NewServerFromEmbedFS(host string, port int, webappFS embed.FS, webappPath string) (*Server, error) {
	// Initialize the database connection
	db.Init()

	serverConfig := &ServerConfig{
		Host:       host,
		Port:       port,
		WebappFS:   webappFS,
		WebappPath: webappPath,
	}

	return NewServer(serverConfig)
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

	// Register the API routes as defined in OpenAPI spec, "/api/" means to handle all routes with that prefix
	s.Mux.Handle("/api/", openApiHandler)

	return nil
}

// configureWebappRoutes sets up the webapp SPA serving
func (s *Server) configureWebappRoutes() error {
	log.Printf("Configuring webapp routes...")

	// Try to serve webapp from embedded FS
	webappFS, err := fs.Sub(s.WebappFS, s.WebappPath)
	if err != nil {
		log.Printf("Warning: Could not access embedded webapp files: %v", err)
		log.Println("Trying to serve webapp from filesystem...")

		// Fallback to serving from filesystem during development
		s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Skip API routes - they should be handled by the API handler
			if r.URL.Path == "/api" || (len(r.URL.Path) > 4 && r.URL.Path[:5] == "/api/") {
				// Don't handle API routes here, let them fall through to API handler
				http.NotFound(w, r)
				return
			}

			s.serveSPAFromFileSystem(w, r, s.WebappPath)
		})

		log.Println("Serving webapp SPA from filesystem")
	} else {
		// Serve the webapp files from embedded FS
		s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Skip API routes - they should be handled by the API handler
			if r.URL.Path == "/api" || (len(r.URL.Path) > 4 && r.URL.Path[:5] == "/api/") {
				// Don't handle API routes here, let them fall through to API handler
				http.NotFound(w, r)
				return
			}

			s.serveSPAFromEmbeddedFS(w, r, webappFS)
		})

		log.Println("Serving webapp SPA from embedded files")
	}

	return nil
}

// serveSPAFromFileSystem serves files from filesystem for SPA
// If the file doesn't exist, it serves index.html for client-side routing
func (s *Server) serveSPAFromFileSystem(w http.ResponseWriter, r *http.Request, rootDir string) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Try to serve the requested file
	filePath := rootDir + path
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, serve index.html for SPA routing
		http.ServeFile(w, r, rootDir+"/index.html")
		return
	}

	// File exists, serve it
	http.ServeFile(w, r, filePath)
}

// serveSPAFromEmbeddedFS serves files from embedded FS for SPA
// If the file doesn't exist, it serves index.html for client-side routing
func (s *Server) serveSPAFromEmbeddedFS(w http.ResponseWriter, r *http.Request, webappFS fs.FS) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Remove leading slash for embedded FS
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Try to open the requested file
	file, err := webappFS.Open(path)
	if err != nil {
		// File doesn't exist, serve index.html for SPA routing
		indexFile, indexErr := webappFS.Open("index.html")
		if indexErr != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		// Set appropriate content type for HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Copy index.html content to response
		indexContent, copyErr := fs.ReadFile(webappFS, "index.html")
		if copyErr != nil {
			http.Error(w, "Failed to read index.html", http.StatusInternalServerError)
			return
		}

		w.Write(indexContent)
		return
	}
	defer file.Close()

	// File exists, serve it using http.FileServer
	fileServer := http.FileServer(http.FS(webappFS))
	fileServer.ServeHTTP(w, r)
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
