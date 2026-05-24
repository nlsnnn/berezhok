package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/adapters/yookassa"
	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	payoutRepos "github.com/nlsnnn/berezhok/internal/modules/payout/repository"
	payoutServices "github.com/nlsnnn/berezhok/internal/modules/payout/service"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	runCalc := flag.Bool("run-calculation", false, "run calculation job once and exit")
	runDispatch := flag.Bool("run-dispatch", false, "run dispatch job once and exit")
	runPoll := flag.Bool("run-poll", false, "poll processing payouts once and exit")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("payout-worker starting", slog.String("env", cfg.Env))

	db, err := postgresql.New(ctx, cfg.Db)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	queries := sqlc.New(db)
	payoutRepo := payoutRepos.NewPayoutRepo(queries, db)

	minAmount, err := decimal.NewFromString(cfg.MinPayoutAmount)
	if err != nil {
		log.Error("invalid min payout amount", sl.Err(err))
		os.Exit(1)
	}

	calcSvc := payoutServices.NewCalculationService(payoutRepo, minAmount, log)

	ykPayoutClient := yookassa.NewPayoutHandlerClient(cfg.YookassaPayout)
	payoutAdapter := yookassa.NewPayoutAdapter(ykPayoutClient)
	dispatchSvc := payoutServices.NewDispatchService(payoutRepo, payoutAdapter, log)

	if *runCalc {
		runCalculation(ctx, cfg, calcSvc, log)
		return
	}

	if *runDispatch {
		if err := dispatchSvc.DispatchPending(ctx, 100); err != nil {
			log.Error("dispatch failed", sl.Err(err))
			os.Exit(1)
		}

		return
	}

	if *runPoll {
		if err := dispatchSvc.PollProcessing(ctx, 100); err != nil {
			log.Error("poll failed", sl.Err(err))
			os.Exit(1)
		}

		return
	}

	c := cron.New(cron.WithLogger(cron.PrintfLogger(slog.NewLogLogger(log.Handler(), slog.LevelInfo))))

	_, _ = c.AddFunc(cfg.CalculationCron, func() {
		runCalculation(ctx, cfg, calcSvc, log)
	})

	_, _ = c.AddFunc(cfg.DispatchCron, func() {
		if err := dispatchSvc.DispatchPending(ctx, 20); err != nil {
			log.Error("dispatch job failed", sl.Err(err))
		}
	})

	_, _ = c.AddFunc(cfg.PollCron, func() {
		if err := dispatchSvc.PollProcessing(ctx, 50); err != nil {
			log.Error("poll job failed", sl.Err(err))
		}
	})

	c.Start()
	log.Info("payout-worker started",
		slog.String("calculation_cron", cfg.CalculationCron),
		slog.String("dispatch_cron", cfg.DispatchCron),
		slog.String("poll_cron", cfg.PollCron),
	)

	<-ctx.Done()
	log.Info("shutting down payout-worker")
	<-c.Stop().Done()
	log.Info("payout-worker stopped")
}

func runCalculation(ctx context.Context, cfg *config.Config, calcSvc *payoutServices.CalculationService, log *slog.Logger) {
	now := periodEnd()
	start := now.AddDate(0, 0, -cfg.PeriodDays)

	log.Info("running calculation", slog.String("period_start", start.String()), slog.String("period_end", now.String()))

	n, err := calcSvc.CalculateForPeriod(ctx, start, now)
	if err != nil {
		log.Error("calculation failed", sl.Err(err))

		return
	}

	log.Info("calculation completed", slog.Int("payouts_created", n))
}

// periodEnd returns truncated-to-day UTC now (exclusive upper bound of period).
func periodEnd() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
