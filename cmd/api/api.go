package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	authHandlers "github.com/nlsnnn/berezhok/internal/modules/auth/handlers"
	authServices "github.com/nlsnnn/berezhok/internal/modules/auth/service"
	partnerHandlers "github.com/nlsnnn/berezhok/internal/modules/partner/handlers"
	partnerRepos "github.com/nlsnnn/berezhok/internal/modules/partner/repository"
	partnerServices "github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/config"
	"github.com/nlsnnn/berezhok/internal/shared/jwt"
	middlewares "github.com/nlsnnn/berezhok/internal/shared/middleware"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Shared infrastructure
	queries := sqlc.New(app.db)
	v := validator.New()
	jwtService := jwt.NewTokenService([]byte("supersecretkey"))

	// Partner module — repositories
	partnerRepo := partnerRepos.NewPartnerRepo(queries)
	employeeRepo := partnerRepos.NewEmployeeRepo(queries)
	appRepo := partnerRepos.NewApplicationRepo(queries)
	locationRepo := partnerRepos.NewLocationRepo(queries)

	// Partner module — services
	partnerSvc := partnerServices.NewPartnerService(partnerRepo, employeeRepo)
	employeeSvc := partnerServices.NewEmployeeService(employeeRepo)
	appSvc := partnerServices.NewApplicationService(appRepo, partnerSvc, employeeSvc)
	locationSvc := partnerServices.NewLocationService(locationRepo)

	// Partner module — handlers
	partHandler := partnerHandlers.NewPartnerHandler(partnerSvc, app.log)
	appHandler := partnerHandlers.NewApplicationHandler(app.log, appSvc)
	locationHandler := partnerHandlers.NewLocationHandler(app.log, v, &locationSvc, partnerSvc)

	// Auth module
	partnerAuthSvc := authServices.NewPartnerAuthenticator(employeeRepo, jwtService)
	authHandler := authHandlers.NewAuthHandler(v, app.log, partnerAuthSvc)

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	r.Route("/api/v1/", func(r chi.Router) {
		// == Public Routes ==

		// Auth
		r.Post("/partner/auth/login", authHandler.PartnerLogin)

		// Application
		r.Post("/applications", appHandler.Create)
		r.Get("/applications", appHandler.List)
		r.Get("/applications/{id}", appHandler.GetByID)
		r.Delete("/applications/{id}", appHandler.Delete)
		r.Post("/applications/{id}/approve", appHandler.Approve)
		r.Post("/applications/{id}/reject", appHandler.Reject)

		// == Partner Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("partner"))

			r.Post("/partner/change-password", partHandler.ChangePassword)
			r.Get("/partner/profile", partHandler.Profile)

			// Location
			r.Get("/partner/locations", locationHandler.List)
			r.Post("/partner/locations", locationHandler.Create)
		})

		// == Admin Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("admin"))
		})
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
