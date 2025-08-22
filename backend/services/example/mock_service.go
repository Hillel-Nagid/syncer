package example

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"syncer.net/core/services"
)

// MockService is an example service implementation showing how to use the base service
type MockService struct {
	*services.BaseService[map[string]any]
	clientID     string
	clientSecret string
	redirectURL  string
}

// NewMockService creates a new mock service instance
func NewMockService() *MockService {
	baseService := services.NewBaseService[map[string]any](services.BaseServiceConfig{
		Name:              "mock",
		DisplayName:       "Mock Service",
		Category:          services.CategoryMusic, // Example category
		Scopes:            []string{"read", "write"},
		RequestsPerSecond: 10,
		BurstSize:         15,
		HTTPTimeout:       30 * time.Second,
	})

	return &MockService{
		BaseService:  baseService,
		clientID:     "mock-client-id",
		clientSecret: "mock-client-secret",
		redirectURL:  "http://localhost:8080/auth/callback",
	}
}

// GetAuthURL implements the OAuth authorization URL generation
func (m *MockService) GetAuthURL(state string, redirectURL string) (string, error) {
	m.LogInfo("Generating auth URL for state: %s", state)

	baseURL := "https://api.mockservice.com/oauth/authorize"
	params := url.Values{}
	params.Set("client_id", m.clientID)
	params.Set("redirect_uri", redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "read write")
	params.Set("state", state)

	authURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	return authURL, nil
}

// ExchangeCode implements the OAuth code exchange for tokens
func (m *MockService) ExchangeCode(code string, redirectURL string) (*services.OAuthTokens, error) {
	m.LogInfo("Exchanging code for tokens: %s", code)

	// In a real implementation, you would make an HTTP request to exchange the code
	// For this mock, we return dummy tokens
	tokens := &services.OAuthTokens{
		AccessToken:  "mock-access-token-" + code,
		RefreshToken: "mock-refresh-token-" + code,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour),
		Scope:        "read write",
	}

	return tokens, nil
}

// RefreshTokens implements token refreshing
func (m *MockService) RefreshTokens(refreshToken string) (*services.OAuthTokens, error) {
	m.LogInfo("Refreshing tokens with refresh token: %s", refreshToken)

	// In a real implementation, you would make an HTTP request to refresh tokens
	tokens := &services.OAuthTokens{
		AccessToken:  "mock-refreshed-access-token",
		RefreshToken: refreshToken, // Keep the same refresh token
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour),
		Scope:        "read write",
	}

	return tokens, nil
}

// GetUserProfile implements user profile retrieval
func (m *MockService) GetUserProfile(ctx context.Context, tokens *services.OAuthTokens) (*services.UserProfile, error) {
	m.LogInfo("Getting user profile")

	// Validate tokens first
	valid, err := m.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	// In a real implementation, you would make an API call to get user profile
	profile := &services.UserProfile{
		ExternalID:  "mock-user-123",
		Username:    "mockuser",
		Email:       "mock@example.com",
		DisplayName: "Mock User",
		AvatarURL:   "https://avatar.example.com/mock.jpg",
		Verified:    true,
		Metadata: map[string]any{
			"account_type": "premium",
			"country":      "US",
		},
	}

	return profile, nil
}

// SyncUserData implements data synchronization
func (m *MockService) SyncUserData(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) (*services.SyncResult[map[string]any], error) {
	m.LogInfo("Syncing user data since: %v", lastSync)

	// Validate tokens first
	valid, err := m.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	// Wait for rate limiting
	if err := m.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	// Mock some sync items
	var items []services.SyncItem[map[string]any]

	// Example: Mock some tracks
	for i := 1; i <= 5; i++ {
		item := m.CreateSyncItem(
			fmt.Sprintf("track-%d", i),
			"track",
			services.ActionCreate,
			map[string]any{
				"title":  fmt.Sprintf("Mock Track %d", i),
				"artist": "Mock Artist",
				"album":  "Mock Album",
			},
			time.Now().Add(-time.Duration(i)*time.Hour),
		)
		items = append(items, item)
	}

	// Mock some playlists
	for i := 1; i <= 2; i++ {
		item := m.CreateSyncItem(
			fmt.Sprintf("playlist-%d", i),
			"playlist",
			services.ActionUpdate,
			map[string]any{
				"name":        fmt.Sprintf("Mock Playlist %d", i),
				"description": "A mock playlist",
				"track_count": 10 + i,
			},
			time.Now().Add(-time.Duration(i*2)*time.Hour),
		)
		items = append(items, item)
	}

	// Create sync result
	result := &services.SyncResult[map[string]any]{
		Success:      true,
		ItemsAdded:   m.CountItemsByAction(items, services.ActionCreate),
		ItemsUpdated: m.CountItemsByAction(items, services.ActionUpdate),
		ItemsDeleted: m.CountItemsByAction(items, services.ActionDelete),
		Items:        items,
		Errors:       []services.SyncError{}, // No errors in this mock
		Metadata: map[string]any{
			"sync_time":    time.Now(),
			"service":      "mock",
			"items_synced": len(items),
			"api_version":  "v1",
		},
	}

	m.LogInfo("Sync completed successfully: %d items processed", len(items))
	return result, nil
}

// HealthCheck implements service health checking
func (m *MockService) HealthCheck(ctx context.Context) error {
	m.LogInfo("Performing health check")

	// In a real implementation, you might ping the service API
	// For this mock, we'll always return healthy
	return nil
}
