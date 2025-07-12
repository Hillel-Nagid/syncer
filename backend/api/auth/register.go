package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"syncer.net/core/email"
	"syncer.net/core/users"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthResponse struct {
	User *users.User `json:"user"`
}

func (s *AuthService) RegisterWithEmail(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, verificationToken, err := users.CreateUserWithEmail(s.db, req.Email, req.Password, req.FullName)
	if err != nil {
		if err == users.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create session with IP and User Agent
	tokens, err := s.sessionService.CreateSession(user, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Send verification email
	ctx := context.Background()
	emailData := email.VerificationEmailData{
		Email:             user.PrimaryEmail,
		FullName:          user.FullName,
		VerificationToken: verificationToken,
	}

	err = s.emailService.SendVerificationEmail(ctx, emailData)
	if err != nil {
		// Log the error but don't fail the registration
		log.Printf("Failed to send verification email to %s: %v", user.PrimaryEmail, err)
		// You might want to add a flag to indicate the email couldn't be sent
	}

	// Set secure cookies instead of returning tokens in response
	s.setSecureCookies(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusCreated, AuthResponse{
		User: user,
	})
}

func (s *AuthService) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := users.VerifyEmail(s.db, req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (s *AuthService) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user information first
	user, err := users.GetUserByEmail(s.db, req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	verificationToken, err := users.ResendVerificationEmail(s.db, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	emailData := email.VerificationEmailData{
		Email:             user.PrimaryEmail,
		FullName:          user.FullName,
		VerificationToken: verificationToken,
	}

	err = s.emailService.SendVerificationEmail(ctx, emailData)
	if err != nil {
		log.Printf("Failed to resend verification email to %s: %v", user.PrimaryEmail, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func generateRandomState() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
