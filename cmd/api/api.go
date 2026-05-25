package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/adapters/rabbitmq"
	redisAdapter "github.com/nlsnnn/berezhok/internal/adapters/redis"
	"github.com/nlsnnn/berezhok/internal/adapters/s3/yandex"
	smsAdapter "github.com/nlsnnn/berezhok/internal/adapters/sms"
	"github.com/nlsnnn/berezhok/internal/adapters/yookassa"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	adminHandlers "github.com/nlsnnn/berezhok/internal/modules/admin/handlers"
	adminRepos "github.com/nlsnnn/berezhok/internal/modules/admin/repository"
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
	orderAdapters "github.com/nlsnnn/berezhok/internal/modules/order/adapters"
	orderHandlers "github.com/nlsnnn/berezhok/internal/modules/order/handlers"
	orderRepos "github.com/nlsnnn/berezhok/internal/modules/order/repository"
	orderServices "github.com/nlsnnn/berezhok/internal/modules/order/service"
	partneradapters "github.com/nlsnnn/berezhok/internal/modules/partner/adapters"
	partnerHandlers "github.com/nlsnnn/berezhok/internal/modules/partner/handlers"
	partnerRepos "github.com/nlsnnn/berezhok/internal/modules/partner/repository"
	partnerServices "github.com/nlsnnn/berezhok/internal/modules/partner/service"
	paymentHandlers "github.com/nlsnnn/berezhok/internal/modules/payment/handlers"
	paymentRepos "github.com/nlsnnn/berezhok/internal/modules/payment/repository"
	paymentServices "github.com/nlsnnn/berezhok/internal/modules/payment/service"
	payoutHandlers "github.com/nlsnnn/berezhok/internal/modules/payout/handlers"
	payoutRepos "github.com/nlsnnn/berezhok/internal/modules/payout/repository"
	payoutServices "github.com/nlsnnn/berezhok/internal/modules/payout/service"
	reviewHandlers "github.com/nlsnnn/berezhok/internal/modules/review/handlers"
	reviewRepos "github.com/nlsnnn/berezhok/internal/modules/review/repository"
	reviewServices "github.com/nlsnnn/berezhok/internal/modules/review/service"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
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
		AllowedOrigins: []string{"http://localhost", "http://localhost:5173", "http://localhost:3000", "http://localhost:8000", "http://localhost:64234"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	// Shared infrastructure
	queries := sqlc.New(app.pool)
	v := validator.New()
	jwtService := jwt.NewTokenService([]byte("supersecretkey"))
	notificationPublisher := rabbitmq.NewNotificationPublisher(app.rabbitmq)
	orderPublisher := rabbitmq.NewOrderPublisher(app.rabbitmq)

	// Partner module - adapters
	notificationAdapter := partneradapters.NewNotificationAdapter(notificationPublisher)

	// Partner module — repositories
	partnerRepo := partnerRepos.NewPartnerRepo(queries)
	employeeRepo := partnerRepos.NewEmployeeRepo(queries)
	appRepo := partnerRepos.NewApplicationRepo(queries)
	locationRepo := partnerRepos.NewLocationRepo(queries)
	pinsRepo := partnerRepos.NewPinsRepo(queries, app.pool)

	// Partner module — services
	partnerSvc := partnerServices.NewPartnerService(partnerRepo, employeeRepo)
	employeeSvc := partnerServices.NewEmployeeService(employeeRepo, partnerRepo, locationRepo, notificationAdapter)
	locationSvc := partnerServices.NewLocationService(locationRepo)
	pinsSvc := partnerServices.NewPinsService(pinsRepo, locationRepo)
	appSvc := partnerServices.NewApplicationService(appRepo, partnerSvc, employeeSvc, locationSvc, notificationAdapter)

	// Partner module — handlers
	partHandler := partnerHandlers.NewPartnerHandler(partnerSvc, app.log)
	appHandler := partnerHandlers.NewApplicationHandler(app.log, appSvc)
	employeeHandler := partnerHandlers.NewEmployeeHandler(app.log, employeeSvc)
	locationHandler := partnerHandlers.NewLocationHandler(app.log, v, locationSvc, partnerSvc)
	pinsHandler := partnerHandlers.NewPinsHandler(app.log, v, pinsSvc)

	// Catalog module — repositories
	boxRepo := catalogRepos.NewBoxRepo(queries)

	// Catalog module — services
	boxSvc := catalogServices.NewBoxService(boxRepo, locationSvc, partnerSvc)

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
	orderRepo := orderRepos.NewOrderRepo(queries, app.pool)
	orderBoxProvider := orderAdapters.NewCatalogBoxProvider(boxSvc)
	orderEventAdapter := orderAdapters.NewOrderEventAdapter(orderPublisher)

	// Payment module
	orderStatusUpdater := orderServices.NewOrderStatusUpdater(orderRepo, app.log, orderEventAdapter)
	yookassaAdapter := yookassa.NewAdapter(yookassa.New(app.cfg.Yookassa))
	paymentRepo := paymentRepos.NewPaymentRepo(queries)
	paymentSvc := paymentServices.NewPaymentService(paymentRepo, yookassaAdapter, orderStatusUpdater)
	webhookHandler := paymentHandlers.NewWebhookHandler(paymentSvc, app.log, v)

	// Order module — services
	orderSvc := orderServices.NewOrderService(orderRepo, orderBoxProvider, paymentSvc, app.log, orderEventAdapter)

	// Order module — handlers
	orderHandler := orderHandlers.NewOrderHandler(orderSvc, app.log, v)

	// Payout module
	payoutRepo := payoutRepos.NewPayoutRepo(queries, app.pool)
	payoutAdapter := yookassa.NewPayoutAdapter(yookassa.NewPayoutHandlerClient(app.cfg.YookassaPayout))
	payoutDestSvc := payoutServices.NewDestinationService(payoutRepo, payoutAdapter)
	payoutQuerySvc := payoutServices.NewQueryService(payoutRepo)
	payoutDestHandler := payoutHandlers.NewDestinationHandler(payoutDestSvc, v, app.log)
	payoutListHandler := payoutHandlers.NewPayoutHandler(payoutQuerySvc, app.log)

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
	adminRepo := adminRepos.NewAdminRepo(queries)
	partnerAuthSvc := authServices.NewPartnerAuthenticator(employeeRepo, jwtService)
	customerAuthSvc := authServices.NewCustomerAuthenticator(customerRepo, jwtService, smsSvc)
	adminAuthSvc := authServices.NewAdminAuthenticator(adminRepo, jwtService)
	authHandler := authHandlers.NewAuthHandler(v, app.log, partnerAuthSvc, customerAuthSvc, adminAuthSvc)
	adminApplicationHandler := adminHandlers.NewApplicationHandler(app.log, v, appSvc, adminRepo)
	adminOpsHandler := adminHandlers.NewOpsHandler(app.log, v, queries, adminRepo, notificationAdapter)

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
		r.Post("/auth/partner/login", authHandler.PartnerLogin)
		r.Post("/auth/customer/send-code", authHandler.CustomerSendCode)
		r.Post("/auth/customer/login", authHandler.CustomerLogin)
		r.Post("/auth/admin/login", authHandler.AdminLogin)

		// Application
		r.Post("/applications", appHandler.Create)

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

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerPasswordChange))
				r.Post("/partner/change-password", partHandler.ChangePassword)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerOrdersView))

				// Orders
				r.Get("/partner/orders", orderHandler.ListPartnerOrders)
				r.Get("/partner/orders/by-code/{pickup_code}", orderHandler.GetPartnerOrderByPickupCode)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerOrdersPickup))

				r.Post("/partner/orders/{order_id}/pickup", orderHandler.PartnerPickupOrder)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerProfileManage))
				r.Get("/partner/profile", partHandler.Profile)
				r.Patch("/partner/profile", partHandler.UpdateProfile)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerLegalInfoManage))
				r.Post("/partner/legal-info", partHandler.AddLegalInfo)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerDashboardView))
				r.Get("/partner/dashboard", partHandler.Dashboard)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerStatsView))
				r.Get("/partner/stats", partHandler.Stats)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerEmployeesManage))
				// Employees
				r.Get("/partner/employees", employeeHandler.List)
				r.Post("/partner/employees", employeeHandler.Create)
				r.Delete("/partner/employees/{id}", employeeHandler.Delete)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerLocationsManage))
				// Location
				r.Get("/partner/locations", locationHandler.List)
				r.Post("/partner/locations", locationHandler.Create)
				// Pins
				r.Get("/partner/pins", pinsHandler.ListAvailable)
				r.Get("/partner/locations/{id}/pins", pinsHandler.GetLocationPins)
				r.Put("/partner/locations/{id}/pins", pinsHandler.UpdateLocationPins)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerBoxesManage))
				// Surprise Box
				r.Post("/partner/boxes", boxHandler.Create)
				r.Get("/partner/boxes/{id}", boxHandler.GetByID)
				r.Put("/partner/boxes/{id}", boxHandler.Update)
				r.Delete("/partner/boxes/{id}", boxHandler.Delete)
				r.Get("/partner/boxes", boxHandler.GetAllByPartnerID)
				r.Get("/locations/{location_id}/boxes", boxHandler.GetAllByLocationID)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerMediaUpload))
				// Media
				r.Post("/media/upload", mediaHandler.Upload)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerPayoutsView))
				r.Get("/partner/payouts", payoutListHandler.List)
				r.Get("/partner/payouts/{id}", payoutListHandler.GetByID)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequirePermission(authz.PermissionPartnerPayoutsManage))
				r.Get("/partner/payout-destination", payoutDestHandler.Get)
				r.Put("/partner/payout-destination", payoutDestHandler.Upsert)
				r.Get("/partner/sbp-banks", payoutDestHandler.GetSBPBanks)
			})
		})

		// == Admin Routes ==
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth("admin"))

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminPermission(authz.PermissionAdminOpsView))
				r.Get("/admin/pins", adminOpsHandler.ListLocationPins)
				r.Get("/admin/me", adminOpsHandler.Me)
				r.Get("/admin/applications", adminApplicationHandler.List)
				r.Get("/admin/applications/{id}", adminApplicationHandler.GetByID)
				r.Get("/admin/audit", adminOpsHandler.ListAudit)
				r.Get("/admin/partners", adminOpsHandler.ListPartners)
				r.Get("/admin/partners/{id}", adminOpsHandler.GetPartner)
				r.Get("/admin/locations", adminOpsHandler.ListLocations)
				r.Get("/admin/locations/{id}", adminOpsHandler.GetLocation)
				r.Get("/admin/boxes", adminOpsHandler.ListBoxes)
				r.Get("/admin/boxes/{id}", adminOpsHandler.GetBox)
				r.Get("/admin/customers", adminOpsHandler.ListCustomers)
				r.Get("/admin/customers/{id}", adminOpsHandler.GetCustomer)
				r.Get("/admin/orders", adminOpsHandler.ListOrders)
				r.Get("/admin/orders/{id}", adminOpsHandler.GetOrder)
				r.Get("/admin/payments", adminOpsHandler.ListPayments)
				r.Get("/admin/payments/{id}", adminOpsHandler.GetPayment)
				r.Get("/admin/payments/{id}/events", adminOpsHandler.ListPaymentEvents)
				r.Get("/admin/stats", adminOpsHandler.Stats)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminPermission(authz.PermissionAdminOpsManage))
				r.Delete("/admin/applications/{id}", adminApplicationHandler.Delete)
				r.Post("/admin/applications/{id}/approve", adminApplicationHandler.Approve)
				r.Post("/admin/applications/{id}/reject", adminApplicationHandler.Reject)
				r.Patch("/admin/partners/{id}", adminOpsHandler.UpdatePartner)
				r.Post("/admin/partners/{partner_id}/legal-info/verify", adminOpsHandler.VerifyPartnerLegalInfo)
				r.Post("/admin/partners/{partner_id}/legal-info/reject", adminOpsHandler.RejectPartnerLegalInfo)
				r.Patch("/admin/locations/{id}/status", adminOpsHandler.UpdateLocationStatus)
				r.Patch("/admin/boxes/{id}/status", adminOpsHandler.UpdateBoxStatus)
				// Pins management
				r.Post("/admin/pins", adminOpsHandler.CreateLocationPin)
				r.Delete("/admin/pins/{code}", adminOpsHandler.DeleteLocationPin)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminPermission(authz.PermissionAdminAdminsManage))
				r.Get("/admin/admins", adminOpsHandler.ListAdmins)
				r.Post("/admin/admins", adminOpsHandler.CreateAdmin)
				r.Patch("/admin/admins/{id}", adminOpsHandler.UpdateAdmin)
				r.Post("/admin/admins/{id}/deactivate", adminOpsHandler.DeactivateAdmin)
			})
		})

		// == Webhook Routes ==
		r.Group(func(r chi.Router) {
			r.Use(webhookMiddleware.IPFilterMiddleware)
			r.Post("/webhooks/yookassa", webhookHandler.Yookassa)
		})
	})

	r.Route("/api/internal/v1", func(r chi.Router) {
		r.Get("/orders/{order_id}/chat-access", orderHandler.InternalChatAccess)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	app.log.Info("starting server", slog.String("address", app.cfg.Address))

	app.srv = &http.Server{
		Addr:         app.cfg.Address,
		Handler:      h,
		ReadTimeout:  app.cfg.Timeout,
		WriteTimeout: app.cfg.Timeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}

	app.log.Info("server has started", slog.String("address", app.cfg.Address))

	return app.srv.ListenAndServe()
}

func (app *application) shutdown(ctx context.Context) error {
	if app.srv == nil {
		app.log.Warn("shutdown called before server initialization")
		return nil
	}

	if err := app.srv.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.log.Info("server has stopped gracefully")
	return nil
}

type application struct {
	cfg      *config.Config
	pool     *pgxpool.Pool
	log      *slog.Logger
	s3       *yandex.Storage
	redis    *redis.Client
	rabbitmq *rabbitmq.Client
	srv      *http.Server
}
