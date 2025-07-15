package handlers

import (
	"time"

	"github.com/jmaister/gots-template/api"
)

// --- Strict API Server Implementation ---

// StrictApiServer provides an implementation of the api.StrictServerInterface.
type StrictApiServer struct {
	// StartTime records when the server was started for uptime calculation
	StartTime time.Time
}

// NewStrictApiServer creates a new StrictApiServer.
func NewStrictApiServer() *StrictApiServer {
	return &StrictApiServer{
		StartTime: time.Now(),
	}
}

// Ensure StrictApiServer implements StrictServerInterface
var _ api.StrictServerInterface = (*StrictApiServer)(nil)
