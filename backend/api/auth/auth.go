package auth

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
	"syncer.net/core/auth"
	"syncer.net/core/email"
)

type AuthService struct {
	db                *sqlx.DB
	googleOauthConfig *oauth2.Config
	jwtService        *auth.JWTService
	sessionService    *auth.SessionService
	emailService      *email.EmailService
}

func NewAuthService(db *sqlx.DB, googleOauthConfig *oauth2.Config, jwtService *auth.JWTService, sessionService *auth.SessionService, emailService *email.EmailService) *AuthService {
	return &AuthService{
		db:                db,
		googleOauthConfig: googleOauthConfig,
		jwtService:        jwtService,
		sessionService:    sessionService,
		emailService:      emailService,
	}
}
