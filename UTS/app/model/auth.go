package model

import "time"

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,max=72,strongpassword"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // detik
}

type RefreshToken struct {
	ID        int64
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type AuthUser struct {
	UserID   int    `json:"user_id"`
	Email string `json:"email"`
	Role     string `json:"role"`
}
