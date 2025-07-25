package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/jmaister/gots-template/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Create a new StrictApiServer instance
	s := NewStrictApiServer()

	t.Run("SuccessfulHealthCheck", func(t *testing.T) {
		// Setup: Create a health check request
		ctx := context.Background()
		req := api.HealthCheckRequestObject{}

		// Record the time before the call
		beforeCall := time.Now()

		// Execute the health check
		resp, err := s.HealthCheck(ctx, req)

		// Record the time after the call
		afterCall := time.Now()

		// Assert: No error should occur
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		// Assert: Response should be of the correct type
		healthResp, ok := resp.(api.HealthCheck200JSONResponse)
		assert.True(t, ok, "Response should be HealthCheck200JSONResponse")

		// Assert: Status should be "ok"
		assert.Equal(t, "ok", healthResp.Status)

		// Assert: Version should be set
		assert.NotNil(t, healthResp.Version)
		assert.Equal(t, "1.0.0", *healthResp.Version)

		// Assert: Uptime should be set
		assert.NotNil(t, healthResp.Uptime)
		// Parse the uptime to verify it's a valid duration
		uptimeDuration, err := time.ParseDuration(*healthResp.Uptime)
		assert.NoError(t, err, "Uptime should be a valid duration string")
		// Uptime should be positive and reasonable (less than a few seconds for a test)
		assert.True(t, uptimeDuration >= 0, "Uptime should be positive")
		assert.True(t, uptimeDuration < 10*time.Second, "Uptime should be reasonable for a test")

		// Assert: Timestamp should be within a reasonable range (between before and after the call)
		assert.True(t, healthResp.Timestamp.After(beforeCall) || healthResp.Timestamp.Equal(beforeCall),
			"Timestamp should be after or equal to the time before the call")
		assert.True(t, healthResp.Timestamp.Before(afterCall) || healthResp.Timestamp.Equal(afterCall),
			"Timestamp should be before or equal to the time after the call")

		// Assert: Timestamp should not be zero
		assert.False(t, healthResp.Timestamp.IsZero(), "Timestamp should not be zero")
	})

}
