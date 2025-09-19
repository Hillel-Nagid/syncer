package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// BaseService provides common functionality for all service implementations
type BaseService struct {
	name        string
	displayName string
	category    ServiceCategory
	scopes      []string
	rateLimiter *rate.Limiter
	httpClient  *http.Client
	logger      *log.Logger
}

// BaseServiceConfig contains configuration for creating a base service
type BaseServiceConfig struct {
	Name              string
	DisplayName       string
	Category          ServiceCategory
	Scopes            []string
	RequestsPerSecond int
	BurstSize         int
	HTTPTimeout       time.Duration
	Logger            *log.Logger
}

// NewBaseService creates a new base service with the given configuration
func NewBaseService(config BaseServiceConfig) *BaseService {
	if config.Logger == nil {
		config.Logger = log.New(log.Writer(), fmt.Sprintf("[%s] ", config.Name), log.LstdFlags)
	}

	if config.RequestsPerSecond == 0 {
		config.RequestsPerSecond = 5
	}

	if config.BurstSize == 0 {
		config.BurstSize = config.RequestsPerSecond
	}

	if config.HTTPTimeout == 0 {
		config.HTTPTimeout = 30 * time.Second
	}

	return &BaseService{
		name:        config.Name,
		displayName: config.DisplayName,
		category:    config.Category,
		scopes:      config.Scopes,
		rateLimiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(config.RequestsPerSecond)), config.BurstSize),
		httpClient: &http.Client{
			Timeout: config.HTTPTimeout,
		},
		logger: config.Logger,
	}
}

func (b *BaseService) Name() string {
	return b.name
}

func (b *BaseService) DisplayName() string {
	return b.displayName
}

func (b *BaseService) Category() ServiceCategory {
	return b.category
}

func (b *BaseService) RequiredScopes() []string {
	return b.scopes
}

func (b *BaseService) GetRateLimit() *RateLimit {
	return &RateLimit{
		RequestsPerSecond: int(b.rateLimiter.Limit()),
		BurstSize:         b.rateLimiter.Burst(),
		ResetWindow:       time.Second,
	}
}

func (b *BaseService) HealthCheck(ctx context.Context) error {
	//TODO: Implement health check
	return nil
}

// WaitForRateLimit waits for rate limiting if necessary
func (b *BaseService) WaitForRateLimit(ctx context.Context) error {
	return b.rateLimiter.Wait(ctx)
}

// CreateAuthenticatedRequest creates an HTTP request with OAuth authorization
func (b *BaseService) CreateAuthenticatedRequest(ctx context.Context, method, url string, tokens *OAuthTokens) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if tokens.TokenType == "" {
		tokens.TokenType = "Bearer"
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	req.Header.Set("User-Agent", fmt.Sprintf("Syncer/1.0 (%s)", b.name))
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// DoRequest performs an HTTP request with rate limiting
func (b *BaseService) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := b.WaitForRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	b.logger.Printf("HTTP %s %s -> %d", req.Method, req.URL.String(), resp.StatusCode)

	return resp, nil
}

// ValidateTokens provides a basic token validation
func (b *BaseService) ValidateTokens(tokens *OAuthTokens) (bool, error) {
	if tokens == nil {
		return false, fmt.Errorf("tokens cannot be nil")
	}

	if tokens.AccessToken == "" {
		return false, fmt.Errorf("access token is required")
	}

	if !tokens.ExpiresAt.IsZero() && time.Now().After(tokens.ExpiresAt) {
		return false, fmt.Errorf("access token is expired")
	}

	return true, nil
}

// CountItemsByAction is a utility function to count sync items by action
func (b *BaseService) CountItemsByAction(items []SyncItem, action SyncAction) int {
	//TODO: Look closer into
	count := 0
	for _, item := range items {
		if item.Action == action {
			count++
		}
	}
	return count
}

// CreateSyncItem is a helper to create a sync item
func (b *BaseService) CreateSyncItem(externalID, itemType string, action SyncAction, data any) SyncItem {
	return SyncItem{
		ExternalID: externalID,
		ItemType:   itemType,
		Action:     action,
		Data:       data,
	}
}

// CreateSyncError is a helper to create a sync error
func (b *BaseService) CreateSyncError(errorType, message, itemID, context string) SyncError {
	return SyncError{
		Type:    errorType,
		Error:   message,
		ItemID:  itemID,
		Context: context,
	}
}

// LogInfo logs an info message
func (b *BaseService) LogInfo(message string, args ...any) {
	b.logger.Printf("[INFO] "+message, args...)
}

// LogWarn logs a warning message
func (b *BaseService) LogWarn(message string, args ...any) {
	b.logger.Printf("[WARN] "+message, args...)
}

// LogError logs an error message
func (b *BaseService) LogError(message string, args ...any) {
	b.logger.Printf("[ERROR] "+message, args...)
}

func (b *BaseService) GetAuthURL(state string, redirectURL string) (string, error) {
	return "", fmt.Errorf("GetAuthURL not implemented for service %s", b.name)
}

func (b *BaseService) ExchangeCode(code string, redirectURL string) (*OAuthTokens, error) {
	return nil, fmt.Errorf("ExchangeCode not implemented for service %s", b.name)
}

func (b *BaseService) RefreshTokens(refreshToken string) (*OAuthTokens, error) {
	return nil, fmt.Errorf("RefreshTokens not implemented for service %s", b.name)
}

func (b *BaseService) GetUserData(ctx context.Context, tokens *OAuthTokens, lastSync time.Time) (*UserDataResult, error) {
	return nil, fmt.Errorf("GetUserData not implemented for service %s", b.name)
}

func (b *BaseService) GetUserProfile(ctx context.Context, tokens *OAuthTokens) (*UserProfile, error) {
	return nil, fmt.Errorf("GetUserProfile not implemented for service %s", b.name)
}
