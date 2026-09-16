package repository

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"latihan-fiber/app/model"
// 	"latihan-fiber/helper"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type AchievementRepository interface {
// 	FindByID(ctx context.Context, id int) (model.Achievement, error)
// }

// type achievementPostgreRepository struct {
// 	pool *pgxpool.Pool
// }

// func NewAchievementRepository(pool *pgxpool.Pool) AchievementRepository {
// 	return &achievementPostgreRepository{pool: pool}
// }

// func buildFilterAchievment(q helper.ListQuery) (string, []any) {
// 	where := " WHERE 1 = 1"
// 	args := []any{}
// 	if q.Filter != nil {
// 		where += fmt.Sprintf(" AND id = $%d", len(args)+1)
// 		args = append(args, *q.IsActive)
// 	}
// 	return where, args
// }

// func (r *achievementPostgreRepository) FindByID(
// 	ctx context.Context, id int, q helper.ListQuery,
// ) ([]model.Achievement, error) {
// 	where, args := buildFilter(q)
// 	sqlText := fmt.Sprintf(
// 		`SELECT id, name, champion, student_id
// 		FROM achievement`,
// 		where,
// 	)
// 	args = append(args, q.Limit, q.Offset())
// 	rows, err := r.pool.Query(ctx, sqlText, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("Get Achievements: %w", err)
// 	}
// 	defer rows.Close()
// 	hasil := []model.Student{}
// 	for rows.Next() {
// 		var s model.Student
// 		if err := rows.Scan(&s.ID, &s.Username, &s.NIM,
// 			&s.IsActive, &s.CreatedAt); err != nil {
// 			return nil, 0, fmt.Errorf("Student rows: %w", err)
// 		}
// 		hasil = append(hasil, s)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, 0, fmt.Errorf("Query output: %w", err)
// 	}
// 	return hasil, total, nil
// }