# SIAKAD Mini — RESTful API (UTS Praktikum Backend Lanjut)

Backend SIAKAD sederhana: mengelola **mahasiswa**, **mata kuliah**, dan **KRS (enrollments)**.
Stack: **Go + Fiber v2 + PostgreSQL (pgx v5) + JWT + validator v10**.

- Base URL lokal: `http://localhost:3000`
- Prefix API: `/api/v1`
- Auth: `Authorization: Bearer <access_token>` — wajib di semua endpoint kecuali `POST /auth/login` (+ `POST /auth/refresh`, `POST /auth/logout` memakai refresh token di body).

## Daftar Isi

- [Setup](#setup)
- [Format Response & Status Code](#format-response--status-code)
- [RBAC](#rbac)
- [Business Rule KRS](#business-rule-krs)
- [Dokumentasi Endpoint](#dokumentasi-endpoint)
  - [1. Login](#1-post-apiv1authlogin)
  - [2. Profil saya](#2-get-apiv1authme)
  - [3. List students](#3-get-apiv1students)
  - [4. Create student](#4-post-apiv1students)
  - [5. Detail student](#5-get-apiv1studentsid)
  - [6. Update student](#6-put-apiv1studentsid)
  - [7. Delete student](#7-delete-apiv1studentsid)
  - [8. List courses](#8-get-apiv1courses)
  - [9. Create enrollment](#9-post-apiv1enrollments)
  - [10. Delete enrollment](#10-delete-apiv1enrollmentsid)

## Setup

```bash
go mod tidy
cp .env.example .env   # isi kredensial DB + JWT_SECRET (min. 32 karakter)
psql -U postgres -c "CREATE DATABASE siakad;"
psql -U postgres -d siakad -f migrations/001_main_tables.sql
# + jalankan seeder: min. 1 admin, 20 mahasiswa (password = hash NIM), 10 mata kuliah
go run main.go
```

`.env` yang dipakai:

```env
APP_PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=***
DB_NAME=siakad
DB_SSLMODE=disable
DB_MAX_CONNS=10
JWT_SECRET=<min-32-karakter>
JWT_ISSUER=praktikum-backend
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7
ALLOWED_ORIGINS=http://localhost:5173
```

Health check: `GET /api/v1/health` → `200` bila server + database jalan.

## Format Response & Status Code

Sukses:

```json
{ "success": true, "message": "login success", "data": { "...": "..." }, "meta": { "page": 1, "limit": 10, "total": 20, "total_pages": 2 } }
```

- List memakai `meta`; create memakai header `Location: /api/v1/.../{id}`; delete memakai `204` tanpa body.
- Error selalu: `{ "success": false, "code": "VALIDATION_ERROR", "message": "...", "errors": { "nim": "..." }, "request_id": "..." }`.

| Status | Arti / kapan keluar |
| --- | --- |
| 200 | OK: login, me, list, detail, update |
| 201 | Created: `POST /students`, `POST /enrollments` |
| 204 | No Content: `DELETE /students/{id}`, `DELETE /enrollments/{id}` |
| 400 | Bad Request: body bukan JSON valid, `:id` bukan angka |
| 401 | Unauthorized: tanpa/expired/invalid token, kredensial salah |
| 403 | Forbidden: role tidak punya permission / akses milik orang lain |
| 404 | Not Found: data tidak ada / sudah soft-delete |
| 409 | Conflict: duplikat (email/nim, atau MK sama di TA sama) |
| 413 | Payload Too Large: body > 1 MB |
| 415 | Unsupported Media Type: `Content-Type` bukan `application/json` |
| 422 | Validation: validasi gagal / kuota penuh / SKS melebihi batas |
| 429 | Login > 5x/menit per IP (`Retry-After: 60`) |
| 500 | Internal error (tanpa stack trace di production) |

## RBAC

| Permission | Admin | Mahasiswa |
| --- | --- | --- |
| `student:list`, `student:create`, `student:read:any`, `student:update:any`, `student:delete:any` | ✅ | ❌ |
| `student:read:own` | ❌ | ✅ (data sendiri) |
| `course:list` | ✅ | ✅ |
| `enrollment:create`, `enrollment:delete:own` | ❌ | ✅ (KRS sendiri) |

Kepemilikan dicek via `student.user_id == current.user_id` (`CanAccessStudent` / `FindByUserID`).

## Business Rule KRS

- Batas SKS dari IPK terakhir: `>= 3.00 → 24`, `2.50–2.99 → 21`, `< 2.50 → 18`.
- Unik `(student_id, course_id, tahun_akademik)` — MK sama tidak bisa diambil 2x di TA sama (409).
- Kuota penuh tidak bisa diambil (422); pengecekan memakai `SELECT ... FOR UPDATE` dalam transaction.
- SKS melebihi batas → 422 dengan pesan sisa SKS, mis. `"SKS melebihi batas. Sisa SKS Anda: 3"`.

---

## Dokumentasi Endpoint

### 1. POST /api/v1/auth/login

Publik. Rate limit 5x/menit/IP.

```http
POST /api/v1/auth/login
Content-Type: application/json

{ "email": "admin@siakad.id", "password": "Admin12345" }
```

- 200:

```json
{
  "success": true, "message": "login success",
  "data": { "access_token": "<jwt>", "refresh_token": "<opaque>", "token_type": "Bearer", "expires_in": 900 }
}
```

- Error: `422` validasi, `401` kredensial salah, `429` kebanyakan percobaan.

Endpoint token lain (di luar 10 wajib tapi tersedia): `POST /auth/refresh` dan `POST /auth/logout` memakai body `{ "refresh_token": "..." }`; refresh memakai rotasi (revoke lama → terbit baru).

### 2. GET /api/v1/auth/me

Semua role yang login → 200.

```http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

- 200 (mahasiswa): user + student + permissions:

```json
{
  "success": true, "message": "successfully got profile information",
  "data": {
    "user": { "id": 7, "email": "mhs01@mail.id", "role": "mahasiswa" },
    "student": { "id": 5, "nim": "187221000001", "nama": "Rina", "prodi": "SI", "angkatan": "2022" },
    "permissions": ["course:list", "enrollment:create", "enrollment:delete:own", "student:read:own"]
  }
}
```

- 200 (admin): sama tanpa blok `student`.
- Error: `401` tanpa/expired/invalid token.

### 3. GET /api/v1/students

Admin → 200 + `meta`. Query: `page` (default 1), `limit` (default 10, maks 50), `search` (nama/nim), `prodi`, `angkatan`, `sort`, `order`.

```http
GET /api/v1/students?page=1&limit=10&search=rina&prodi=SI
Authorization: Bearer <admin_token>
```

```json
{
  "success": true, "message": "Successful fetching Students data",
  "data": [{ "id": 5, "nama": "Rina", "nim": "187221000001", "prodi": "SI", "angkatan": "2022", "ipk_terakhir": 3.45 }],
  "meta": { "page": 1, "limit": 10, "total": 20, "total_pages": 2 }
}
```

- Error: `401`, `403` bukan admin. Soft-deleted (`deleted_at` terisi) tidak muncul.

### 4. POST /api/v1/students

Admin → 201 + header `Location`.

```http
POST /api/v1/students
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "email": "mhs21@mail.id", "password": "Mhs12345",
  "nama": "Budi", "nim": "187221000021", "prodi": "SI", "angkatan": "2023"
}
```

Validasi: `email` required+email, `password` min 8 tanpa spasi, `nama/nim/prodi/angkatan` required (maks 80/50/30/4). Proses: hash password → buat `users` (role mahasiswa) → buat `students`.

- Error: `403`, `422`/`409` duplikat atau validasi gagal.

### 5. GET /api/v1/students/{id}

Admin (data mana pun) + mahasiswa (data sendiri) → 200.

```http
GET /api/v1/students/5
Authorization: Bearer <token>
```

```json
{ "success": true, "message": "Student found!", "data": { "id": 5, "nim": "187221000001", "nama": "Rina", "user_id": 7 } }
```

- Error: `400` id tidak valid, `401`, `403` akses data orang lain, `404` tidak ada/sudah dihapus.

### 6. PUT /api/v1/students/{id}

Admin → 200. Body: `nama`, `nim`, `prodi`, `angkatan`, `ipk_terakhir` (nim tidak boleh diubah menurut penugasan).

```http
PUT /api/v1/students/5
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "nama": "Rina Putri", "nim": "187221000001", "prodi": "SI", "angkatan": "2022", "ipk_terakhir": 3.6 }
```

- Error: `400`/`404`/`422`, `403` bukan admin.

### 7. DELETE /api/v1/students/{id}

Admin → 204 (soft delete via `deleted_at`).

```http
DELETE /api/v1/students/5
Authorization: Bearer <admin_token>
```

- Efek: hilang dari list/detail dan tidak dapat login.
- Error: `401`/`403`/`404`.

### 8. GET /api/v1/courses

Semua role yang login → 200. Tiap MK menyertakan `terisi` dan `sisa_kuota` (dihitung dari `enrollments`).

```http
GET /api/v1/courses
Authorization: Bearer <token>
```

```json
{
  "success": true, "message": "Successful fetching course data",
  "data": [{ "id": 3, "kode_mk": "IF301", "nama_mk": "Basis Data", "sks": 3, "semester": 3, "kuota": 30, "terisi": 28, "sisa_kuota": 2 }]
}
```

- Error: `401`.

### 9. POST /api/v1/enrollments

Mahasiswa → 201.

```http
POST /api/v1/enrollments
Authorization: Bearer <mahasiswa_token>
Content-Type: application/json

{ "course_id": 3, "tahun_akademik": "2026/2027-Ganjil" }
```

```json
{
  "success": true, "message": "Enrollment successfully created",
  "data": { "id": 12, "student_id": 5, "course_id": 3, "tahun_akademik": "2026/2027-Ganjil" }
}
```

- Error: `403` bukan mahasiswa, `409` sudah diambil di TA sama, `422` kuota penuh / SKS melebihi batas.

### 10. DELETE /api/v1/enrollments/{id}

Mahasiswa (milik sendiri) → 204; kuota otomatis bertambah karena dihitung ulang.

```http
DELETE /api/v1/enrollments/12
Authorization: Bearer <mahasiswa_token>
```

- Error: `403` milik mahasiswa lain, `404` tidak ditemukan.
