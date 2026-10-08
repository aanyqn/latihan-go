CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL
);

CREATE INDEX IF NOT EXISTS users_email_lower_idx ON users (LOWER(email));

CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    nim VARCHAR(50) NOT NULL,
    nama VARCHAR(100) NOT NULL,
    prodi VARCHAR(30) NOT NULL,
    angkatan VARCHAR(4) NOT NULL,
    ipk_terakhir FLOAT,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS students_nim ON students (nim);

ALTER TABLE students
ADD CONSTRAINT students_fkey_users FOREIGN KEY(user_id) REFERENCES users(id);

CREATE TABLE IF NOT EXISTS course (
    id SERIAL PRIMARY KEY,
    kode_mk VARCHAR(10) NOT NULL,
    nama_mk VARCHAR(100) NOT NULL,
    sks INT NOT NULL,
    semester INT NOT NULL,
    kuota INT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS course_kode_mk ON course (kode_mk);

CREATE TABLE IF NOT EXISTS enrollments (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL,
    course_id INT NOT NULL,
    tahun_akademik VARCHAR(6),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE enrollments
ADD CONSTRAINT enrollments_fkey_students FOREIGN KEY(student_id) REFERENCES students(id);

ALTER TABLE enrollments
ADD CONSTRAINT enrollments_fkey_courses FOREIGN KEY(course_id) REFERENCES courses(id);

CREATE UNIQUE INDEX IF NOT EXISTS enrollments_students_courses_year ON enrollments (student_id, course_id, tahun_akademik);

CREATE TABLE IF NOT EXISTS refresh_tokens (
 id BIGSERIAL PRIMARY KEY,
 user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_hash TEXT NOT NULL UNIQUE,
 expires_at TIMESTAMPTZ NOT NULL,
 revoked_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx
 ON refresh_tokens (user_id)

CREATE TABLE IF NOT EXISTS permissions (
	name VARCHAR(50) PRIMARY KEY,
	description VARCHAR(150) NOT NULL
);

CREATE INDEX IF NOT EXISTS permissions_name_lower ON permissions (LOWER(name));

CREATE TABLE IF NOT EXISTS role_permissions (
	role_name VARCHAR(20),
	permissions_name VARCHAR(50) NOT NULL REFERENCES permissions (name) ON DELETE CASCADE,
	PRIMARY KEY (role_name, permissions_name)
);

CREATE INDEX IF NOT EXISTS role_permissions_role_lower ON role_permissions (LOWER(role_name));

INSERT INTO permissions (name, description)
VALUES
    ('student:list', 'Melihat daftar seluruh mahasiswa dengan pagination dan filter'),
    ('student:create', 'Menambah mahasiswa baru sekaligus akun usernya'),
    ('student:read:any', 'Melihat detail data mahasiswa mana pun'),
    ('student:read:own', 'Melihat detail data mahasiswa milik sendiri'),
    ('student:update:any', 'Memperbarui data mahasiswa mana pun'),
    ('student:delete:any', 'Menghapus (soft delete) data mahasiswa mana pun'),
    ('course:list', 'Melihat daftar mata kuliah beserta sisa kuota'),
    ('enrollment:create', 'Mengambil mata kuliah (menambah KRS)'),
    ('enrollment:delete:own', 'Membatalkan mata kuliah dari KRS milik sendiri')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permissions_name)
VALUES
    ('admin', 'student:list'),
    ('admin', 'student:create'),
    ('admin', 'student:read:any'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete:any'),
    ('admin', 'course:list'),

    ('mahasiswa', 'student:read:own'),
    ('mahasiswa', 'course:list'),
    ('mahasiswa', 'enrollment:create'),
    ('mahasiswa', 'enrollment:delete:own')
ON CONFLICT DO NOTHING;