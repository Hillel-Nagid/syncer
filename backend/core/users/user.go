package users

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           string         `json:"id" db:"id"`
	PrimaryEmail string         `json:"primary_email" db:"primary_email"`
	FullName     string         `json:"full_name" db:"full_name"`
	AvatarURL    sql.NullString `json:"avatar_url,omitempty" db:"avatar_url"`
	LastLogin    sql.NullTime   `json:"last_login,omitempty" db:"last_login"`
	Online       bool           `json:"online" db:"online"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

type UserAuthMethod struct {
	ID                  string         `json:"id" db:"id"`
	UserID              string         `json:"user_id" db:"user_id"`
	Provider            Provider       `json:"provider" db:"provider"`
	ProviderID          sql.NullString `json:"provider_id" db:"provider_id"`
	ProviderEmail       sql.NullString `json:"provider_email" db:"provider_email"`
	PasswordHash        sql.NullString `json:"password_hash,omitempty" db:"password_hash"`
	EmailVerified       bool           `json:"email_verified" db:"email_verified"`
	IsPrimary           bool           `json:"is_primary" db:"is_primary"`
	VerificationToken   sql.NullString `json:"verification_token,omitempty" db:"verification_token"`
	VerificationExpires sql.NullTime   `json:"verification_expires,omitempty" db:"verification_expires"`
	ResetToken          sql.NullString `json:"reset_token,omitempty" db:"reset_token"`
	ResetExpires        sql.NullTime   `json:"reset_expires,omitempty" db:"reset_expires"`
	Metadata            []byte         `json:"metadata,omitempty" db:"metadata"`
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
}

// GoogleUserInfo represents the user info from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// UserProfile represents comprehensive user profile information
type UserProfile struct {
	User *User `json:"user"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalSessions     int    `json:"total_sessions"`
	ActiveSessions    int    `json:"active_sessions"`
	LastSessionIP     string `json:"last_session_ip,omitempty"`
	LastSessionAgent  string `json:"last_session_agent,omitempty"`
	ConnectedServices int    `json:"connected_services"`
}

func (u *User) UpdateLastLogin(db *sqlx.DB) error {
	now := time.Now()
	u.LastLogin = sql.NullTime{Time: now, Valid: true}
	_, err := db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", now, u.ID)
	return err
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sqlx.DB, userID string) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their primary email
func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE primary_email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByProviderEmail retrieves a user by email from any auth method
func GetUserByProviderEmail(db *sqlx.DB, email string) (*User, error) {
	var user User
	query := `
		SELECT u.* FROM users u
		JOIN user_auth_methods uam ON u.id = uam.user_id
		WHERE uam.provider_email = $1
		LIMIT 1
	`
	err := db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserAuthMethods retrieves all authentication methods for a user
func GetUserAuthMethods(db *sqlx.DB, userID string) ([]UserAuthMethod, error) {
	var methods []UserAuthMethod
	err := db.Select(&methods, "SELECT * FROM user_auth_methods WHERE user_id = $1", userID)
	return methods, err
}

// GetUserAuthMethod retrieves a specific authentication method for a user
func GetUserAuthMethod(db *sqlx.DB, userID, provider string) (*UserAuthMethod, error) {
	var method UserAuthMethod
	err := db.Get(&method, "SELECT * FROM user_auth_methods WHERE user_id = $1 AND provider = $2", userID, provider)
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// CreateUserAuthMethod creates a new authentication method for a user
func CreateUserAuthMethod(db *sqlx.DB, method *UserAuthMethod) error {
	query := `
		INSERT INTO user_auth_methods (
			user_id, provider, provider_id, provider_email, password_hash, 
			email_verified, is_primary, verification_token, verification_expires, 
			reset_token, reset_expires, metadata
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	return db.Get(method, query,
		method.UserID, method.Provider, method.ProviderID, method.ProviderEmail,
		method.PasswordHash, method.EmailVerified, method.IsPrimary,
		method.VerificationToken, method.VerificationExpires,
		method.ResetToken, method.ResetExpires, method.Metadata)
}

// GetUserAuthMethodByEmail retrieves an authentication method by email
func GetUserAuthMethodByEmail(db *sqlx.DB, email string, provider string) (*UserAuthMethod, error) {
	var method UserAuthMethod
	err := db.Get(&method, "SELECT * FROM user_auth_methods WHERE provider_email = $1 AND provider = $2", email, provider)
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// GetPrimaryAuthMethod retrieves the primary authentication method for a user
func GetPrimaryAuthMethod(db *sqlx.DB, userID string) (*UserAuthMethod, error) {
	var method UserAuthMethod
	err := db.Get(&method, "SELECT * FROM user_auth_methods WHERE user_id = $1 AND is_primary = true", userID)
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// GetUserStats retrieves user statistics
func GetUserStats(db *sqlx.DB, userID string) (*UserStats, error) {
	stats := &UserStats{}

	// Get session statistics
	err := db.Get(&stats.TotalSessions, "SELECT COUNT(*) FROM user_sessions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	err = db.Get(&stats.ActiveSessions, "SELECT COUNT(*) FROM user_sessions WHERE user_id = $1 AND is_active = true AND expires_at > NOW()", userID)
	if err != nil {
		return nil, err
	}

	// Get last session info
	var lastSession struct {
		IPAddress sql.NullString `db:"ip_address"`
		UserAgent sql.NullString `db:"user_agent"`
	}
	err = db.Get(&lastSession, `
		SELECT ip_address, user_agent 
		FROM user_sessions 
		WHERE user_id = $1 
		ORDER BY last_used_at DESC 
		LIMIT 1
	`, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if lastSession.IPAddress.Valid {
		stats.LastSessionIP = lastSession.IPAddress.String
	}
	if lastSession.UserAgent.Valid {
		stats.LastSessionAgent = lastSession.UserAgent.String
	}

	// Get connected services count
	err = db.Get(&stats.ConnectedServices, "SELECT COUNT(*) FROM user_services WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
