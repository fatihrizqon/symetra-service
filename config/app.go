package config

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/middleware"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/route"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Cors     *CORSConfig
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	JWT      *JWTService
}

func Bootstrap(config *BootstrapConfig) {
	config.App.Use(config.Cors.Handler())

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepository := repository.NewUserRepository(config.DB)
	authRepository := repository.NewAuthRepository(config.DB)
	tokenRepository := repository.NewTokenRepository(config.DB)
	coaGroupRepository := repository.NewCOAGroupRepository(config.DB)
	coaSubGroupRepository := repository.NewCOASubGroupRepository(config.DB)
	coaRepository := repository.NewCOARepository(config.DB)
	journalEntryRepository := repository.NewJournalEntryRepository(config.DB)
	customerRepository := repository.NewCustomerRepository(config.DB)
	vendorRepository := repository.NewVendorRepository(config.DB)
	reportRepository := repository.NewReportRepository(config.DB)

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate)
	dashboardService := service.NewDashboardService(userRepository)
	coaGroupService := service.NewCOAGroupService(coaGroupRepository, config.Validate)
	coaSubGroupService := service.NewCOASubGroupService(coaSubGroupRepository, config.Validate)
	coaService := service.NewCOAService(coaRepository, config.Validate)
	journalEntryService := service.NewJournalEntryService(journalEntryRepository, config.Validate)
	customerService := service.NewCustomerService(customerRepository, config.Validate)
	vendorService := service.NewVendorService(vendorRepository, config.Validate)
	reportService := service.NewReportService(reportRepository)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	coaGroupHandler := handler.NewCOAGroupHandler(coaGroupService)
	coaSubGroupHandler := handler.NewCOASubGroupHandler(coaSubGroupService)
	coaHandler := handler.NewCOAHandler(coaService)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryService)
	// Revenue and Expense handlers reuse the same journalEntryService —
	// they only differ in how they scope the type filter and journal number prefix.
	revenueHandler := handler.NewRevenueHandler(journalEntryService)
	expenseHandler := handler.NewExpenseHandler(journalEntryService)
	customerHandler := handler.NewCustomerHandler(customerService)
	vendorHandler := handler.NewVendorHandler(vendorService)
	reportHandler := handler.NewReportHandler(reportService)

	// ── Middleware ─────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:                 config.App,
		UserHandler:         userHandler,
		AuthHandler:         authHandler,
		AuthMiddleware:      authMiddleware,
		DashboardHandler:    dashboardHandler,
		COAGroupHandler:     coaGroupHandler,
		COASubGroupHandler:  coaSubGroupHandler,
		COAHandler:          coaHandler,
		JournalEntryHandler: journalEntryHandler,
		RevenueHandler:      revenueHandler,
		ExpenseHandler:      expenseHandler,
		CustomerHandler:     customerHandler,
		VendorHandler:       vendorHandler,
		ReportHandler:       reportHandler,
	}

	routeConfig.Setup()
}
