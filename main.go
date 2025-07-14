package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// webapp/dist will be embedded here after building the frontend
//
//go:embed webapp/dist
var webappEmbedFS embed.FS

// Version information - will be injected during build by GoReleaser
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	buildOS   = "unknown"
	buildArch = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "tg",
	Short: "GOTS Template CLI",
	Long:  `A CLI for managing and running the GOTS Template application.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version, commit, build date, and build platform information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GOTS Template\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		fmt.Printf("Build OS: %s\n", buildOS)
		fmt.Printf("Build Arch: %s\n", buildArch)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the GOTS Template application",
	Long:  `Starts the GOTS Template application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

// --- Main Function ---

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServer() {
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Include file/line number
	log.Println("Starting server...")

	// Create a simple HTTP server
	mux := http.NewServeMux()

	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve webapp files as SPA
	webappFS, err := fs.Sub(webappEmbedFS, "webapp/dist")
	if err != nil {
		log.Printf("Warning: Could not access embedded webapp files: %v", err)
		log.Println("Trying to serve webapp from filesystem...")

		// Fallback to serving from filesystem during development
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// For API routes, return 404
			if r.URL.Path == "/api" || len(r.URL.Path) > 4 && r.URL.Path[:5] == "/api/" {
				http.NotFound(w, r)
				return
			}

			serveSPAFromFileSystem(w, r, "webapp/dist")
		})

		log.Println("Serving webapp SPA from filesystem")
	} else {
		// Serve the webapp files from embedded FS
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// For API routes, return 404
			if r.URL.Path == "/api" || len(r.URL.Path) > 4 && r.URL.Path[:5] == "/api/" {
				http.NotFound(w, r)
				return
			}

			serveSPAFromEmbeddedFS(w, r, webappFS)
		})

		log.Println("Serving webapp SPA from embedded files")
	}

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Printf("Server listening on http://localhost%s", server.Addr)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}

	log.Println("Server shut down gracefully.")
}

// serveSPAFromFileSystem serves files from filesystem for SPA
// If the file doesn't exist, it serves index.html for client-side routing
func serveSPAFromFileSystem(w http.ResponseWriter, r *http.Request, rootDir string) {
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
func serveSPAFromEmbeddedFS(w http.ResponseWriter, r *http.Request, webappFS fs.FS) {
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
		if _, copyErr := fs.ReadFile(webappFS, "index.html"); copyErr != nil {
			http.Error(w, "Failed to read index.html", http.StatusInternalServerError)
			return
		}

		indexContent, _ := fs.ReadFile(webappFS, "index.html")
		w.Write(indexContent)
		return
	}
	defer file.Close()

	// File exists, serve it using http.FileServer
	fileServer := http.FileServer(http.FS(webappFS))
	fileServer.ServeHTTP(w, r)
}

// --- Server Functions ---
