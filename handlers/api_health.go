package handlers

import (
	"context"
	"time"

	"github.com/jmaister/gots-template/api"
)

// HealthCheck implements the HealthCheck operation for the api.StrictServerInterface.
func (s *StrictApiServer) HealthCheck(ctx context.Context, request api.HealthCheckRequestObject) (api.HealthCheckResponseObject, error) {
	// Calculate actual uptime since server start
	uptime := time.Since(s.StartTime)
	uptimeStr := uptime.String()

	// You can customize the version or make it dynamic
	version := "1.0.0"

	response := api.HealthCheck200JSONResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   &version,
		Uptime:    &uptimeStr,
	}
	return response, nil
}
