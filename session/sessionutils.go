package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

const HTTPRequestContextKey contextKey = "httpRequest"

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

func GetUserId(ctx context.Context) (string, error) {
	req, _ := ctx.Value(HTTPRequestContextKey).(*http.Request)
	if req == nil {
		return "", fmt.Errorf("request context is missing")
	}

	userID := req.Header.Get(HeaderUserId)
	if userID == "" {
		log.Printf("Error: X-User-ID header is missing")
		return "", fmt.Errorf("missing user header")
	}

	return userID, nil
}

func GetUserData(ctx context.Context) (Session, error) {
	req, _ := ctx.Value(HTTPRequestContextKey).(*http.Request)
	if req == nil {
		return Session{}, fmt.Errorf("request context is missing")
	}

	userData := req.Header.Get(HeaderUserData)
	if userData == "" {
		log.Printf("Error: X-User-Data header is missing")
		return Session{}, fmt.Errorf("missing user header")
	}

	var session Session
	err := json.Unmarshal([]byte(userData), &session)
	if err != nil {
		log.Printf("Error: Failed to unmarshal user data: %v", err)
		return Session{}, fmt.Errorf("invalid user data")
	}

	return session, nil
}

// ForwardAuthorizationCookie is a RequestEditorFn that forwards the user's session token as a cookie
// Use this for operations that need to authenticate as a specific user
var ForwardAuthorizationCookie client.RequestEditorFn = func(ctx context.Context, req *http.Request) error {
	userData, err := GetUserData(ctx)
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
