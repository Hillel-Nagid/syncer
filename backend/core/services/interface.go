package services

import (
	"context"
	"time"
)

// ServiceProvider defines the interface that all service implementations must implement
type ServiceProvider[T any] interface {
	// Metadata
	Name() string
	DisplayName() string
	Category() ServiceCategory
	RequiredScopes() []string

	// OAuth Flow
	GetAuthURL(state string, redirectURL string) (string, error)
	ExchangeCode(code string, redirectURL string) (*OAuthTokens, error)
	RefreshTokens(refreshToken string) (*OAuthTokens, error)
	ValidateTokens(tokens *OAuthTokens) (bool, error)

	// Data Synchronization
	SyncUserData(ctx context.Context, tokens *OAuthTokens, lastSync time.Time) (*SyncResult[T], error)
	GetUserProfile(ctx context.Context, tokens *OAuthTokens) (*UserProfile, error)

	// Health and Status
	HealthCheck(ctx context.Context) error
	GetRateLimit() *RateLimit
}

// ServiceCategory defines the type of service
type ServiceCategory string

const (
	CategoryMusic    ServiceCategory = "music"
	CategoryCalendar ServiceCategory = "calendar"
)

// OAuthTokens represents OAuth2 tokens with expiration
type OAuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope,omitempty"`
}

// SyncResult represents the result of a synchronization operation
type SyncResult[T any] struct {
	Success       bool           `json:"success"`
	ItemsAdded    int            `json:"items_added"`
	ItemsUpdated  int            `json:"items_updated"`
	ItemsDeleted  int            `json:"items_deleted"`
	Items         []SyncItem[T]  `json:"items"`
	NextPageToken string         `json:"next_page_token,omitempty"`
	Errors        []SyncError    `json:"errors,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// SyncItem represents a single item that was synchronized
type SyncItem[T any] struct {
	ExternalID   string     `json:"external_id"`
	ItemType     string     `json:"item_type"`
	Action       SyncAction `json:"action"`
	Data         T          `json:"data"`
	LastModified time.Time  `json:"last_modified"`
	Checksum     string     `json:"checksum,omitempty"`
}

// SyncAction defines what action was performed on an item
type SyncAction string

const (
	ActionCreate SyncAction = "create"
	ActionUpdate SyncAction = "update"
	ActionDelete SyncAction = "delete"
)

// SyncError represents an error that occurred during synchronization
type SyncError struct {
	Type    string `json:"type"`
	Error   string `json:"error"`
	ItemID  string `json:"item_id,omitempty"`
	Context string `json:"context,omitempty"`
}

// UserProfile represents a user profile from an external service
type UserProfile struct {
	ExternalID  string         `json:"external_id"`
	Username    string         `json:"username,omitempty"`
	Email       string         `json:"email,omitempty"`
	DisplayName string         `json:"display_name,omitempty"`
	AvatarURL   string         `json:"avatar_url,omitempty"`
	Verified    bool           `json:"verified,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// RateLimit represents rate limiting information for a service
type RateLimit struct {
	RequestsPerSecond int           `json:"requests_per_second"`
	RequestsPerMinute int           `json:"requests_per_minute"`
	RequestsPerHour   int           `json:"requests_per_hour"`
	BurstSize         int           `json:"burst_size"`
	ResetWindow       time.Duration `json:"reset_window"`
}

// ServiceInfo represents metadata about a service
type ServiceInfo struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Category    ServiceCategory `json:"category"`
	Scopes      []string        `json:"scopes"`
	Available   bool            `json:"available"`
}

// AuthInitiation represents the result of initiating an OAuth flow
type AuthInitiation struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// CallbackResult represents the result of handling an OAuth callback
type CallbackResult struct {
	UserID      string       `json:"user_id"`
	ServiceName string       `json:"service_name"`
	Profile     *UserProfile `json:"profile,omitempty"`
}

// SyncJobRequest represents a request to sync a user's service
type SyncJobRequest struct {
	UserServiceID string       `json:"user_service_id"`
	Priority      SyncPriority `json:"priority"`
	RequestedAt   time.Time    `json:"requested_at"`
}

// SyncPriority defines the priority level for sync jobs
type SyncPriority int

const (
	PriorityLow    SyncPriority = 0
	PriorityMedium SyncPriority = 1
	PriorityHigh   SyncPriority = 2
	PriorityUrgent SyncPriority = 3
)
