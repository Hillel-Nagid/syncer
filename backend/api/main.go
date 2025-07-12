package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"syncer.net/api/auth"
	"syncer.net/api/middlewares"
	coreAuth "syncer.net/core/auth"
	"syncer.net/core/email"
)

func main() {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	jwtService := coreAuth.NewJWTService(jwtSecret)

	sessionService := coreAuth.NewSessionService(db, jwtService)

	emailService := email.NewEmailService()

	googleOauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	router := gin.Default()
	router.Use(middlewares.DBMiddleware(db))
	router.Use(middlewares.SecurityHeadersMiddleware())
	router.Use(middlewares.CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// CSRF token endpoint
	router.GET("/csrf-token", middlewares.CSRFTokenEndpoint())

	authRoutes := router.Group("/auth")
	{
		authService := auth.NewAuthService(db, googleOauthConfig, jwtService, sessionService, emailService)

		// Apply CSRF protection to state-changing auth routes
		authRoutes.POST("/register", middlewares.CSRFMiddleware(), authService.RegisterWithEmail)
		authRoutes.POST("/login", middlewares.CSRFMiddleware(), authService.LoginWithEmail)
		authRoutes.POST("/verify-email", middlewares.CSRFMiddleware(), authService.VerifyEmail)
		authRoutes.POST("/resend-verification", middlewares.CSRFMiddleware(), authService.ResendVerificationEmail)
		authRoutes.POST("/forgot-password", middlewares.CSRFMiddleware(), authService.ForgotPassword)
		authRoutes.POST("/reset-password", middlewares.CSRFMiddleware(), authService.ResetPassword)
		authRoutes.GET("/validate-reset-token", authService.ValidateResetToken)
		authRoutes.GET("/google", authService.LoginWithGoogle)
		authRoutes.GET("/google/callback", authService.GoogleCallback)
		authRoutes.POST("/refresh", middlewares.CSRFMiddleware(), authService.RefreshToken)

		protectedAuthRoutes := authRoutes.Group("/")
		protectedAuthRoutes.Use(middlewares.AuthMiddleware(jwtService, sessionService))
		protectedAuthRoutes.Use(middlewares.CSRFMiddleware())
		{
			protectedAuthRoutes.POST("/logout", authService.Logout)
			protectedAuthRoutes.POST("/logout-all", authService.LogoutAll)
		}
	}

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware(jwtService, sessionService))
	protectedRoutes.Use(middlewares.CSRFMiddleware())
	{
		protectedRoutes.GET("/profile", func(c *gin.Context) {
			userID := c.GetString("user_id")
			userEmail := c.GetString("user_email")
			c.JSON(200, gin.H{
				"user_id": userID,
				"email":   userEmail,
			})
		})
	}

	serviceRoutes := protectedRoutes.Group("/services")
	{
		serviceRoutes.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Services endpoint"})
		})
	}

	router.POST("/contact", middlewares.CSRFMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Contact endpoint"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
