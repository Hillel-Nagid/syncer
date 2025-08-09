package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"syncer.net/utils"
)

type Provider string

const (
	EmailProvider  Provider = "email"
	GoogleProvider Provider = "google"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// CreateUserWithEmail creates a new user with email authentication
func CreateUserWithEmail(db *sqlx.DB, email, password, fullName string) (*User, string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	var existingMethod UserAuthMethod
	err = tx.Get(&existingMethod, "SELECT id FROM user_auth_methods WHERE provider_email = $1", email)
	if err == nil {
		return nil, "", ErrUserAlreadyExists
	}
	if err != sql.ErrNoRows {
		return nil, "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	verificationToken, err := utils.GenerateToken()
	if err != nil {
		return nil, "", err
	}

	user := &User{
		ID:           uuid.New().String(),
		PrimaryEmail: email,
		FullName:     fullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = tx.NamedExec(`
		INSERT INTO users (id, primary_email, full_name, created_at, updated_at)
		VALUES (:id, :primary_email, :full_name, :created_at, :updated_at)
	`, user)
	if err != nil {
		return nil, "", err
	}

	authMethod := &UserAuthMethod{
		ID:                  uuid.New().String(),
		UserID:              user.ID,
		Provider:            EmailProvider,
		ProviderEmail:       sql.NullString{String: email, Valid: true},
		PasswordHash:        sql.NullString{String: string(hashedPassword), Valid: true},
		EmailVerified:       false,
		IsPrimary:           true,
		VerificationToken:   sql.NullString{String: verificationToken, Valid: true},
		VerificationExpires: sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true},
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	_, err = tx.NamedExec(`
		INSERT INTO user_auth_methods (
			id, user_id, provider, provider_email, password_hash, 
			email_verified, is_primary, verification_token, verification_expires,
			created_at, updated_at
		) VALUES (
			:id, :user_id, :provider, :provider_email, :password_hash,
			:email_verified, :is_primary, :verification_token, :verification_expires,
			:created_at, :updated_at
		)
	`, authMethod)
	if err != nil {
		return nil, "", err
	}

	if err = tx.Commit(); err != nil {
		return nil, "", err
	}

	return user, verificationToken, nil
}

// CreateUserWithGoogle creates a new user with Google OAuth or links to existing user
func CreateUserWithGoogle(db *sqlx.DB, googleInfo *GoogleUserInfo) (*User, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var existingMethod UserAuthMethod
	err = tx.Get(&existingMethod, "SELECT * FROM user_auth_methods WHERE provider = 'google' AND provider_id = $1", googleInfo.ID)
	if err == nil {
		var user User
		err = tx.Get(&user, "SELECT * FROM users WHERE id = $1", existingMethod.UserID)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	var existingUser User
	err = tx.Get(&existingUser, "SELECT u.* FROM users u JOIN user_auth_methods uam ON u.id = uam.user_id WHERE uam.provider_email = $1", googleInfo.Email)
	if err == nil {
		metadata, _ := json.Marshal(googleInfo)

		authMethod := &UserAuthMethod{
			ID:            uuid.New().String(),
			UserID:        existingUser.ID,
			Provider:      GoogleProvider,
			ProviderID:    sql.NullString{String: googleInfo.ID, Valid: true},
			ProviderEmail: sql.NullString{String: googleInfo.Email, Valid: true},
			EmailVerified: googleInfo.VerifiedEmail,
			IsPrimary:     false,
			Metadata:      metadata,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err = tx.NamedExec(`
			INSERT INTO user_auth_methods (
				id, user_id, provider, provider_id, provider_email,
				email_verified, is_primary, metadata, created_at, updated_at
			) VALUES (
				:id, :user_id, :provider, :provider_id, :provider_email,
				:email_verified, :is_primary, :metadata, :created_at, :updated_at
			)
		`, authMethod)
		if err != nil {
			return nil, err
		}

		if err = tx.Commit(); err != nil {
			return nil, err
		}
		return &existingUser, nil
	}

	user := &User{
		ID:           uuid.New().String(),
		PrimaryEmail: googleInfo.Email,
		FullName:     googleInfo.Name,
		AvatarURL:    sql.NullString{String: googleInfo.Picture, Valid: true},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = tx.NamedExec(`
		INSERT INTO users (id, primary_email, full_name, avatar_url, created_at, updated_at)
		VALUES (:id, :primary_email, :full_name, :avatar_url, :created_at, :updated_at)
	`, user)
	if err != nil {
		return nil, err
	}

	metadata, _ := json.Marshal(googleInfo)
	authMethod := &UserAuthMethod{
		ID:            uuid.New().String(),
		UserID:        user.ID,
		Provider:      GoogleProvider,
		ProviderID:    sql.NullString{String: googleInfo.ID, Valid: true},
		ProviderEmail: sql.NullString{String: googleInfo.Email, Valid: true},
		EmailVerified: googleInfo.VerifiedEmail,
		IsPrimary:     true,
		Metadata:      metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = tx.NamedExec(`
		INSERT INTO user_auth_methods (
			id, user_id, provider, provider_id, provider_email,
			email_verified, is_primary, metadata, created_at, updated_at
		) VALUES (
			:id, :user_id, :provider, :provider_id, :provider_email,
			:email_verified, :is_primary, :metadata, :created_at, :updated_at
		)
	`, authMethod)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

// VerifyEmailPassword verifies email and password for login
func VerifyEmailPassword(db *sqlx.DB, email, password string) (*User, error) {
	var authMethod UserAuthMethod
	err := db.Get(&authMethod, `
		SELECT * FROM user_auth_methods 
		WHERE provider = 'email' AND provider_email = $1
	`, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !authMethod.PasswordHash.Valid {
		return nil, ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(authMethod.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	var user User
	err = db.Get(&user, "SELECT * FROM users WHERE id = $1", authMethod.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// VerifyEmail verifies a user's email address using the verification token
func VerifyEmail(db *sqlx.DB, token string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var authMethod UserAuthMethod
	err = tx.Get(&authMethod, `
		SELECT * FROM user_auth_methods 
		WHERE verification_token = $1 AND provider = 'email'
	`, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invalid verification token")
		}
		return err
	}

	// Check if token has expired
	if !authMethod.VerificationExpires.Valid || time.Now().After(authMethod.VerificationExpires.Time) {
		return errors.New("verification token has expired")
	}

	// Check if already verified
	if authMethod.EmailVerified {
		return errors.New("email already verified")
	}

	// Update the auth method to mark email as verified
	_, err = tx.Exec(`
		UPDATE user_auth_methods 
		SET email_verified = true, verification_token = NULL, verification_expires = NULL, updated_at = NOW()
		WHERE id = $1
	`, authMethod.ID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ResendVerificationEmail generates a new verification token for a user
func ResendVerificationEmail(db *sqlx.DB, email string) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var authMethod UserAuthMethod
	err = tx.Get(&authMethod, `
		SELECT * FROM user_auth_methods 
		WHERE provider_email = $1 AND provider = 'email'
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// Check if already verified
	if authMethod.EmailVerified {
		return "", errors.New("email already verified")
	}

	// Generate new verification token
	newToken, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	// Update the auth method with new token
	_, err = tx.Exec(`
		UPDATE user_auth_methods 
		SET verification_token = $1, verification_expires = $2, updated_at = NOW()
		WHERE id = $3
	`, newToken, time.Now().Add(24*time.Hour), authMethod.ID)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return newToken, nil
}

// RequestPasswordReset generates a password reset token for a user
func RequestPasswordReset(db *sqlx.DB, email string) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var authMethod UserAuthMethod
	err = tx.Get(&authMethod, `
		SELECT * FROM user_auth_methods 
		WHERE provider_email = $1 AND provider = 'email'
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	resetToken, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
		UPDATE user_auth_methods 
		SET reset_token = $1, reset_expires = $2, updated_at = $3
		WHERE id = $4
	`, resetToken, time.Now().Add(1*time.Hour), time.Now(), authMethod.ID)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return resetToken, nil
}

// ResetPassword validates a reset token and updates the user's password
func ResetPassword(db *sqlx.DB, token, newPassword string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var authMethod UserAuthMethod
	err = tx.Get(&authMethod, `
		SELECT * FROM user_auth_methods 
		WHERE reset_token = $1 AND provider = 'email'
	`, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invalid reset token")
		}
		return err
	}

	if !authMethod.ResetExpires.Valid || time.Now().After(authMethod.ResetExpires.Time) {
		return errors.New("reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE user_auth_methods 
		SET password_hash = $1, reset_token = NULL, reset_expires = NULL, updated_at = NOW()
		WHERE id = $2
	`, string(hashedPassword), authMethod.ID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetUserByResetToken retrieves a user by their reset token (for getting user info before reset)
func GetUserByResetToken(db *sqlx.DB, token string) (*User, error) {
	var user User
	query := `
		SELECT u.* FROM users u
		JOIN user_auth_methods uam ON u.id = uam.user_id
		WHERE uam.reset_token = $1 AND uam.provider = 'email' AND uam.reset_expires > $2
	`
	err := db.Get(&user, query, token, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}
	return &user, nil
}
