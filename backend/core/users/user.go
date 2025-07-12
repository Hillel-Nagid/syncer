package users

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	PrimaryEmail string    `json:"primary_email" db:"primary_email"`
	FullName     string    `json:"full_name" db:"full_name"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
	Online       bool      `json:"online" db:"online"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserAuthMethod struct {
	ID                  string    `json:"id" db:"id"`
	UserID              string    `json:"user_id" db:"user_id"`
	Provider            Provider  `json:"provider" db:"provider"`
	ProviderID          string    `json:"provider_id" db:"provider_id"`
	ProviderEmail       string    `json:"provider_email" db:"provider_email"`
	PasswordHash        string    `json:"password_hash,omitempty" db:"password_hash"`
	EmailVerified       bool      `json:"email_verified" db:"email_verified"`
	IsPrimary           bool      `json:"is_primary" db:"is_primary"`
	VerificationToken   string    `json:"verification_token,omitempty" db:"verification_token"`
	VerificationExpires time.Time `json:"verification_expires,omitempty" db:"verification_expires"`
	ResetToken          string    `json:"reset_token,omitempty" db:"reset_token"`
	ResetExpires        time.Time `json:"reset_expires,omitempty" db:"reset_expires"`
	Metadata            []byte    `json:"metadata,omitempty" db:"metadata"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// GoogleUserInfo represents the user info from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (u *User) UpdateLastLogin(db *sqlx.DB) error {
	u.LastLogin = time.Now()
	_, err := db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", u.LastLogin, u.ID)
	return err
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
