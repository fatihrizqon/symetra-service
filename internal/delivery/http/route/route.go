package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	// CompanyMiddleware validates X-Company-ID and injects company_id + role into locals.
	// Apply this to any route that reads/writes company-scoped data.
	CompanyMiddleware fiber.Handler

	UserHandler         *handler.UserHandler
	AuthHandler         *handler.AuthHandler
	DashboardHandler    *handler.DashboardHandler
	COAGroupHandler     *handler.COAGroupHandler
	COASubGroupHandler  *handler.COASubGroupHandler
	COAHandler          *handler.COAHandler
	JournalEntryHandler *handler.JournalEntryHandler
	RevenueHandler      *handler.RevenueHandler
	ExpenseHandler      *handler.ExpenseHandler
	ReportHandler       *handler.ReportHandler
	FiscalHandler       *handler.FiscalHandler
	CompanyHandler      *handler.CompanyHandler // ← NEW
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Post("/api/v1/users", rc.UserHandler.Create) // public registration
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	// All routes below require a valid JWT
	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)
	rc.App.Get("/api/v1/dashboard/overview", rc.DashboardHandler.Overview)

	// ── Users (admin / self) ──────────────────────────────────────────────────
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

	// ── Company (no X-Company-ID needed for create / mine list) ──────────────
	// Create a new company — any authenticated user can do this
	rc.App.Post("/api/v1/companies", rc.CompanyHandler.CreateCompany)
	// List MY companies — returns companies the caller belongs to
	rc.App.Get("/api/v1/companies/mine", rc.CompanyHandler.FindMyCompanies)
	// List ALL companies — superadmin only (permission gate in handler)
	rc.App.Get("/api/v1/companies",
		middleware.NewRequirePermission("*"),
		rc.CompanyHandler.FindAllCompanies,
	)

	// Routes below require X-Company-ID — apply CompanyMiddleware as group prefix
	// Pattern: rc.App.Method(path, rc.CompanyMiddleware, [permMiddleware,] handler)

	// Company detail / settings
	rc.App.Get("/api/v1/companies/:id",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("company:read"),
		rc.CompanyHandler.FindCompanyById,
	)
	rc.App.Put("/api/v1/companies/:id",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("company:update"),
		rc.CompanyHandler.UpdateCompany,
	)
	rc.App.Delete("/api/v1/companies/:id",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("company:delete"),
		rc.CompanyHandler.DeleteCompany,
	)

	// Company members
	rc.App.Get("/api/v1/companies/:id/members",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("users:read"),
		rc.CompanyHandler.FindMembers,
	)
	rc.App.Post("/api/v1/companies/:id/members",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("users:manage"),
		rc.CompanyHandler.AssignMember,
	)
	rc.App.Put("/api/v1/companies/:id/members/:user_id/role",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("users:manage"),
		rc.CompanyHandler.UpdateMemberRole,
	)
	rc.App.Delete("/api/v1/companies/:id/members/:user_id",
		rc.CompanyMiddleware,
		middleware.NewRequirePermission("users:manage"),
		rc.CompanyHandler.RemoveMember,
	)

	// ── COA (company-scoped) ──────────────────────────────────────────────────
	rc.App.Post("/api/v1/coa_groups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAGroupHandler.Create)
	rc.App.Get("/api/v1/coa_groups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_groups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_groups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_groups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_groups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAGroupHandler.SelectDropdownList)

	rc.App.Post("/api/v1/coa_subgroups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COASubGroupHandler.Create)
	rc.App.Get("/api/v1/coa_subgroups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COASubGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_subgroups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COASubGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_subgroups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COASubGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_subgroups/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COASubGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_subgroups", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COASubGroupHandler.SelectDropdownList)

	rc.App.Post("/api/v1/coa", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAHandler.Create)
	rc.App.Get("/api/v1/coa", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAHandler.FindAll)
	rc.App.Get("/api/v1/coa/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAHandler.FindById)
	rc.App.Put("/api/v1/coa/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAHandler.Update)
	rc.App.Delete("/api/v1/coa/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:manage"), rc.COAHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa", rc.CompanyMiddleware, middleware.NewRequirePermission("coa:read"), rc.COAHandler.SelectDropdownList)

	// ── Journal Entries (company-scoped) ──────────────────────────────────────
	rc.App.Post("/api/v1/journal_entries", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Create)
	rc.App.Get("/api/v1/journal_entries", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.JournalEntryHandler.FindAll)
	rc.App.Get("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.JournalEntryHandler.FindById)
	rc.App.Put("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Update)
	rc.App.Delete("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Delete)
	rc.App.Put("/api/v1/journal_entries/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.JournalEntryHandler.Post)
	rc.App.Put("/api/v1/journal_entries/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.JournalEntryHandler.Void)

	// ── Revenues (company-scoped) ─────────────────────────────────────────────
	rc.App.Post("/api/v1/revenues", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Create)
	rc.App.Get("/api/v1/revenues", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.RevenueHandler.FindAll)
	rc.App.Get("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.RevenueHandler.FindById)
	rc.App.Put("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Update)
	rc.App.Delete("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Delete)
	rc.App.Put("/api/v1/revenues/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.RevenueHandler.Post)
	rc.App.Put("/api/v1/revenues/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.RevenueHandler.Void)

	// ── Expenses (company-scoped) ─────────────────────────────────────────────
	rc.App.Post("/api/v1/expenses", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Create)
	rc.App.Get("/api/v1/expenses", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.ExpenseHandler.FindAll)
	rc.App.Get("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.ExpenseHandler.FindById)
	rc.App.Put("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Update)
	rc.App.Delete("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Delete)
	rc.App.Put("/api/v1/expenses/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.ExpenseHandler.Post)
	rc.App.Put("/api/v1/expenses/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.ExpenseHandler.Void)

	// ── Reports (company-scoped) ──────────────────────────────────────────────
	rc.App.Get("/api/v1/reports/trial-balance", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.TrialBalance)
	rc.App.Get("/api/v1/reports/profit-loss", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.ProfitLoss)
	rc.App.Get("/api/v1/reports/balance-sheet", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.BalanceSheet)
	rc.App.Get("/api/v1/reports/cash-flow", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.CashFlow)
	rc.App.Get("/api/v1/reports/equity-statement", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.EquityStatement)

	// ── Fiscal (company-scoped) ───────────────────────────────────────────────
	rc.App.Post("/api/v1/fiscal/years", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.CreateFiscalYear)
	rc.App.Get("/api/v1/fiscal/years", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindAllFiscalYears)
	rc.App.Get("/api/v1/fiscal/years/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindFiscalYearById)
	rc.App.Put("/api/v1/fiscal/years/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.UpdateFiscalYear)
	rc.App.Put("/api/v1/fiscal/years/:id/activate", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.ActivateFiscalYear)
	rc.App.Delete("/api/v1/fiscal/years/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.DeleteFiscalYear)
	rc.App.Get("/api/v1/fiscal/periods", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindAllPeriods)
	rc.App.Get("/api/v1/fiscal/periods/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindPeriodById)
	rc.App.Put("/api/v1/fiscal/periods/:id/close", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.ClosePeriod)
	rc.App.Put("/api/v1/fiscal/periods/:id/reopen", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.ReopenPeriod)
	rc.App.Put("/api/v1/fiscal/periods/:id/lock", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.LockPeriod)
	rc.App.Get("/api/v1/fiscal/periods/:id/logs", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindPeriodLogs)
}
