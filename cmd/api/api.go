package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	redisAdapter "github.com/nlsnnn/berezhok/internal/adapters/redis"
	"github.com/nlsnnn/berezhok/internal/adapters/s3/yandex"
	smsAdapter "github.com/nlsnnn/berezhok/internal/adapters/sms"
	"github.com/nlsnnn/berezhok/internal/adapters/yookassa"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	authHandlers "github.com/nlsnnn/berezhok/internal/modules/auth/handlers"
	authServices "github.com/nlsnnn/berezhok/internal/modules/auth/service"
	catalogHandlers "github.com/nlsnnn/berezhok/internal/modules/catalog/handlers"
	catalogRepos "github.com/nlsnnn/berezhok/internal/modules/catalog/repository"
	catalogServices "github.com/nlsnnn/berezhok/internal/modules/catalog/service"
	customerHandlers "github.com/nlsnnn/berezhok/internal/modules/customer/handlers"
	customerRepos "github.com/nlsnnn/berezhok/internal/modules/customer/repository"
	customerServices "github.com/nlsnnn/berezhok/internal/modules/customer/service"
	mediaHandlers "github.com/nlsnnn/berezhok/internal/modules/media/handlers"
	mediaRepos "github.com/nlsnnn/berezhok/internal/modules/media/repository"
	mediaServices "github.com/nlsnnn/berezhok/internal/modules/media/service"
	orderHandlers "github.com/nlsnnn/berezhok/internal/modules/order/handlers"
	orderRepos "github.com/nlsnnn/berezhok/internal/modules/order/repository"
	orderServices "github.com/nlsnnn/berezhok/internal/modules/order/service"
	partnerHandlers "github.com/nlsnnn/berezhok/internal/modules/partner/handlers"
	partnerRepos "github.com/nlsnnn/berezhok/internal/modules/partner/repository"
	partnerServices "github.com/nlsnnn/berezhok/internal/modules/partner/service"
	paymentHandlers "github.com/nlsnnn/berezhok/internal/modules/payment/handlers"
	paymentRepos "github.com/nlsnnn/berezhok/internal/modules/payment/repository"
	paymentServices "github.com/nlsnnn/berezhok/internal/modules/payment/service"
	reviewHandlers "github.com/nlsnnn/berezhok/internal/modules/review/handlers"
	reviewRepos "github.com/nlsnnn/berezhok/internal/modules/review/repository"
	reviewServices "github.com/nlsnnn/berezhok/internal/modules/review/service"
	"github.com/nlsnnn/berezhok/internal/shared/config"
	"github.com/nlsnnn/berezhok/internal/shared/jwt"
	middlewares "github.com/nlsnnn/berezhok/internal/shared/middleware"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost", "http://localhost:5173", "http://localhost:3000", "http://localhost:8000", "http://localhost:63333"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Shared infrastructure
	queries := sqlc.New(app.pool)
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
	locationSvc := partnerServices.NewLocationService(locationRepo)
	appSvc := partnerServices.NewApplicationService(appRepo, partnerSvc, employeeSvc, locationSvc)

	// Partner module — handlers
	partHandler := partnerHandlers.NewPartnerHandler(partnerSvc, app.log)
	appHandler := partnerHandlers.NewApplicationHandler(app.log, appSvc)
	locationHandler := partnerHandlers.NewLocationHandler(app.log, v, locationSvc, partnerSvc)

	// Catalog module — repositories
	boxRepo := catalogRepos.NewBoxRepo(queries)

	// Catalog module — services
	boxSvc := catalogServices.NewBoxService(boxRepo, locationSvc)

	// Catalog module — handlers
	boxHandler := catalogHandlers.NewBoxHandler(boxSvc, app.log, v)

	// Media module — repositories
	mediaRepo := mediaRepos.NewMediaRepo(queries)

	// Media module — services
	mediaSvc := mediaServices.NewMediaService(app.s3, mediaRepo, app.log)

	// Media module — handlers
	mediaHandler := mediaHandlers.NewMediaHandler(mediaSvc, app.log)

	// Customer module — repositories
	customerRepo := customerRepos.NewUserRepo(queries)
	customerLocationRepo := customerRepos.NewLocationRepo(queries)

	// Customer module — services
	customerSvc := customerServices.NewCustomerService(customerRepo)
	customerLocationSvc := customerServices.NewLocationService(customerLocationRepo)

	// Order module — repositories
	orderRepo := orderRepos.NewOrderRepo(queries)

	// Payment module
	orderStatusUpdater := orderServices.NewOrderStatusUpdater(orderRepo, app.log)
	yookassaAdapter := yookassa.NewAdapter(yookassa.New(app.cfg.Yookassa))
	paymentRepo := paymentRepos.NewPaymentRepo(queries)
	paymentSvc := paymentServices.NewPaymentService(paymentRepo, yookassaAdapter, orderStatusUpdater)
	webhookHandler := paymentHandlers.NewWebhookHandler(paymentSvc, app.log, v)

	// Order module — services
	orderSvc := orderServices.NewOrderService(orderRepo, boxSvc, paymentSvc, app.log)

	// Order module — handlers
	orderHandler := orderHandlers.NewOrderHandler(orderSvc, app.log, v)

	// Review module
	reviewRepo := reviewRepos.NewReviewRepo(queries)
	reviewSvc := reviewServices.NewReviewService(reviewRepo, orderSvc)
	reviewHandler := reviewHandlers.NewReviewHandler(reviewSvc, app.log, v)

	// Customer module — handlers
	customerHandler := customerHandlers.NewCustomerHandler(customerSvc, app.log, v)
	customerLocationHandler := customerHandlers.NewLocationHandler(customerLocationSvc, app.log)

	// SMS module
	smsStorage := redisAdapter.NewSMSStorage(app.redis)
	smsSender := smsAdapter.NewConsoleSender()
	smsSvc := authServices.NewSMSService(smsStorage, smsSender)

	// Auth module
	partnerAuthSvc := authServices.NewPartnerAuthenticator(employeeRepo, jwtService)
	customerAuthSvc := authServices.NewCustomerAuthenticator(customerRepo, jwtService, smsSvc)
	authHandler := authHandlers.NewAuthHandler(v, app.log, partnerAuthSvc, customerAuthSvc)

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)
	webhookMiddleware := middlewares.NewWebhookMiddleware([]string{
		"127.0.0.1/32", // IPv4 localhost
		"::1/128",      // IPv6 localhost
		"185.71.76.0/27",
		"185.71.77.0/27",
		"77.75.153.0/25",
		"77.75.156.11",
		"77.75.156.35",
		"77.75.154.128/25",
		"2a02:5180::/32",
	}, app.log) // IP-адреса Yookassa

	r.Route("/api/v1/", func(r chi.Router) {
		// == Public Routes ==

		// Auth
		r.Post("/partner/auth/login", authHandler.PartnerLogin)
		r.Post("/auth/customer/send-code", authHandler.CustomerSendCode)
		r.Post("/auth/customer/login", authHandler.CustomerLogin)

		// Application
		r.Post("/applications", appHandler.Create)
		r.Get("/applications", appHandler.List)
		r.Get("/applications/{id}", appHandler.GetByID)
		r.Delete("/applications/{id}", appHandler.Delete)
		r.Post("/applications/{id}/approve", appHandler.Approve)
		r.Post("/applications/{id}/reject", appHandler.Reject)

		// == Customer Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("customer"))

			// Profile
			r.Get("/customer/profile", customerHandler.GetProfile)
			r.Patch("/customer/profile", customerHandler.UpdateProfile)

			// Locations
			r.Get("/customer/locations", customerLocationHandler.SearchLocations)
			r.Get("/customer/locations/{location_id}", customerLocationHandler.GetLocationDetails)

			// Orders
			r.Post("/customer/orders", orderHandler.CreateOrder)
			r.Get("/customer/orders", orderHandler.ListOrders)
			r.Get("/customer/orders/{order_id}", orderHandler.GetOrder)
			r.Post("/customer/orders/{order_id}/confirm-pickup", orderHandler.ConfirmPickup)
			r.Post("/customer/orders/{order_id}/dispute", orderHandler.CreateDispute)

			// Reviews
			r.Post("/customer/reviews", reviewHandler.CreateReview)
			r.Get("/customer/locations/{location_id}/reviews", reviewHandler.ListLocationReviews)
		})

		// == Partner Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("partner"))

			r.Post("/partner/change-password", partHandler.ChangePassword)
			r.Get("/partner/profile", partHandler.Profile)

			// Location
			r.Get("/partner/locations", locationHandler.List)
			r.Post("/partner/locations", locationHandler.Create)

			// Surprise Box
			r.Post("/partner/boxes", boxHandler.Create)
			r.Get("/partner/boxes/{id}", boxHandler.GetByID)
			r.Put("/partner/boxes/{id}", boxHandler.Update)
			r.Delete("/partner/boxes/{id}", boxHandler.Delete)
			r.Get("/partner/boxes", boxHandler.GetAllByPartnerID)
			r.Get("/locations/{location_id}/boxes", boxHandler.GetAllByLocationID)

			// Media
			r.Post("/media/upload", mediaHandler.Upload)

			// Orders
			r.Get("/partner/orders/by-code/{pickup_code}", orderHandler.GetPartnerOrderByPickupCode)
			r.Post("/partner/orders/{order_id}/pickup", orderHandler.PartnerPickupOrder)
		})

		// == Admin Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("admin"))
		})

		// == Webhook Routes ==
		r.Group(func(r chi.Router) {
			r.Use(webhookMiddleware.IPFilterMiddleware)
			r.Post("/webhooks/yookassa", webhookHandler.Yookassa)
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
	cfg   *config.Config
	pool  *pgxpool.Pool
	log   *slog.Logger
	s3    *yandex.Storage
	redis *redis.Client
}
