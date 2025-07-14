package handlers

import (
	"github.com/jmaister/gots-template/api"
)

// --- Strict API Server Implementation ---

// StrictApiServer provides an implementation of the api.StrictServerInterface.
type StrictApiServer struct {
	// Currently only implements health check endpoint
}

// NewStrictApiServer creates a new StrictApiServer.
func NewStrictApiServer() *StrictApiServer {
	return &StrictApiServer{}
}

// Ensure StrictApiServer implements StrictServerInterface
var _ api.StrictServerInterface = (*StrictApiServer)(nil)
