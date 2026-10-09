package repository

import (
	"context"
	"errors"
	"fmt"
	"siakad/app/model"
	"siakad/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository interface {
	FindAll(ctx context.Context, q helper.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	FindByUserID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

var sortColumnStudent = map[string]string{
	"id": "id",
	"nama": "nama",
	"ipk_terakhir": "ipk_terakhir",
}

type studentPostgreRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgreRepository{pool: pool}
}

func studentFilters(q helper.ListQuery) (string, []any) {
	where := " WHERE 1 = 1 AND s.deleted_at IS NULL"
	args := []any{}
	if q.Search != "" {
		where += fmt.Sprintf(" AND (s.nama ILIKE $%d OR s.nim ILIKE $%d)",
			len(args)+1, len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}
	if q.Prodi != "" {
		where += fmt.Sprintf(" AND (s.prodi ILIKE $%d)",
			len(args)+1)
		args = append(args, "%"+q.Prodi+"%")
	}
	if q.Angkatan != "" {
		where += fmt.Sprintf(" AND (s.angkatan ILIKE $%d)",
			len(args)+1)
		args = append(args, "%"+q.Angkatan+"%")
	}
	return where, args
}

func (r *studentPostgreRepository) FindAll(
	ctx context.Context, q helper.ListQuery,
) ([]model.Student, int, error) {
	where, args := studentFilters(q)

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM students s LEFT JOIN users u ON u.id = s.user_id "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Count student: %w", err)
	}
	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}
	sqlText := fmt.Sprintf(
		`SELECT s.id, s.nama, s.nim, s.prodi, s.angkatan, s.ipk_terakhir, u.email, u.role
		FROM students s
		LEFT JOIN users u
		ON s.user_id = u.id %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortColumnStudent[q.Sort], arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())
	rows, err := r.pool.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Get student: %w", err)
	}
	defer rows.Close()
	hasil := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.Nama, &s.NIM,
			&s.Prodi, &s.Angkatan, &s.IPKTerakhir, &s.User.Email, &s.User.Role); err != nil {
			return nil, 0, fmt.Errorf("Student rows: %w", err)
		}
		hasil = append(hasil, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Query output: %w", err)
	}
	return hasil, total, nil
}

func (r *studentPostgreRepository) FindByID(
	ctx context.Context, id int,
) (model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx,
		`SELECT s.id, s.nama, s.nim, s.prodi, s.angkatan, s.ipk_terakhir, s.user_id, u.email, u.role
		FROM students s
		LEFT JOIN users u
		ON s.user_id = u.id WHERE s.id = $1 AND s.deleted_at IS NULL`, id,
	).Scan(&s.ID, &s.Nama, &s.NIM, &s.Prodi, &s.Angkatan, &s.IPKTerakhir, &s.UserID, &s.User.Email, &s.User.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, helper.ErrNotFound
		}
		return model.Student{}, fmt.Errorf("Get Student: %w", err)
	}
	return s, nil
}

func (r *studentPostgreRepository) FindByUserID(
	ctx context.Context, id int,
) (model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx,
		`SELECT s.id, s.nama, s.nim, s.prodi, s.angkatan, s.ipk_terakhir, s.user_id, u.id, u.email, u.role
		FROM students s
		LEFT JOIN users u
		ON s.user_id = u.id WHERE s.user_id = $1 AND s.deleted_at IS NULL`, id,
	).Scan(&s.ID, &s.Nama, &s.NIM, &s.Prodi, &s.Angkatan, &s.IPKTerakhir, &s.UserID, &s.User.ID, &s.User.Email, &s.User.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, helper.ErrNotFound
		}
		return model.Student{}, fmt.Errorf("Get Student: %w", err)
	}
	return s, nil
}

func (r *studentPostgreRepository) Create(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO students (nama, nim, prodi, angkatan, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, nama, user_id`,
		s.Nama, s.NIM, s.Prodi, s.Angkatan, s.User.ID,
	).Scan(&s.ID, &s.Nama, &s.UserID)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, helper.ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Creating Student: %w", err)
	}
	return s, nil
}

func (r *studentPostgreRepository) Update(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`UPDATE students SET nama = $1, nim = $2, prodi = $3, angkatan=$4, ipk_terakhir=$5
		WHERE id = $6
		RETURNING id, nama, nim, prodi, angkatan, ipk_terakhir, user_id`,
		s.Nama, s.NIM, s.Prodi, s.Angkatan, s.IPKTerakhir, s.ID,
	).Scan(&s.ID, &s.Nama, &s.NIM, &s.Prodi, &s.Angkatan, &s.IPKTerakhir, &s.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, helper.ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, helper.ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Update student: %w", err)
	}
	return s, nil
}

func (r *studentPostgreRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `UPDATE students SET deleted_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete student: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return helper.ErrNotFound
	}
	return nil
}
