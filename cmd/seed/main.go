package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

//go:embed seed.sql
var seedFS embed.FS

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	if cfg.Env != envLocal && cfg.Env != envDev && os.Getenv("SEED_FORCE") != "1" {
		slog.Error("seed is allowed only in local/dev; set SEED_FORCE=1 to override", slog.String("env", cfg.Env))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := postgresql.New(ctx, cfg.Db)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	seedSQL, err := seedFS.ReadFile("seed.sql")
	if err != nil {
		slog.Error("failed to read seed.sql", slog.Any("error", err))
		os.Exit(1)
	}

	passwordHash, err := auth.Hash("test12345")
	if err != nil {
		slog.Error("failed to hash password", slog.Any("error", err))
		os.Exit(1)
	}

	script := strings.ReplaceAll(string(seedSQL), "__PARTNER_PASSWORD_HASH__", passwordHash)

	tx, err := db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", slog.Any("error", err))
		os.Exit(1)
	}

	if _, err = tx.Exec(ctx, script); err != nil {
		_ = tx.Rollback(ctx)
		slog.Error("failed to execute seed script", slog.Any("error", err))
		os.Exit(1)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("failed to commit seed transaction", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println("DB seed completed.")
	fmt.Println("Partner accounts:")
	fmt.Println("  owner@berezhok.local / test12345")
	fmt.Println("  coffee.owner@berezhok.local / test12345")
	fmt.Println("  employee@berezhok.local / test12345")
	fmt.Println("Customer phones:")
	fmt.Println("  +79990000001")
	fmt.Println("  +79990000002")
	fmt.Println("  +79990000003")
}
