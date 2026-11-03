package model

import "time"

// Portfolio adalah bentuk portofolio yang dikenal seluruh aplikasi.
//
// Sama seperti model.Review, struct ini tidak memakai satu pun tipe dari
// driver MongoDB. ID disimpan sebagai string; penerjemahan ke bson.ObjectID
// dilakukan repository.
type Portfolio struct {
	ID           string        `json:"id"`
	StudentID    int           `json:"student_id"`
	Headline     string        `json:"headline"`
	Bio          string        `json:"bio"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// Project disimpan DI DALAM dokumen Portfolio (embedding). Proyek selalu
// dibaca bersama portofolionya dan jumlahnya dibatasi 20.
type Project struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Technologies []string `json:"technologies"`
	Link         string   `json:"link"`
	Year         int      `json:"year"`
}

// Certificate juga disarangkan di dalam dokumen Portfolio.
type Certificate struct {
	Name     string    `json:"name"`
	Issuer   string    `json:"issuer"`
	IssuedAt time.Time `json:"issued_at"`
}

// ProjectRequest adalah body untuk menambah proyek (POST) maupun proyek di
// dalam body PUT. Memakai struct bertipe tetap membuat operator MongoDB di
// dalam JSON ditolak sebelum menyentuh database.
type ProjectRequest struct {
	Title        string   `json:"title"        validate:"required,min=2,max=150"`
	Technologies []string `json:"technologies" validate:"max=15,dive,min=1,max=50"`
	Link         string   `json:"link"         validate:"omitempty,url,max=500"`
	Year         int      `json:"year"         validate:"omitempty,min=1900,max=2100"`
}

// CertificateRequest memakai tanggal berbentuk string agar client dapat
// mengirim "2006-01-02" maupun RFC3339. Service yang mengubahnya menjadi
// time.Time sebelum disimpan sebagai BSON Date.
type CertificateRequest struct {
	Name     string `json:"name"      validate:"required,min=2,max=150"`
	Issuer   string `json:"issuer"    validate:"required,min=2,max=150"`
	IssuedAt string `json:"issued_at" validate:"required"`
}

// ReplacePortfolioRequest adalah body PUT /students/:id/portfolio. PUT
// mengganti seluruh isi portofolio, termasuk projects dan certificates.
type ReplacePortfolioRequest struct {
	Headline     string               `json:"headline"     validate:"required,min=3,max=120"`
	Bio          string               `json:"bio"          validate:"max=2000"`
	Projects     []ProjectRequest     `json:"projects"     validate:"max=20,dive"`
	Certificates []CertificateRequest `json:"certificates" validate:"max=50,dive"`
}

// TechnologyStat adalah satu baris hasil aggregation teknologi terpopuler.
type TechnologyStat struct {
	Technology string `json:"technology"`
	Total      int    `json:"total"`
}
