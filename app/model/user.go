package model

import "time"

// User
type User struct {
	ID        int       `json:"id"`
	Role      string    `json:"role"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReplaceUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

type PatchUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}
