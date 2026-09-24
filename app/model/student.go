package model

import "time"

// Student
type Student struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	NIM       string    `json:"nim"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
}

type CreateStudentRequest struct {
	Username string   `json:"username"`
	NIM      string   `json:"nim"`
	Grade    *float64 `json:"grade"`
}

type ReplaceStudentRequest struct {
	Username *string  `json:"username"`
	NIM      *string  `json:"nim"`
	Grade    *float64 `json:"grade"`
	IsActive *bool    `json:"is_active"`
}

type PatchStudentRequest struct {
	Username *string  `json:"username,omitempty"`
	NIM      *string  `json:"nim"`
	Grade    *float64 `json:"grade,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}
