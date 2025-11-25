package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	APIKey    string    `json:"apiKey"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	APIKey   string `json:"apiKey"`
	Username string `json:"username"`
}

// UpdateUserRequest represents fields that can be updated for a user.
// Pointer fields are used so the service can detect which fields were provided.
type UpdateUserRequest struct {
	ID       int     `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	APIKey   *string `json:"apiKey,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}
