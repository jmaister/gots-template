package session

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrictInjectHTTPRequestMiddleware(t *testing.T) {
	t.Run("InjectRequestWithUserID", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set(HeaderUserId, "user123")

		var capturedCtx context.Context
		nextCalled := false

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			nextCalled = true
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)
		assert.True(t, nextCalled)

		// Verify context was injected
		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.NotNil(t, reqCtx)
		assert.Equal(t, "user123", reqCtx.UserID)
		assert.True(t, reqCtx.IsAuthenticated)
		assert.NotNil(t, reqCtx.RawRequest)
		assert.NotEmpty(t, reqCtx.RequestID)
		assert.False(t, reqCtx.RequestTime.IsZero())
	})

	t.Run("InjectRequestWithSessionData", func(t *testing.T) {
		// Setup
		session := Session{
			UserID:   "user123",
			Username: "testuser",
			Email:    "test@example.com",
		}
		sessionJSON, _ := json.Marshal(session)

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set(HeaderUserId, "user123")
		req.Header.Set(HeaderUserData, string(sessionJSON))

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.Equal(t, "user123", reqCtx.UserID)
		assert.NotEmpty(t, reqCtx.UserDataRaw)

		// Lazy parsing: Session should be nil until GetUserData is called
		assert.Nil(t, reqCtx.Session)

		// Now parse it
		userData, err := reqCtx.GetUserData()
		assert.NoError(t, err)
		assert.Equal(t, "testuser", userData.Username)
		assert.Equal(t, "test@example.com", userData.Email)

		// Session should now be cached
		assert.NotNil(t, reqCtx.Session)
	})

	t.Run("InjectRequestWithoutAuth", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/public", nil)

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.Empty(t, reqCtx.UserID)
		assert.False(t, reqCtx.IsAuthenticated)
		assert.Empty(t, reqCtx.UserDataRaw)
	})

	t.Run("InjectRequestWithExistingRequestID", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-Request-ID", "existing-req-id-123")

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.Equal(t, "existing-req-id-123", reqCtx.RequestID)
	})

	t.Run("InjectRequestGeneratesRequestID", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/test", nil)
		// No X-Request-ID header

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.NotEmpty(t, reqCtx.RequestID)
		// Should be a valid UUID format
		assert.Len(t, reqCtx.RequestID, 36) // UUID length with dashes
	})

	t.Run("InjectRequestWithInvalidJSON", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set(HeaderUserId, "user123")
		req.Header.Set(HeaderUserData, "invalid json {{{")

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)

		// Assert: Middleware should not fail, just store raw data
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.Equal(t, "invalid json {{{", reqCtx.UserDataRaw)

		// Parsing should fail when GetUserData is called
		_, err = reqCtx.GetUserData()
		assert.Error(t, err)
	})

	t.Run("RequestTimeIsRecent", func(t *testing.T) {
		// Setup
		beforeTest := time.Now()
		req := httptest.NewRequest("GET", "/api/test", nil)

		var capturedCtx context.Context

		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		middleware := StrictInjectHTTPRequestMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		_, err := middleware(context.Background(), w, req, nil)
		afterTest := time.Now()

		// Assert
		assert.NoError(t, err)

		reqCtx, err := GetRequestContext(capturedCtx)
		assert.NoError(t, err)
		assert.True(t, reqCtx.RequestTime.After(beforeTest) || reqCtx.RequestTime.Equal(beforeTest))
		assert.True(t, reqCtx.RequestTime.Before(afterTest) || reqCtx.RequestTime.Equal(afterTest))
	})
}

func TestStrictCORSMiddleware(t *testing.T) {
	t.Run("SetsCORSHeaders", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("GET", "/api/test", nil)

		nextCalled := false
		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			nextCalled = true
			return "response", nil
		}

		middleware := StrictCORSMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		result, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "response", result)
		assert.True(t, nextCalled)

		// Check CORS headers
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("HandlesOPTIONSRequest", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest("OPTIONS", "/api/test", nil)

		nextCalled := false
		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			nextCalled = true
			return "response", nil
		}

		middleware := StrictCORSMiddleware(next, "testOperation")

		// Execute
		w := httptest.NewRecorder()
		result, err := middleware(context.Background(), w, req, nil)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.False(t, nextCalled, "Next handler should not be called for OPTIONS")
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Check CORS headers are still set
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
