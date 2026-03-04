package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

func New(ctx context.Context, dbConfig config.Db) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%v/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
