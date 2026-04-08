// ── PATCH untuk file: config/app.go ──────────────────────────────────────────
//
// Diff dari versi existing ke versi baru.
// Perubahan:
//   1. Tambah vendorRepository dan vendorService
//   2. Update NewPurchaseOrderService — tambah param vendorRepo (FIX FRAUD-02)
//   3. Update NewBillService — tambah param vendorRepo (FIX FRAUD-02/03)
//   4. Tambah VendorHandler di RouteConfig
//
// ─────────────────────────────────────────────────────────────────────────────

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
	companyRepository := repository.NewCompanyRepository(config.DB)
	companyMemberRepository := repository.NewCompanyMemberRepository(config.DB)
	companyConfigRepository := repository.NewCompanyConfigurationRepository(config.DB)
	customerRepository := repository.NewCustomerRepository(config.DB)
	// NEW: vendorRepository diperlukan untuk validasi kepemilikan vendor (FRAUD-02)
	vendorRepository := repository.NewVendorRepository(config.DB)
	quotationRepository := repository.NewQuotationRepository(config.DB)
	invoiceRepository := repository.NewInvoiceRepository(config.DB)
	purchaseOrderRepository := repository.NewPurchaseOrderRepository(config.DB)
	billRepository := repository.NewBillRepository(config.DB)

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate)
	dashboardService := service.NewDashboardService(userRepository, companyMemberRepository)
	coaGroupService := service.NewCOAGroupService(coaGroupRepository, config.Validate)
	coaSubGroupService := service.NewCOASubGroupService(coaSubGroupRepository, config.Validate)
	coaService := service.NewCOAService(coaRepository, config.Validate)
	// FIX [BUG-02 & FLOW-06]: Pass fiscalPeriodRepository agar JE service bisa ValidatePeriodOpen dan set FiscalPeriodId
	journalEntryService := service.NewJournalEntryService(journalEntryRepository, invoiceRepository, fiscalPeriodRepository, config.Validate)
	reportService := service.NewReportService(reportRepository)
	fiscalService := service.NewFiscalService(
		fiscalYearRepository,
		fiscalPeriodRepository,
		journalEntryRepository,
		coaRepository,
		coaGroupRepository,
		reportRepository,
		companyConfigRepository,
		config.Validate,
	)
	companyService := service.NewCompanyService(companyRepository, companyMemberRepository, userRepository, config.Validate)
	companyConfigService := service.NewCompanyConfigurationService(companyConfigRepository, config.Validate)
	customerService := service.NewCustomerService(customerRepository, config.Validate)
	// NEW: vendorService
	vendorService := service.NewVendorService(vendorRepository, config.Validate)
	quotationService := service.NewQuotationService(quotationRepository, companyConfigRepository, config.Validate)
	invoiceService := service.NewInvoiceService(
		invoiceRepository, quotationRepository, companyConfigRepository,
		journalEntryRepository, fiscalPeriodRepository, config.Validate,
	)
	// FIX [FRAUD-02]: NewPurchaseOrderService sekarang menerima vendorRepository
	purchaseOrderService := service.NewPurchaseOrderService(
		purchaseOrderRepository,
		companyConfigRepository,
		vendorRepository, // ← tambahan untuk validasi vendor cross-company
		config.Validate,
	)
	// FIX [FRAUD-02/03]: NewBillService sekarang menerima vendorRepository
	billService := service.NewBillService(
		billRepository,
		purchaseOrderRepository,
		companyConfigRepository,
		vendorRepository, // ← tambahan untuk validasi vendor cross-company
		journalEntryRepository,
		fiscalPeriodRepository,
		config.Validate,
	)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	// FIX [BUG-03]: Baca flag production dari config agar cookie Secure=true di HTTPS
	// Tambahkan "app.production": true di config.json untuk environment produksi
	isProduction := config.Config.GetBool("app.production")
	authHandler := handler.NewAuthHandler(authService, isProduction)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	coaGroupHandler := handler.NewCOAGroupHandler(coaGroupService)
	coaSubGroupHandler := handler.NewCOASubGroupHandler(coaSubGroupService)
	coaHandler := handler.NewCOAHandler(coaService)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryService)
	revenueHandler := handler.NewRevenueHandler(journalEntryService)
	expenseHandler := handler.NewExpenseHandler(journalEntryService)
	reportHandler := handler.NewReportHandler(reportService)
	fiscalHandler := handler.NewFiscalHandler(fiscalService)
	companyHandler := handler.NewCompanyHandler(companyService)
	companyConfigHandler := handler.NewCompanyConfigurationHandler(companyConfigService)
	customerHandler := handler.NewCustomerHandler(customerService)
	// NEW: vendorHandler
	vendorHandler := handler.NewVendorHandler(vendorService)
	quotationHandler := handler.NewQuotationHandler(quotationService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	purchaseOrderHandler := handler.NewPurchaseOrderHandler(purchaseOrderService)
	billHandler := handler.NewBillHandler(billService)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)
	companyMiddleware := middleware.NewCompany(companyMemberRepository)

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:               config.App,
		AuthMiddleware:    authMiddleware,
		CompanyMiddleware: companyMiddleware,

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
		CompanyHandler:       companyHandler,
		CompanyConfigHandler: companyConfigHandler,
		CustomerHandler:      customerHandler,
		// NEW: VendorHandler
		VendorHandler:        vendorHandler,
		QuotationHandler:     quotationHandler,
		InvoiceHandler:       invoiceHandler,
		PurchaseOrderHandler: purchaseOrderHandler,
		BillHandler:          billHandler,
	}

	routeConfig.Setup()
}
