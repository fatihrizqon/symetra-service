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
	reportRepository := repository.NewReportRepository(config.DB)

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate)
	dashboardService := service.NewDashboardService(userRepository)
	coaGroupService := service.NewCOAGroupService(coaGroupRepository, config.Validate)
	coaSubGroupService := service.NewCOASubGroupService(coaSubGroupRepository, config.Validate)
	coaService := service.NewCOAService(coaRepository, config.Validate)
	journalEntryService := service.NewJournalEntryService(journalEntryRepository, config.Validate)
	reportService := service.NewReportService(reportRepository)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	coaGroupHandler := handler.NewCOAGroupHandler(coaGroupService)
	coaSubGroupHandler := handler.NewCOASubGroupHandler(coaSubGroupService)
	coaHandler := handler.NewCOAHandler(coaService)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryService)
	revenueHandler := handler.NewRevenueHandler(journalEntryService)
	expenseHandler := handler.NewExpenseHandler(journalEntryService)
	reportHandler := handler.NewReportHandler(reportService)

	authMiddleware := middleware.NewAuth(tokenRepository)

	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,

		UserHandler:         userHandler,
		AuthHandler:         authHandler,
		DashboardHandler:    dashboardHandler,
		COAGroupHandler:     coaGroupHandler,
		COASubGroupHandler:  coaSubGroupHandler,
		COAHandler:          coaHandler,
		JournalEntryHandler: journalEntryHandler,
		RevenueHandler:      revenueHandler,
		ExpenseHandler:      expenseHandler,
		ReportHandler:       reportHandler,
	}

	routeConfig.Setup()
}
