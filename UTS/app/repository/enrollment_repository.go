package repository

import (
	"context"
	"fmt"
	"siakad/app/model"
	"siakad/helper"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, e model.Enrollment, ipk float64) (model.Enrollment, error)
	Delete(ctx context.Context, id int) error
}

type enrollmentPostgreRepository struct {
	pool *pgxpool.Pool
}

func NewEnrollmentRepository(pool *pgxpool.Pool) EnrollmentRepository {
	return &enrollmentPostgreRepository{pool: pool}
}

func (r *enrollmentPostgreRepository) Create(
	ctx context.Context, e model.Enrollment, ipk float64,
) (model.Enrollment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.Enrollment{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var sks, kuota int
	err = tx.QueryRow(ctx, "SELECT sks, kuota FROM courses WHERE id = $1 FOR UPDATE", e.CourseID).Scan(&sks, &kuota)
	if err != nil {
		return model.Enrollment{}, helper.BadRequest("Mata kuliah tidak ditemukan")
	}

	var enrolledCount int
	err = tx.QueryRow(ctx, "SELECT count(*) FROM enrollments WHERE course_id = $1 AND tahun_akademik ILIKE $2", e.CourseID, e.TahunAkademik).Scan(&enrolledCount)
	if err != nil {
		return model.Enrollment{}, err
	}
	if enrolledCount >= kuota {
		return model.Enrollment{}, helper.ErrQuotaFull
	}

	maxSKS := 18
	if ipk >= 3.00 {
		maxSKS = 24
	} else if ipk >= 2.50 {
		maxSKS = 21
	}

	var totalSKS int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(c.sks), 0)
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE e.student_id = $1 AND e.tahun_akademik = $2
	`, e.StudentID, e.TahunAkademik).Scan(&totalSKS)
	if err != nil {
		return model.Enrollment{}, err
	}

	if totalSKS+sks > maxSKS {
		sisaSKS := maxSKS - totalSKS
		if sisaSKS < 0 {
			sisaSKS = 0
		}
		return model.Enrollment{}, helper.UnprocessableEntity(fmt.Sprintf("SKS melebihi batas. Sisa SKS Anda: %d", sisaSKS))
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		VALUES ($1, $2, $3)
		RETURNING id, student_id, course_id, tahun_akademik`,
		e.StudentID, e.CourseID, e.TahunAkademik,
	).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Enrollment{}, helper.ErrDuplicate
		}
		return model.Enrollment{}, fmt.Errorf("Creating Enrollment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Enrollment{}, fmt.Errorf("commit transaction: %w", err)
	}

	return e, nil
}

func (e *enrollmentPostgreRepository) Delete(ctx context.Context, id int) error {
	tag, err := e.pool.Exec(ctx, `DELETE FROM enrollments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return helper.ErrNotFound
	}
	return nil
}
