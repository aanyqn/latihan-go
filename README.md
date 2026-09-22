# Backend Praktikum Latihan Fiber

Repositori ini berisi backend REST API menggunakan **Golang (Fiber)** dan **PostgreSQL**.

## Daftar Isi

- [Kontrak API](#kontrak-api-api-contract)
  - [1. Students](#1-students)
  - [2. Users](#2-users)
  - [3. Auth](#3-auth-authentication)
- [Instalasi](#instalasi)
  - [1. Variabel Environment (.env)](#1-variabel-environment-env)
  - [2. Database](#2-database)

---

## Kontrak API (API Contract)

> **Catatan Keamanan:** Mulai dari penyelesaian Tugas 5, seluruh endpoint pada bagian **Students** dan **Users** membutuhkan header **`Authorization: Bearer <access_token>`** yang valid. Jika token tidak disertakan, kedaluwarsa, atau tidak valid, API akan mengembalikan status `401 Unauthorized`.

### 1. Students

| Metode | Endpoint | Parameter / Query | Contoh Body Permintaan | Status yang Mungkin Dikembalikan | Contoh Respons (Sukses) |
|---|---|---|---|---|---|
| **GET** | `/api/v1/students/` | `page`, `limit`, `search`, `sort`, `order`, `is_active`, `grade_min`, `grade_max` | Tidak ada | `200 OK`, `401 Unauthorized` | `{"success":true,"message":"Berhasil","data":[{"id":1,"username":"Andi","nim":"123","grade":85.5,"is_active":true}],"meta":{...}}` |
| **POST** | `/api/v1/students/` | Tidak ada | `{"username":"Gina","nim":"98765"}` | `201 Created`, `400`, `401`, `409`, `422` | `{"success":true,"message":"Data berhasil dibuat","data":{"id":2,...}}` |
| **GET** | `/api/v1/students/:id` | `id` (integer) | Tidak ada | `200 OK`, `400`, `401`, `404` | `{"success":true,"message":"Data ditemukan","data":{"id":2,...}}` |
| **PUT** | `/api/v1/students/:id` | `id` (integer) | `{"username":"Gina Updated","nim":"98765","grade":90.0,"is_active":true}` | `200 OK`, `400`, `401`, `404`, `422` | `{"success":true,"message":"Data berhasil diubah utuh","data":{...}}` |
| **PATCH** | `/api/v1/students/:id` | `id` (integer) | `{"is_active":true}` | `200 OK`, `400`, `401`, `404`, `422` | `{"success":true,"message":"Data berhasil diperbarui","data":{...}}` |
| **DELETE**| `/api/v1/students/:id` | `id` (integer) | Tidak ada | `204 No Content`, `400`, `401`, `404` | Tidak ada response body |

### 2. Users

| Metode | Endpoint | Parameter / Query | Contoh Body Permintaan | Status yang Mungkin Dikembalikan | Contoh Respons (Sukses) |
|---|---|---|---|---|---|
| **GET** | `/api/v1/users/` | `page`, `limit`, `search`, `sort`, `order`, `is_active` | Tidak ada | `200 OK`, `401 Unauthorized` | `{"success":true,"message":"Berhasil","data":[{"id":1,"username":"budi","email":"budi@ex.com","is_active":true}],"meta":{...}}` |
| **POST** | `/api/v1/users/` | Tidak ada | `{"username":"budi","email":"budi@ex.com","password":"Password1"}` | `201 Created`, `400`, `401`, `409`, `422` | `{"success":true,"message":"Data berhasil dibuat","data":{"id":2,...}}` |
| **GET** | `/api/v1/users/:id` | `id` (integer) | Tidak ada | `200 OK`, `400`, `401`, `404` | `{"success":true,"message":"Data ditemukan","data":{"id":2,...}}` |
| **PUT** | `/api/v1/users/:id` | `id` (integer) | `{"username":"budi ubah","email":"baru@ex.com","is_active":false}` | `200 OK`, `400`, `401`, `404`, `422` | `{"success":true,"message":"Data berhasil diubah utuh","data":{...}}` |
| **PATCH** | `/api/v1/users/:id` | `id` (integer) | `{"is_active":false}` | `200 OK`, `400`, `401`, `404`, `422` | `{"success":true,"message":"Data berhasil diperbarui","data":{...}}` |
| **DELETE**| `/api/v1/users/:id` | `id` (integer) | Tidak ada | `204 No Content`, `400`, `401`, `404` | Tidak ada response body |

### 3. Auth (Authentication)

| Metode | Endpoint | Contoh Body Permintaan | Deskripsi & Status Response |
|---|---|---|---|
| **POST** | `/api/v1/auth/register` | `{"username":"admin","email":"a@a.com","password":"Password1"}` | Mendaftar user baru (role selalu "user"). (`201`, `400`, `409`, `422`) |
| **POST** | `/api/v1/auth/login` | `{"username":"admin","password":"Password1"}` | Login dan mendapatkan pasangan token. Max 5 percobaan/menit. (`200`, `400`, `401`, `403`, `429`) |
| **POST** | `/api/v1/auth/refresh` | `{"refresh_token":"<token>"}` | Memperbarui Access Token. (`200`, `400`, `401`) |
| **POST** | `/api/v1/auth/logout` | `{"refresh_token":"<token>"}` | Mencabut *refresh token* dari database. (`200`) |
| **GET**  | `/api/v1/auth/me` | *Butuh Header Bearer Token* | Mengambil data profil user yang sedang diakses saat ini. (`200`, `401`) |

---

## Instalasi

### 1. Variabel Environment (.env)

Sebelum menjalankan aplikasi, buat file `.env` di *root folder* proyek (sejajar dengan `main.go`). Salin konfigurasi berikut dan sesuaikan dengan kredensial PostgreSQL di komputer Anda:

```env
# Konfigurasi Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_postgres_kamu  # Ubah sesuai password psql lokalmu!
DB_NAME=praktikum_backend
DB_SSLMODE=disable
DB_MAX_CONNS=10

# Konfigurasi App (Opsional)
APP_PORT=3000
```

### 2. Database

**a. Buat database baru**

```bash
psql -U postgres -c "CREATE DATABASE praktikum_backend;"
```

**b. Contoh file migrasi** (`migrations/001_create_users.sql`)

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Keunikan username tanpa membedakan huruf besar dan kecil.
-- Inilah yang menggantikan pemeriksaan manual di pertemuan 2.
CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key
    ON users (LOWER(username));

CREATE INDEX IF NOT EXISTS users_email_lower_idx
    ON users (LOWER(email));
```

**c. Jalankan migrasi dan periksa hasil**

```bash
psql -U postgres -d praktikum_backend -f migrations/001_create_users.sql
psql -U postgres -d praktikum_backend -c "\d users"
```

**d. Konfigurasi database dan connection pool** (`database/database.go`)

```go
package database

import (
	"context"
	"fmt"
	"latihan-fiber/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", ""),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_NAME", "praktikum_backend"),
		config.GetEnv("DB_SSLMODE", "disable"),
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("database config isn't valid: %w", err)
	}

	cfg.MaxConns = int32(config.GetEnvInt("DB_MAX_CONNS", 10))
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return pool, nil
}
```