package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	client "github.com/jmaister/taronja-gateway-clients/go"
)

const (
	HeaderUserId   = "X-User-Id"
	HeaderUserData = "X-User-Data"
	CookieSession  = "tg_session_token"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const requestContextKey contextKey = "requestContext"

// RequestContext holds all request-related data including the raw request and parsed session.
// Session data is parsed lazily on first access to avoid unnecessary overhead.
type RequestContext struct {
	RawRequest      *http.Request
	UserID          string
	UserDataRaw     string   // Raw JSON string for lazy parsing
	Session         *Session // Parsed lazily on first GetUserData() call
	IsAuthenticated bool
	RequestID       string     // Request ID for tracing
	RequestTime     time.Time  // When the request started
	mu              sync.Mutex // Protects Session field during lazy parsing
}

// Session is a copy of db.Session from Taronja Gateway
type Session struct {
	ID              int64      `json:"ID"`
	CreatedAt       time.Time  `json:"CreatedAt"`
	UpdatedAt       time.Time  `json:"UpdatedAt"`
	DeletedAt       *time.Time `json:"DeletedAt"`
	Token           string     `json:"Token"`
	UserID          string     `json:"UserID"`
	Username        string     `json:"Username"`
	Email           string     `json:"Email"`
	IsAuthenticated bool       `json:"IsAuthenticated"`
	IsAdmin         bool       `json:"IsAdmin"`
	ValidUntil      time.Time  `json:"ValidUntil"`
	Provider        string     `json:"Provider"`
	ClosedOn        *time.Time `json:"ClosedOn"`
	LastActivity    time.Time  `json:"LastActivity"`
	SessionName     string     `json:"SessionName"`
	CreatedFrom     string     `json:"CreatedFrom"`
	IPAddress       string     `json:"IPAddress"`
	UserAgent       string     `json:"UserAgent"`
	Referrer        string     `json:"Referrer"`
	BrowserFamily   string     `json:"BrowserFamily"`
	BrowserVersion  string     `json:"BrowserVersion"`
	OSFamily        string     `json:"OSFamily"`
	OSVersion       string     `json:"OSVersion"`
	DeviceFamily    string     `json:"DeviceFamily"`
	DeviceBrand     string     `json:"DeviceBrand"`
	DeviceModel     string     `json:"DeviceModel"`
	GeoLocation     string     `json:"GeoLocation"`
	Latitude        float64    `json:"Latitude"`
	Longitude       float64    `json:"Longitude"`
	City            string     `json:"City"`
	ZipCode         string     `json:"ZipCode"`
	Country         string     `json:"Country"`
	CountryCode     string     `json:"CountryCode"`
	Region          string     `json:"Region"`
	Continent       string     `json:"Continent"`
	JA4Fingerprint  string     `json:"JA4Fingerprint"`
}

// GetRequestContext retrieves the complete request context from the context.
// This provides access to both the raw request and parsed session data.
func GetRequestContext(ctx context.Context) (*RequestContext, error) {
	reqCtx, ok := ctx.Value(requestContextKey).(*RequestContext)
	if !ok || reqCtx == nil {
		return nil, fmt.Errorf("request context is missing")
	}
	return reqCtx, nil
}

// GetUserId retrieves the user ID from the request context.
// This is a simple lookup with no header parsing.
func (rc *RequestContext) GetUserId() (string, error) {
	if rc.UserID == "" {
		return "", fmt.Errorf("user not authenticated")
	}
	return rc.UserID, nil
}

// GetUserData retrieves the session data from the request context.
// This function performs lazy parsing: the JSON is unmarshaled only on the first call,
// then cached for subsequent calls within the same request.
func (rc *RequestContext) GetUserData() (Session, error) {
	// Lazy parsing with mutex protection
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// If already parsed, return cached session
	if rc.Session != nil {
		return *rc.Session, nil
	}

	// Parse session data if not already done
	if rc.UserDataRaw == "" {
		return Session{}, fmt.Errorf("session data not available")
	}

	var session Session
	if err := json.Unmarshal([]byte(rc.UserDataRaw), &session); err != nil {
		log.Printf("Error: Failed to unmarshal user data: %v", err)
		return Session{}, fmt.Errorf("invalid session data: %w", err)
	}

	// Cache the parsed session
	rc.Session = &session
	return session, nil
}

// GetRawRequest retrieves the raw HTTP request from the request context.
// Use this when you need direct access to headers, cookies, or other request details.
func (rc *RequestContext) GetRawRequest() (*http.Request, error) {
	if rc.RawRequest == nil {
		return nil, fmt.Errorf("raw request not available")
	}
	return rc.RawRequest, nil
}

// ForwardAuthorizationCookie is a RequestEditorFn that forwards the user's session token as a cookie
// Use this for operations that need to authenticate as a specific user
var ForwardAuthorizationCookie client.RequestEditorFn = func(ctx context.Context, req *http.Request) error {
	reqCtx, err := GetRequestContext(ctx)
	if err != nil {
		return err
	}
	userData, err := reqCtx.GetUserData()
	if err != nil {
		return err
	}
	// Add token as cookie for authentication
	req.AddCookie(&http.Cookie{Name: CookieSession, Value: userData.Token})
	return nil
}

// WithAdminToken creates a RequestEditorFn that adds the server/admin token to the Authorization header
// Use this for admin-level operations that regular users cannot perform (counters, user listings, etc.)
// This explicitly removes any session cookies to ensure clean admin authentication
func WithAdminToken(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		// Explicitly remove any session cookies to prevent mixing admin and user auth
		// This ensures clean admin-level authentication
		req.Header.Del("Cookie")
		return nil
	}
}
