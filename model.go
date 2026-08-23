package main

import "time"

// Student
type Student struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
	NIM      string  `json:"nim"`
}

type CreateStudentRequest struct {
	Username string `json:"username"`
	NIM      string `json:"nim"`
}

type ReplaceStudentRequest struct {
	Username string   `json:"username"`
	NIM      string   `json:"nim"`
	Grade    *float64 `json:"grade"`
	IsActive bool     `json:"is_active"`
}

type PatchStudentRequest struct {
	Username *string  `json:"username,omitempty"`
	NIM      *string   `json:"nim"`
	Grade    *float64 `json:"grade,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// User
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
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

type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	isActive *bool
	GradeMin *float64
	GradeMax *float64
}
