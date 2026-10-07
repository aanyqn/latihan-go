package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	LoadPermissions(ctx context.Context) (map[string][]string, error)
}

type rolePostgresRepository struct{ pool *pgxpool.Pool }

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &rolePostgresRepository{pool: pool}
}

func (r *rolePostgresRepository) LoadPermissions(ctx context.Context) (map[string][]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT role_name, COALESCE(permissions_name, '') FROM role_permissions ORDER BY role_name, permissions_name`)
	if err != nil {
		return nil, fmt.Errorf("Get permission: %w", err)
	}
	defer rows.Close()

	result := map[string][]string{}
	for rows.Next() {
		var role, permission string
		if err := rows.Scan(&role, &permission); err != nil {
			return nil, fmt.Errorf("Reading row permission: %w", err)
		}
		if _, ok := result[role]; !ok {
			result[role] = []string{}
		}
		if permission != "" {
			result[role] = append(result[role], permission)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Reading query result permission: %w", err)
	}
	return result, nil
}
