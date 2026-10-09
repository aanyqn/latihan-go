package repository

import (
	"context"
	"fmt"
	"siakad/app/model"
	"siakad/helper"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepository interface {
	FindAll(ctx context.Context, q helper.ListQuery) ([]model.Course, int, error)
}

var sortColumnCourse = map[string]string{
	"id":      "id",
	"nama_mk": "nama_mk",
	"kode_mk": "kode_mk",
}

type coursePostgreRepository struct {
	pool *pgxpool.Pool
}

func NewCourseRepository(pool *pgxpool.Pool) CourseRepository {
	return &coursePostgreRepository{pool: pool}
}

func courseFilters(q helper.ListQuery) (string, []any) {
	where := " WHERE 1 = 1"
	args := []any{}
	if q.Search != "" {
		where += fmt.Sprintf(" AND (nama_mk ILIKE $%d OR kode_mk ILIKE $%d)",
			len(args)+1, len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}
	return where, args
}

func (r *coursePostgreRepository) FindAll(ctx context.Context, q helper.ListQuery) ([]model.Course, int, error) {
	where, args := courseFilters(q)

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM courses "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"Count courses: %w", err,
		)
	}

	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}
	sqlText := fmt.Sprintf(
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota
		FROM courses %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortColumnCourse[q.Sort], arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())
	rows, err := r.pool.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Get courses: %w", err)
	}
	defer rows.Close()

	hasil := []model.Course{}
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota); err != nil {
			return nil, 0, fmt.Errorf("Course rows: %w", err)
		}
		hasil = append(hasil, c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Query output: %w", err)
	}
	return hasil, total, nil
}
