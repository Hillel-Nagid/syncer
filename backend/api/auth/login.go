package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"syncer.net/core/users"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Helper function to set secure cookies
func (s *AuthService) setSecureCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"access_token",
		accessToken,
		15*60,
		"/",
		"",
		true,
		true,
	)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60,
		"/",
		"",
		true,
		true,
	)
}

func (s *AuthService) LoginWithGoogle(c *gin.Context) {
	state := generateRandomState()
	authURL := s.googleOauthConfig.AuthCodeURL(state)

	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

func (s *AuthService) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	storedState, _ := c.Cookie("oauth_state")

	// Log state validation for debugging
	log.Printf("OAuth state validation: received=%s, stored=%s", state, storedState)

	if state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	code := c.Query("code")
	token, err := s.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})
		return
	}

	client := s.googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser users.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Log the Google user info for debugging
	log.Printf("Google user info: ID=%s, Email=%s, Name=%s, VerifiedEmail=%t",
		googleUser.ID, googleUser.Email, googleUser.Name, googleUser.VerifiedEmail)

	user, err := users.CreateUserWithGoogle(s.db, &googleUser)
	if err != nil {
		// Log the actual error for debugging
		log.Printf("Failed to create user with Google: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	user.UpdateLastLogin(s.db)

	tokens, err := s.sessionService.CreateSession(user, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	s.setSecureCookies(c, tokens.AccessToken, tokens.RefreshToken)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	c.Redirect(http.StatusFound, frontendURL)
}

func (s *AuthService) LoginWithEmail(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.VerifyEmailPassword(s.db, req.Email, req.Password)
	if err != nil {
		if err == users.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	user.UpdateLastLogin(s.db)

	tokens, err := s.sessionService.CreateSession(user, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	s.setSecureCookies(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, AuthResponse{
		User: user,
	})
}
