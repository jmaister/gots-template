package session

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// StrictInjectHTTPRequestMiddleware injects the *http.Request and request metadata into the context.
// Session data is stored as raw JSON and parsed lazily only when GetUserData() is called.
func StrictInjectHTTPRequestMiddleware(next strictnethttp.StrictHTTPHandlerFunc, operationName string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		// Generate or extract request ID for tracing
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create request context with raw data (lazy parsing)
		reqCtx := &RequestContext{
			RawRequest:  r,
			RequestID:   requestID,
			RequestTime: time.Now(),
		}

		// Parse user ID from header (cheap operation)
		if userID := r.Header.Get(HeaderUserId); userID != "" {
			reqCtx.UserID = userID
			reqCtx.IsAuthenticated = true
		}

		// Store raw session data for lazy parsing
		reqCtx.UserDataRaw = r.Header.Get(HeaderUserData)

		// Store the complete context
		ctx = context.WithValue(ctx, requestContextKey, reqCtx)
		return next(ctx, w, r, request)
	}
}

// StrictCORSMiddleware allows any domain for CORS requests.
func StrictCORSMiddleware(next strictnethttp.StrictHTTPHandlerFunc, operationName string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return nil, nil
		}
		return next(ctx, w, r, request)
	}
}
