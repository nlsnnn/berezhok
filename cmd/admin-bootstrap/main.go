package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

const defaultAdminRole = "super_admin"

func main() {
	cfg := config.MustLoad()

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	name := os.Getenv("ADMIN_NAME")
	role := os.Getenv("ADMIN_ROLE")
	if role == "" {
		role = defaultAdminRole
	}

	if email == "" || password == "" || name == "" {
		slog.Error("ADMIN_EMAIL, ADMIN_PASSWORD and ADMIN_NAME are required")
		os.Exit(1)
	}
	if role != "super_admin" && role != "admin" && role != "support" {
		slog.Error("ADMIN_ROLE must be one of super_admin, admin, support", slog.String("role", role))
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

	q := sqlc.New(db)
	if _, err = q.FindAdminByEmail(ctx, email); err == nil {
		fmt.Printf("Admin already exists: %s\n", email)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to check admin", slog.Any("error", err))
		os.Exit(1)
	}

	passwordHash, err := auth.Hash(password)
	if err != nil {
		slog.Error("failed to hash password", slog.Any("error", err))
		os.Exit(1)
	}

	admin, err := q.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		Role:         role,
		IsActive:     true,
	})
	if err != nil {
		slog.Error("failed to create admin", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Printf("Admin created: %s (%s)\n", admin.Email, admin.Role)
}
