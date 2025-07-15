package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jmaister/gots-template/server"
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

	// Create and start the server
	srv, err := server.NewServerFromEmbedFS("127.0.0.1", 8081, webappEmbedFS, "webapp/dist")
	if err != nil {
		log.Fatalf("FATAL: Failed to create server: %v", err)
	}

	err = srv.Start()
	if err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}

	log.Println("Server shut down gracefully.")
}
