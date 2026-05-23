package auth

import "time"

type Parent struct {
	ID             string    `db:"id"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	MarketingOptIn bool      `db:"marketing_opt_in"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type SignupRequest struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	DisplayName    string `json:"display_name" validate:"required"`
	MarketingOptIn bool   `json:"marketing_opt_in"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Claims struct {
	ActorID   string `json:"actor_id"`
	ActorType string `json:"actor_type"` // "parent" or "child"
	Scope     string `json:"scope"`      // e.g. "parent", "kid", or "parent-elevated"
}
