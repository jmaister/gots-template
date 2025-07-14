package handlers

import (
	"context"
	"time"

	"github.com/jmaister/gots-template/api"
)

// HealthCheck implements the HealthCheck operation for the api.StrictServerInterface.
func (s *StrictApiServer) HealthCheck(ctx context.Context, request api.HealthCheckRequestObject) (api.HealthCheckResponseObject, error) {
	// You can customize these values or make them dynamic
	version := "1.0.0"
	uptime := "0s" // You could calculate actual uptime here

	response := api.HealthCheck200JSONResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   &version,
		Uptime:    &uptime,
	}
	return response, nil
}
