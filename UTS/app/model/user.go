package model

type User struct {
	ID        int       `json:"id"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,nospace"`
}
type ReplaceUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=120"`
	IsActive bool   `json:"is_active"`
}

type PatchUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitnil,email,max=120"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}
