package model

import "time"

type Student struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Nama        string    `json:"nama"`
	Prodi       string    `json:"prodi"`
	Angkatan    string    `json:"angkatan"`
	IPKTerakhir *float64   `json:"ipk_terakhir"`
	NIM         string    `json:"nim"`
	DeletedAt   time.Time `json:"deleted_at"`
	User        User      `json:"user"`
}

type CreateStudentRequest struct {
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,nospace"`
	Nama     string `json:"nama" validate:"required,max=80"`
	NIM      string `json:"nim" validate:"required,max=50"`
	Prodi    string `json:"prodi" validate:"required,max=30"`
	Angkatan string `json:"angkatan" validate:"required,max=4"`
}

type ReplaceStudentRequest struct {
	Nama        *string  `json:"nama" validate:"required,max=80"`
	NIM         *string  `json:"nim" validate:"required,max=50"`
	Prodi       *string  `json:"prodi" validate:"required,max=30"`
	Angkatan    *string  `json:"angkatan" validate:"required,max=4"`
	IPKTerakhir *float64 `json:"ipk_terakhir" validate:"required"`
}
