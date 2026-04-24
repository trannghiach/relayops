package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedChallengeFlag(ctx context.Context, pool *pgxpool.Pool, flag string) error {
	if flag == "" {
		return nil
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO challenge_flags (name, value)
		VALUES ('demo_key', $1)
		ON CONFLICT (name) DO UPDATE
		SET value = EXCLUDED.value
	`, flag)

	return err
}
