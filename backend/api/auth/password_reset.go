package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"syncer.net/core/users"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// ForgotPassword initiates a password reset by sending a reset email
func (s *AuthService) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByEmail(s.db, req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, we've sent a password reset link"})
		return
	}

	resetToken, err := users.RequestPasswordReset(s.db, req.Email)
	if err != nil {
		log.Printf("Failed to generate reset token for %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password reset request"})
		return
	}

	ctx := context.Background()
	err = s.emailService.SendPasswordResetEmail(ctx, user.PrimaryEmail, user.FullName, resetToken)
	if err != nil {
		log.Printf("Failed to send password reset email to %s: %v", user.PrimaryEmail, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent successfully"})
}

// ResetPassword validates a reset token and updates the user's password
func (s *AuthService) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := users.ResetPassword(s.db, req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ValidateResetToken validates a reset token without updating the password (for frontend validation)
func (s *AuthService) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := users.GetUserByResetToken(s.db, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"email":   user.PrimaryEmail,
	})
}
