package services

import "time"

type Service struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserService struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ServiceID    string    `json:"service_id" db:"service_id"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	Metadata     []byte    `json:"metadata" db:"metadata"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
