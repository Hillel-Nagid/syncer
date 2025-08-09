package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"syncer.net/core/users"
)

type ProfileResponse struct {
	ID           string    `json:"id" db:"id"`
	PrimaryEmail string    `json:"primary_email" db:"primary_email"`
	FullName     string    `json:"full_name" db:"full_name"`
	AvatarURL    string    `json:"avatar_url,omitempty" db:"avatar_url"`
	LastLogin    time.Time `json:"last_login,omitempty" db:"last_login"`
	Online       bool      `json:"online" db:"online"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (s *AuthService) clearSecureCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}

func (s *AuthService) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found"})
		return
	}

	tokens, err := s.sessionService.RefreshSession(refreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		s.clearSecureCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	s.setSecureCookies(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

func (s *AuthService) Logout(c *gin.Context) {
	sessionToken := c.GetString("session_token")
	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active session"})
		return
	}

	err := s.sessionService.RevokeSession(sessionToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	s.clearSecureCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (s *AuthService) LogoutAll(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}

	err := s.sessionService.RevokeAllUserSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout from all devices"})
		return
	}

	s.clearSecureCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices successfully"})
}

func (s *AuthService) GetProfileBasicInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	db := c.MustGet("db").(*sqlx.DB)

	profile, err := users.GetUserByID(db, userID)
	if err != nil {
		log.Printf("Error retrieving user profile for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, ProfileResponse{
		ID:           profile.ID,
		PrimaryEmail: profile.PrimaryEmail,
		FullName:     profile.FullName,
		AvatarURL:    profile.AvatarURL.String,
		LastLogin:    profile.LastLogin.Time,
		Online:       profile.Online,
		CreatedAt:    profile.CreatedAt,
		UpdatedAt:    profile.UpdatedAt,
	})
}
