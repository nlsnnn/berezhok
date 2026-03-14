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
	partnerServices "github.com/nlsnnn/berezhok/internal/modules/partner/service"
	partnerRepos "github.com/nlsnnn/berezhok/internal/repository/partner"
	"github.com/nlsnnn/berezhok/internal/shared/config"
	"github.com/nlsnnn/berezhok/internal/shared/jwt"
	middlewares "github.com/nlsnnn/berezhok/internal/shared/middleware"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
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

	// General
	queries := sqlc.New(app.db)
	validator := validator.New()
	jwtService := jwt.NewTokenService([]byte("supersecretkey"))

	// Partner Module
	partRepo := sqlc.New(app.db)
	partService := partnerServices.NewPartnerService(partRepo)
	partHandler := partnerHandlers.NewPartnerHandler(partService, app.log)

	// Employee
	employeeRepo := sqlc.New(app.db)
	employeeService := partnerServices.NewEmployeeService(employeeRepo)

	// Application
	appRepo := sqlc.New(app.db)
	appSvc := partnerServices.NewApplicationService(appRepo, partService, employeeService)
	appHandler := partnerHandlers.NewApplicationHandler(app.log, appSvc, partService)

	// Location
	// locationRepo := sqlc.New(app.db)
	locationRepo := partnerRepos.NewLocationRepo(queries)
	locationService := partnerServices.NewLocationService(appRepo, locationRepo)
	locationHandler := partnerHandlers.NewLocationHandler(app.log, validator, &locationService, partService)

	// Auth Module
	partnerAuthService := authServices.NewPartnerAuthenticator(partRepo, jwtService)
	authHandler := authHandlers.NewAuthHandler(validator, app.log, partnerAuthService)

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	r.Route("/api/v1/", func(r chi.Router) {
		// == Public Routes ==

		// Auth
		r.Post("/partner/auth/login", authHandler.PartnerLogin)
		// r.Post("/admin/auth/login", authHandler.AdminLogin)

		// Application
		r.Post("/applications", appHandler.Create)
		r.Get("/applications", appHandler.List)
		r.Get("/applications/{id}", appHandler.GetByID)
		r.Delete("/applications/{id}", appHandler.Delete)
		r.Post("/applications/{id}/approve", appHandler.Approve)
		r.Post("/applications/{id}/reject", appHandler.Reject)

		// == Customer Routes ==

		// == Partner Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("partner"))

			r.Post("/partner/change-password", partHandler.ChangePassword)
			r.Get("/partner/profile", partHandler.Profile)
			// r.Put("/partner/profile", partHandler.UpdateProfile)
			// r.Get("/partner/employees", partHandler.ListEmployees)
			// r.Post("/partner/employees", partHandler.CreateEmployee)
			// r.Delete("/partner/employees/{id}", partHandler.DeleteEmployee)

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
