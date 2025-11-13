package session

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestContext(t *testing.T) {
	t.Run("SuccessfulRetrieval", func(t *testing.T) {
		// Setup: Create a request context
		req, _ := http.NewRequest("GET", "/test", nil)
		reqCtx := &RequestContext{
			RawRequest:      req,
			UserID:          "user123",
			IsAuthenticated: true,
			RequestID:       "req-123",
			RequestTime:     time.Now(),
		}

		// Store it in context
		ctx := context.WithValue(context.Background(), requestContextKey, reqCtx)

		// Execute
		result, err := GetRequestContext(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "user123", result.UserID)
		assert.Equal(t, "req-123", result.RequestID)
		assert.True(t, result.IsAuthenticated)
	})

	t.Run("MissingContext", func(t *testing.T) {
		// Setup: Empty context
		ctx := context.Background()

		// Execute
		result, err := GetRequestContext(ctx)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "request context is missing")
	})
}

func TestRequestContext_GetUserId(t *testing.T) {
	t.Run("AuthenticatedUser", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			UserID:          "user123",
			IsAuthenticated: true,
		}

		// Execute
		userID, err := reqCtx.GetUserId()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user123", userID)
	})

	t.Run("NotAuthenticated", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			UserID:          "",
			IsAuthenticated: false,
		}

		// Execute
		userID, err := reqCtx.GetUserId()

		// Assert
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "user not authenticated")
	})
}

func TestRequestContext_GetUserData(t *testing.T) {
	t.Run("LazyParsingFirstCall", func(t *testing.T) {
		// Setup: Create session JSON
		session := Session{
			UserID:   "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Token:    "token123",
		}
		sessionJSON, _ := json.Marshal(session)

		reqCtx := &RequestContext{
			UserDataRaw: string(sessionJSON),
		}

		// Execute first call
		result, err := reqCtx.GetUserData()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user123", result.UserID)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "token123", result.Token)

		// Verify session is cached
		assert.NotNil(t, reqCtx.Session)
	})

	t.Run("CachedSessionSecondCall", func(t *testing.T) {
		// Setup: Pre-parsed session
		session := Session{
			UserID:   "user123",
			Username: "testuser",
		}
		reqCtx := &RequestContext{
			Session: &session,
		}

		// Execute
		result, err := reqCtx.GetUserData()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user123", result.UserID)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("NoSessionData", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			UserDataRaw: "",
		}

		// Execute
		result, err := reqCtx.GetUserData()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session data not available")
		assert.Empty(t, result.UserID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			UserDataRaw: "invalid json {",
		}

		// Execute
		result, err := reqCtx.GetUserData()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid session data")
		assert.Empty(t, result.UserID)
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		// Setup: Create session JSON
		session := Session{
			UserID:   "user123",
			Username: "testuser",
		}
		sessionJSON, _ := json.Marshal(session)

		reqCtx := &RequestContext{
			UserDataRaw: string(sessionJSON),
		}

		// Execute: Multiple concurrent calls
		done := make(chan bool)
		for i := 0; i < 5; i++ {
			go func() {
				result, err := reqCtx.GetUserData()
				assert.NoError(t, err)
				assert.Equal(t, "user123", result.UserID)
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 5; i++ {
			<-done
		}

		// Assert: Session should be parsed only once and cached
		assert.NotNil(t, reqCtx.Session)
	})
}

func TestRequestContext_GetRawRequest(t *testing.T) {
	t.Run("SuccessfulRetrieval", func(t *testing.T) {
		// Setup
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Custom-Header", "test-value")

		reqCtx := &RequestContext{
			RawRequest: req,
		}

		// Execute
		result, err := reqCtx.GetRawRequest()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "GET", result.Method)
		assert.Equal(t, "/test", result.URL.Path)
		assert.Equal(t, "test-value", result.Header.Get("X-Custom-Header"))
	})

	t.Run("MissingRawRequest", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			RawRequest: nil,
		}

		// Execute
		result, err := reqCtx.GetRawRequest()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "raw request not available")
	})
}

func TestForwardAuthorizationCookie(t *testing.T) {
	t.Run("SuccessfulForward", func(t *testing.T) {
		// Setup: Create session with token
		session := Session{
			Token: "test-token-123",
		}
		sessionJSON, _ := json.Marshal(session)

		reqCtx := &RequestContext{
			UserDataRaw: string(sessionJSON),
		}

		ctx := context.WithValue(context.Background(), requestContextKey, reqCtx)

		// Create a request to modify
		req, _ := http.NewRequest("GET", "/api/test", nil)

		// Execute
		err := ForwardAuthorizationCookie(ctx, req)

		// Assert
		assert.NoError(t, err)

		// Verify cookie was added
		cookies := req.Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, CookieSession, cookies[0].Name)
		assert.Equal(t, "test-token-123", cookies[0].Value)
	})

	t.Run("MissingContext", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		req, _ := http.NewRequest("GET", "/api/test", nil)

		// Execute
		err := ForwardAuthorizationCookie(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request context is missing")
	})

	t.Run("MissingSessionData", func(t *testing.T) {
		// Setup
		reqCtx := &RequestContext{
			UserDataRaw: "",
		}
		ctx := context.WithValue(context.Background(), requestContextKey, reqCtx)
		req, _ := http.NewRequest("GET", "/api/test", nil)

		// Execute
		err := ForwardAuthorizationCookie(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session data not available")
	})
}

func TestWithAdminToken(t *testing.T) {
	t.Run("SuccessfulTokenSet", func(t *testing.T) {
		// Setup
		adminToken := "admin-token-123"
		ctx := context.Background()
		req, _ := http.NewRequest("GET", "/api/admin", nil)
		req.AddCookie(&http.Cookie{Name: "some_cookie", Value: "test"})

		// Execute
		editorFn := WithAdminToken(adminToken)
		err := editorFn(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "Bearer admin-token-123", req.Header.Get("Authorization"))
		assert.Empty(t, req.Header.Get("Cookie"), "Cookie header should be removed")
	})

	t.Run("EmptyToken", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		req, _ := http.NewRequest("GET", "/api/admin", nil)

		// Execute
		editorFn := WithAdminToken("")
		err := editorFn(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, req.Header.Get("Authorization"))
	})
}
