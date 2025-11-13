package session

import (
	"context"
	"net/http"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// InjectHTTPRequestMiddleware injects the *http.Request into the context for endpoint handlers.
func InjectHTTPRequestMiddleware(next func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error), operationName string) func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, HTTPRequestContextKey, r)
		return next(ctx, w, r, request)
	}
}

// StrictInjectHTTPRequestMiddleware adapts InjectHTTPRequestMiddleware to StrictMiddlewareFunc type.
func StrictInjectHTTPRequestMiddleware(next strictnethttp.StrictHTTPHandlerFunc, operationName string) strictnethttp.StrictHTTPHandlerFunc {
	return InjectHTTPRequestMiddleware(next, operationName)
}

// CORSMiddleware allows any domain for CORS requests.
func CORSMiddleware(next func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error), operationName string) func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
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

// StrictCORSMiddleware adapts CORSMiddleware to StrictMiddlewareFunc type.
func StrictCORSMiddleware(next strictnethttp.StrictHTTPHandlerFunc, operationName string) strictnethttp.StrictHTTPHandlerFunc {
	return CORSMiddleware(next, operationName)
}
