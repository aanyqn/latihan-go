INSERT INTO
    permissions (name, description)
VALUES
    (
        'student:list',
        'Melihat daftar seluruh data mahasiswa'
    ),
    (
        'student:read:any',
        'Melihat detail data mahasiswa mana pun'
    ),
    (
        'student:create',
        'Menambahkan data mahasiswa baru'
    ),
    (
        'student:update:any',
        'Mengubah/memperbarui data mahasiswa mana pun'
    ),
    ('student:delete', 'Menghapus data mahasiswa') ON CONFLICT (name) DO NOTHING;

INSERT INTO
    role_permissions (role_name, permissions_name)
VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create') ON CONFLICT DO NOTHING;

ALTER TABLE students
ADD COLUMN owner_id int REFERENCES users (id);

UPDATE students
SET
    owner_id = 7
WHERE
    owner_id IS NULL;

ALTER TABLE students
ALTER COLUMN owner_id
SET
    NOT NULL;