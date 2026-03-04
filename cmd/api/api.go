package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}

func (app *application) run(log *slog.Logger, h http.Handler) error {
	log.Info("starting server", slog.String("address", app.cfg.Address))

	srv := &http.Server{
		Addr:         app.cfg.Address,
		Handler:      h,
		ReadTimeout:  app.cfg.HTTPServer.Timeout,
		WriteTimeout: app.cfg.Timeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}

	log.Info("server has started", slog.String("address", app.cfg.Address))

	return srv.ListenAndServe()
}

type application struct {
	cfg *config.Config
	db  *pgx.Conn
	log *slog.Logger
}
