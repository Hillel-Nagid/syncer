package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"syncer.net/core/security"
)

// OAuthManager handles OAuth flows and token management
type OAuthManager struct {
	registry   *ServiceRegistry
	db         *sqlx.DB
	encryption *security.TokenEncryption
	logger     *log.Logger
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(registry *ServiceRegistry, db *sqlx.DB, encryptionKey []byte, logger *log.Logger) (*OAuthManager, error) {
	if logger == nil {
		logger = log.New(log.Writer(), "[OAuthManager] ", log.LstdFlags)
	}

	encryption, err := security.NewTokenEncryption(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token encryption: %w", err)
	}

	return &OAuthManager{
		registry:   registry,
		db:         db,
		encryption: encryption,
		logger:     logger,
	}, nil
}

// PendingAuth represents a pending OAuth authorization
type PendingAuth struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	ServiceName string    `db:"service_name"`
	State       string    `db:"state"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}

// InitiateAuth starts the OAuth flow for a service
func (o *OAuthManager) InitiateAuth(serviceName, userID, redirectURL string) (*AuthInitiation, error) {
	service, err := o.registry.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	// Generate secure state token
	state, err := o.generateStateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state token: %w", err)
	}

	// Get authorization URL
	authURL, err := service.GetAuthURL(state, redirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth URL: %w", err)
	}

	// Store pending authorization
	err = o.storePendingAuth(userID, serviceName, state, time.Now().Add(10*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("failed to store pending auth: %w", err)
	}

	o.logger.Printf("Initiated OAuth flow for user %s, service %s", userID, serviceName)

	return &AuthInitiation{
		AuthURL: authURL,
		State:   state,
	}, nil
}

// HandleCallback processes the OAuth callback
func (o *OAuthManager) HandleCallback(serviceName, code, state string) (*CallbackResult, error) {
	// Validate state token
	userID, err := o.validateStateToken(serviceName, state)
	if err != nil {
		return nil, fmt.Errorf("invalid state token: %w", err)
	}

	service, err := o.registry.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	// Exchange code for tokens
	tokens, err := service.ExchangeCode(code, "")
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user profile from service
	profile, err := service.GetUserProfile(context.Background(), tokens)
	if err != nil {
		o.logger.Printf("Warning: Failed to get user profile for %s: %v", serviceName, err)
	}

	// Store user tokens
	err = o.storeUserTokens(userID, serviceName, tokens, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to store tokens: %w", err)
	}

	// Clean up pending auth
	o.cleanupPendingAuth(state)

	o.logger.Printf("Successfully completed OAuth flow for user %s, service %s", userID, serviceName)

	return &CallbackResult{
		UserID:      userID,
		ServiceName: serviceName,
		Profile:     profile,
	}, nil
}

// RefreshUserTokens refreshes expired tokens for a user service
func (o *OAuthManager) RefreshUserTokens(userServiceID string) error {
	// Get current tokens
	userService, err := o.getUserService(userServiceID)
	if err != nil {
		return fmt.Errorf("failed to get user service: %w", err)
	}

	// Get service provider
	service, err := o.registry.GetService(userService.ServiceName)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	// Decrypt refresh token
	_, refreshToken, err := o.encryption.DecryptTokens(nil, userService.EncryptedRefreshToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	if refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Refresh tokens
	newTokens, err := service.RefreshTokens(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh tokens: %w", err)
	}

	// Update stored tokens
	err = o.updateUserTokens(userServiceID, newTokens)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}

	o.logger.Printf("Successfully refreshed tokens for user service %s", userServiceID)
	return nil
}

// GetUserTokens retrieves and decrypts tokens for a user service
func (o *OAuthManager) GetUserTokens(userServiceID string) (*OAuthTokens, error) {
	userService, err := o.getUserService(userServiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user service: %w", err)
	}

	// Decrypt tokens
	accessToken, refreshToken, err := o.encryption.DecryptTokens(
		userService.EncryptedAccessToken,
		userService.EncryptedRefreshToken,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt tokens: %w", err)
	}

	return &OAuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    userService.TokenType,
		ExpiresAt:    userService.TokenExpiresAt,
		Scope:        userService.Scopes,
	}, nil
}

// generateStateToken creates a cryptographically secure state token
func (o *OAuthManager) generateStateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// storePendingAuth stores a pending OAuth authorization
func (o *OAuthManager) storePendingAuth(userID, serviceName, state string, expiresAt time.Time) error {
	_, err := o.db.Exec(`
		INSERT INTO pending_oauth_auth (id, user_id, service_name, state, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, service_name) 
		DO UPDATE SET state = $4, expires_at = $5, created_at = NOW()
	`, uuid.New().String(), userID, serviceName, state, expiresAt)

	return err
}

// validateStateToken validates the state token and returns the user ID
func (o *OAuthManager) validateStateToken(serviceName, state string) (string, error) {
	var auth PendingAuth
	err := o.db.Get(&auth, `
		SELECT user_id, service_name, expires_at 
		FROM pending_oauth_auth 
		WHERE service_name = $1 AND state = $2 AND expires_at > NOW()
	`, serviceName, state)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid or expired state token")
		}
		return "", err
	}

	return auth.UserID, nil
}

// cleanupPendingAuth removes a pending authorization
func (o *OAuthManager) cleanupPendingAuth(state string) {
	_, err := o.db.Exec("DELETE FROM pending_oauth_auth WHERE state = $1", state)
	if err != nil {
		o.logger.Printf("Warning: Failed to cleanup pending auth: %v", err)
	}
}

// UserServiceRecord represents the enhanced user service record
type UserServiceRecord struct {
	ID                    string     `db:"id"`
	UserID                string     `db:"user_id"`
	ServiceID             string     `db:"service_id"`
	ServiceName           string     `db:"service_name"`
	EncryptedAccessToken  []byte     `db:"encrypted_access_token"`
	EncryptedRefreshToken []byte     `db:"encrypted_refresh_token"`
	TokenType             string     `db:"token_type"`
	TokenExpiresAt        time.Time  `db:"token_expires_at"`
	LastSyncAt            *time.Time `db:"last_sync_at"`
	SyncFrequency         string     `db:"sync_frequency"`
	SyncEnabled           bool       `db:"sync_enabled"`
	ServiceUserID         string     `db:"service_user_id"`
	ServiceUsername       string     `db:"service_username"`
	Scopes                string     `db:"scopes"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

// storeUserTokens encrypts and stores OAuth tokens
func (o *OAuthManager) storeUserTokens(userID, serviceName string, tokens *OAuthTokens, profile *UserProfile) error {
	// Get service ID
	var serviceID string
	err := o.db.Get(&serviceID, "SELECT id FROM services WHERE name = $1", serviceName)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	// Encrypt tokens
	encryptedAccess, encryptedRefresh, err := o.encryption.EncryptTokens(tokens.AccessToken, tokens.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt tokens: %w", err)
	}

	// Prepare profile data
	serviceUserID := ""
	serviceUsername := ""
	if profile != nil {
		serviceUserID = profile.ExternalID
		serviceUsername = profile.Username
	}

	// Store or update user service
	_, err = o.db.Exec(`
		INSERT INTO user_services (
			id, user_id, service_id, encrypted_access_token, encrypted_refresh_token,
			token_type, token_expires_at, service_user_id, service_username, scopes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id, service_id) 
		DO UPDATE SET 
			encrypted_access_token = $4,
			encrypted_refresh_token = $5,
			token_type = $6,
			token_expires_at = $7,
			service_user_id = $8,
			service_username = $9,
			scopes = $10,
			updated_at = NOW()
	`, uuid.New().String(), userID, serviceID, encryptedAccess, encryptedRefresh,
		tokens.TokenType, tokens.ExpiresAt, serviceUserID, serviceUsername, tokens.Scope)

	return err
}

// getUserService retrieves a user service record
func (o *OAuthManager) getUserService(userServiceID string) (*UserServiceRecord, error) {
	var userService UserServiceRecord
	err := o.db.Get(&userService, `
		SELECT us.*, s.name as service_name
		FROM user_services us
		JOIN services s ON us.service_id = s.id
		WHERE us.id = $1
	`, userServiceID)

	if err != nil {
		return nil, err
	}

	return &userService, nil
}

// updateUserTokens updates encrypted tokens for a user service
func (o *OAuthManager) updateUserTokens(userServiceID string, tokens *OAuthTokens) error {
	// Encrypt tokens
	encryptedAccess, encryptedRefresh, err := o.encryption.EncryptTokens(tokens.AccessToken, tokens.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt tokens: %w", err)
	}

	_, err = o.db.Exec(`
		UPDATE user_services 
		SET encrypted_access_token = $1, encrypted_refresh_token = $2, 
		    token_type = $3, token_expires_at = $4, scopes = $5, updated_at = NOW()
		WHERE id = $6
	`, encryptedAccess, encryptedRefresh, tokens.TokenType, tokens.ExpiresAt, tokens.Scope, userServiceID)

	return err
}
