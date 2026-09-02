package repository

import (
	"context"
	"errors"
	"fmt"
	"latihan-fiber/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, u model.Student) (model.Student, error)
	Update(ctx context.Context, u model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

var sortColumn = map[string]string{
	"id":         "id",
	"username":   "username",
	"nim":        "nim",
	"created_at": "created_at",
}

type studentPostgreRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgreRepository{pool: pool}
}

func buildFilters(q model.ListQuery) (string, []any) {
	where := " WHERE 1 = 1"
	args := []any{}
	if q.Search != "" {
		where += fmt.Sprintf(" AND (username ILIKE $%d OR nim ILIKE $%d)",
			len(args)+1, len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}
	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}
	return where, args
}

func (r *studentPostgreRepository) FindAll(
	ctx context.Context, q model.ListQuery,
) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Count student: %w", err)
	}
	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}
	sqlText := fmt.Sprintf(
		`SELECT id, username, nim, is_active, created_at
		FROM students %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortColumn[q.Sort], arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())
	rows, err := r.pool.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Get student: %w", err)
	}
	defer rows.Close()
	hasil := []model.Student{}
	for rows.Next() {
		var u model.Student
		if err := rows.Scan(&u.ID, &u.Username, &u.NIM,
			&u.IsActive, &u.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("Student rows: %w", err)
		}
		hasil = append(hasil, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Query output: %w", err)
	}
	return hasil, total, nil
}

func (r *studentPostgreRepository) FindByID(
	ctx context.Context, id int,
) (model.Student, error) {
	var u model.Student
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, nim, is_active, created_at
		FROM students WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.NIM, &u.IsActive, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("Get Student: %w", err)
	}
	return u, nil
}

func (r *studentPostgreRepository) Create(
	ctx context.Context, u model.Student,
) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO students (username, nim, grade, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		u.Username, u.NIM, u.Grade, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Saving Student: %w", err)
	}
	return u, nil
}

func (r *studentPostgreRepository) Update(
	ctx context.Context, u model.Student,
) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`UPDATE students SET username = $1, nim = $2, is_active = $3
		WHERE id = $4
		RETURNING id, username, nim, is_active, created_at`,
		u.Username, u.NIM, u.IsActive, u.ID,
	).Scan(&u.ID, &u.Username, &u.NIM, &u.IsActive, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Update student: %w", err)
	}
	return u, nil
}

func (r *studentPostgreRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete student: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
