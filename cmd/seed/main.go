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
	script = strings.ReplaceAll(script, "__ADMIN_PASSWORD_HASH__", passwordHash)

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
	fmt.Println("  owner@berezhok.local         / test12345  (Хлеб и Кофе, owner)")
	fmt.Println("  employee@berezhok.local       / test12345  (Хлеб и Кофе, employee)")
	fmt.Println("  bakery.manager@berezhok.local / test12345  (Хлеб и Кофе, manager — must_change_password)")
	fmt.Println("  dinner.owner@berezhok.local   / test12345  (Вечерняя Трапеза, owner)")
	fmt.Println("  dinner.staff@berezhok.local   / test12345  (Вечерняя Трапеза, employee — must_change_password)")
	fmt.Println("  coffee.owner@berezhok.local   / test12345  (Городской Вкус, owner)")
	fmt.Println("  coffee.manager@berezhok.local / test12345  (Городской Вкус, manager)")
	fmt.Println("  sushi.owner@berezhok.local    / test12345  (Роллы и Суши, owner)")
	fmt.Println("  pizza.owner@berezhok.local    / test12345  (Пицца Плюс, owner — partner suspended)")
	fmt.Println("  artizan.owner@berezhok.local  / test12345  (Пекарня Артизан)")
	fmt.Println("  sushidoma.owner@berezhok.local / test12345 (Суши Дома)")
	fmt.Println("  coffeeugol.owner@berezhok.local / test12345 (Кофе Угол)")
	fmt.Println("  bistro24.owner@berezhok.local / test12345  (Бистро 24)")
	fmt.Println("  zelenaya.owner@berezhok.local / test12345  (Зелёная лавка)")
	fmt.Println("  myasnoy.owner@berezhok.local  / test12345  (Мясной двор)")
	fmt.Println("  yaponka.owner@berezhok.local  / test12345  (Японский квартал)")
	fmt.Println("  pirozhki.owner@berezhok.local / test12345  (Пирожковая №1)")
	fmt.Println("  smoothie.owner@berezhok.local / test12345  (Смузи Бар)")
	fmt.Println("  vostok.owner@berezhok.local   / test12345  (Восточный базар)")
	fmt.Println("Admin accounts:")
	fmt.Println("  admin@berezhok.local          / test12345  (super_admin)")
	fmt.Println("  moderator@berezhok.local      / test12345  (admin)")
	fmt.Println("  support@berezhok.local        / test12345  (support)")
	fmt.Println("Customer phones (+79990000001 .. +79990000012)")
}
