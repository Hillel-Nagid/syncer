package middlewares

import (
	"crypto/subtle"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"syncer.net/core/auth"
	"syncer.net/utils"
)

func AuthMiddleware(jwtService *auth.JWTService, sessionService *auth.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
				c.Abort()
				return
			}
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Validate session if session token exists
		if claims.SessionToken != "" {
			_, err := sessionService.ValidateSession(claims.SessionToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("session_token", claims.SessionToken)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"https://syncer.net",
		}

		if slices.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		c.Writer.Header().Set("X-Frame-Options", "DENY")

		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self' https://accounts.google.com; " +
			"font-src 'self' data:; " +
			"object-src 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"frame-ancestors 'none'; " +
			"upgrade-insecure-requests"
		c.Writer.Header().Set("Content-Security-Policy", csp)

		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		c.Writer.Header().Set("Referrer-Policy", "no-referrer")

		if c.Request.URL.Path == "/api/profile" || c.Request.URL.Path == "/auth/login" {
			c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Writer.Header().Set("Pragma", "no-cache")
			c.Writer.Header().Set("Expires", "0")
		}

		c.Next()
	}
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		csrfCookie, err := c.Cookie("csrf_token")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}

		csrfHeader := c.GetHeader("X-CSRF-Token")
		if csrfHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token required in header"})
			c.Abort()
			return
		}

		if !validateCSRFToken(csrfCookie, csrfHeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CSRFTokenEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("csrf_token", token, 60*60, "/", "", true, false)

		c.JSON(http.StatusOK, gin.H{"csrf_token": token})
	}
}

func validateCSRFToken(cookieToken, headerToken string) bool {
	return subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) == 1
}
