CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    nim VARCHAR(50) NOT NULL,
    grade FLOAT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS students_username_lower_key ON students (LOWER(username));

CREATE UNIQUE INDEX IF NOT EXISTS students_nim ON students (nim);