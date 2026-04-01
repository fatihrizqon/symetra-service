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
	fiscalYearRepository := repository.NewFiscalYearRepository(config.DB)
	fiscalPeriodRepository := repository.NewFiscalPeriodRepository(config.DB)
	companyRepository := repository.NewCompanyRepository(config.DB)                    // ← NEW
	companyMemberRepository := repository.NewCompanyMemberRepository(config.DB)        // ← NEW
	companyConfigRepository := repository.NewCompanyConfigurationRepository(config.DB) // ← NEW
	customerRepository := repository.NewCustomerRepository(config.DB)                  // ← NEW
	quotationRepository := repository.NewQuotationRepository(config.DB)                // ← NEW
	invoiceRepository := repository.NewInvoiceRepository(config.DB)                    // ← NEW
	purchaseOrderRepository := repository.NewPurchaseOrderRepository(config.DB)        // ← NEW
	billRepository := repository.NewBillRepository(config.DB)                          // ← NEW

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate)
	dashboardService := service.NewDashboardService(userRepository)
	coaGroupService := service.NewCOAGroupService(coaGroupRepository, config.Validate)
	coaSubGroupService := service.NewCOASubGroupService(coaSubGroupRepository, config.Validate)
	coaService := service.NewCOAService(coaRepository, config.Validate)
	journalEntryService := service.NewJournalEntryService(journalEntryRepository, invoiceRepository, config.Validate)
	reportService := service.NewReportService(reportRepository)
	fiscalService := service.NewFiscalService(fiscalYearRepository, fiscalPeriodRepository, config.Validate)
	companyService := service.NewCompanyService( // ← NEW
		companyRepository,
		companyMemberRepository,
		userRepository,
		config.Validate,
	)
	companyConfigService := service.NewCompanyConfigurationService(companyConfigRepository, config.Validate)       // ← NEW
	customerService := service.NewCustomerService(customerRepository, config.Validate)                             // ← NEW
	quotationService := service.NewQuotationService(quotationRepository, companyConfigRepository, config.Validate) // ← NEW
	invoiceService := service.NewInvoiceService(                                                                   // ← NEW
		invoiceRepository, quotationRepository, companyConfigRepository,
		journalEntryRepository, fiscalPeriodRepository, config.Validate,
	)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderRepository, companyConfigRepository, config.Validate) // ← NEW
	billService := service.NewBillService(                                                                                     // ← NEW
		billRepository, purchaseOrderRepository, companyConfigRepository,
		journalEntryRepository, fiscalPeriodRepository, config.Validate,
	)

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
	fiscalHandler := handler.NewFiscalHandler(fiscalService)
	companyHandler := handler.NewCompanyHandler(companyService)                          // ← NEW
	companyConfigHandler := handler.NewCompanyConfigurationHandler(companyConfigService) // ← NEW
	customerHandler := handler.NewCustomerHandler(customerService)                       // ← NEW
	quotationHandler := handler.NewQuotationHandler(quotationService)                    // ← NEW
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)                          // ← NEW
	purchaseOrderHandler := handler.NewPurchaseOrderHandler(purchaseOrderService)        // ← NEW
	billHandler := handler.NewBillHandler(billService)                                   // ← NEW

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)
	companyMiddleware := middleware.NewCompany(companyMemberRepository) // ← NEW

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:               config.App,
		AuthMiddleware:    authMiddleware,
		CompanyMiddleware: companyMiddleware, // ← NEW

		UserHandler:          userHandler,
		AuthHandler:          authHandler,
		DashboardHandler:     dashboardHandler,
		COAGroupHandler:      coaGroupHandler,
		COASubGroupHandler:   coaSubGroupHandler,
		COAHandler:           coaHandler,
		JournalEntryHandler:  journalEntryHandler,
		RevenueHandler:       revenueHandler,
		ExpenseHandler:       expenseHandler,
		ReportHandler:        reportHandler,
		FiscalHandler:        fiscalHandler,
		CompanyHandler:       companyHandler,       // ← NEW
		CompanyConfigHandler: companyConfigHandler, // ← NEW
		CustomerHandler:      customerHandler,      // ← NEW
		QuotationHandler:     quotationHandler,     // ← NEW
		InvoiceHandler:       invoiceHandler,       // ← NEW
		PurchaseOrderHandler: purchaseOrderHandler, // ← NEW
		BillHandler:          billHandler,          // ← NEW
	}

	routeConfig.Setup()
}
