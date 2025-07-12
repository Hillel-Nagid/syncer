package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"syncer.net/core/users"
	"syncer.net/utils"
)

type SessionService struct {
	db         *sqlx.DB
	jwtService *JWTService
}

type UserSession struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	SessionToken     string    `json:"session_token" db:"session_token"`
	RefreshToken     string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" db:"refresh_expires_at"`
	IPAddress        string    `json:"ip_address" db:"ip_address"`
	UserAgent        string    `json:"user_agent" db:"user_agent"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	LastUsedAt       time.Time `json:"last_used_at" db:"last_used_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewSessionService(db *sqlx.DB, jwtService *JWTService) *SessionService {
	return &SessionService{
		db:         db,
		jwtService: jwtService,
	}
}

func (s *SessionService) CreateSession(user *users.User, ipAddress, userAgent string) (*AuthTokens, error) {
	sessionToken, err := utils.GenerateToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken()
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateTokenWithSession(user, sessionToken)
	if err != nil {
		return nil, err
	}

	session := &UserSession{
		ID:               uuid.New().String(),
		UserID:           user.ID,
		SessionToken:     sessionToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        time.Now().Add(15 * time.Minute),   // Short-lived access tokens
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		IsActive:         true,
		LastUsedAt:       time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `
		INSERT INTO user_sessions (
			id, user_id, session_token, refresh_token, expires_at, 
			refresh_expires_at, ip_address, user_agent, is_active, 
			last_used_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	_, err = s.db.Exec(query,
		session.ID, session.UserID, session.SessionToken, session.RefreshToken,
		session.ExpiresAt, session.RefreshExpiresAt, session.IPAddress, session.UserAgent,
		session.IsActive, session.LastUsedAt, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}

func (s *SessionService) ValidateSession(sessionToken string) (*UserSession, error) {
	var session UserSession
	query := `
		SELECT * FROM user_sessions 
		WHERE session_token = $1 AND is_active = true AND expires_at > NOW()
	`
	err := s.db.Get(&session, query, sessionToken)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	_, err = s.db.Exec("UPDATE user_sessions SET last_used_at = NOW() WHERE id = $1", session.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) RefreshSession(refreshToken, ipAddress, userAgent string) (*AuthTokens, error) {
	var session UserSession
	query := `
		SELECT * FROM user_sessions 
		WHERE refresh_token = $1 AND is_active = true AND refresh_expires_at > NOW()
	`
	err := s.db.Get(&session, query, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user users.User
	err = s.db.Get(&user, "SELECT * FROM users WHERE id = $1", session.UserID)
	if err != nil {
		return nil, err
	}

	newSessionToken, err := utils.GenerateToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateToken()
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateTokenWithSession(&user, newSessionToken)
	if err != nil {
		return nil, err
	}

	updateQuery := `
		UPDATE user_sessions 
		SET session_token = $1, refresh_token = $2, expires_at = $3, 
			refresh_expires_at = $4, ip_address = $5, user_agent = $6, 
			last_used_at = NOW(), updated_at = NOW()
		WHERE id = $7
	`
	_, err = s.db.Exec(updateQuery,
		newSessionToken, newRefreshToken,
		time.Now().Add(15*time.Minute),
		time.Now().Add(7*24*time.Hour),
		ipAddress, userAgent,
		session.ID)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

func (s *SessionService) RevokeSession(sessionToken string) error {
	_, err := s.db.Exec("UPDATE user_sessions SET is_active = false WHERE session_token = $1", sessionToken)
	return err
}

func (s *SessionService) RevokeAllUserSessions(userID string) error {
	_, err := s.db.Exec("UPDATE user_sessions SET is_active = false WHERE user_id = $1", userID)
	return err
}
