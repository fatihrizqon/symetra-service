package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App                  *fiber.App
	AuthMiddleware       fiber.Handler
	CompanyMiddleware    fiber.Handler
	UserHandler          *handler.UserHandler
	AuthHandler          *handler.AuthHandler
	DashboardHandler     *handler.DashboardHandler
	COAGroupHandler      *handler.COAGroupHandler
	COASubGroupHandler   *handler.COASubGroupHandler
	COAHandler           *handler.COAHandler
	JournalEntryHandler  *handler.JournalEntryHandler
	RevenueHandler       *handler.RevenueHandler
	ExpenseHandler       *handler.ExpenseHandler
	ReportHandler        *handler.ReportHandler
	FiscalHandler        *handler.FiscalHandler
	CompanyHandler       *handler.CompanyHandler
	CompanyConfigHandler *handler.CompanyConfigurationHandler
	CustomerHandler      *handler.CustomerHandler
	VendorHandler        *handler.VendorHandler
	QuotationHandler     *handler.QuotationHandler
	InvoiceHandler       *handler.InvoiceHandler
	PurchaseOrderHandler *handler.PurchaseOrderHandler
	BillHandler          *handler.BillHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)
	rc.App.Get("/api/v1/dashboard/overview", rc.CompanyMiddleware, rc.DashboardHandler.Overview)

	// ── Users ────────────────────────────────────────────────────────────────
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

	// ── Company ───────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/companies", rc.CompanyHandler.CreateCompany)
	rc.App.Get("/api/v1/companies/mine", rc.CompanyHandler.FindMyCompanies)
	rc.App.Get("/api/v1/companies", rc.CompanyHandler.FindAllCompanies)
	rc.App.Get("/api/v1/companies/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("company:read"), rc.CompanyHandler.FindCompanyById)
	rc.App.Put("/api/v1/companies/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("company:update"), rc.CompanyHandler.UpdateCompany)
	rc.App.Delete("/api/v1/companies/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("company:delete"), rc.CompanyHandler.DeleteCompany)
	rc.App.Get("/api/v1/companies/:id/members", rc.CompanyMiddleware, middleware.NewRequirePermission("users:read"), rc.CompanyHandler.FindMembers)
	rc.App.Post("/api/v1/companies/:id/members", rc.CompanyMiddleware, middleware.NewRequirePermission("users:manage"), rc.CompanyHandler.AssignMember)
	rc.App.Put("/api/v1/companies/:id/members/:id/role", rc.CompanyMiddleware, middleware.NewRequirePermission("users:manage"), rc.CompanyHandler.UpdateMemberRole)
	rc.App.Delete("/api/v1/companies/:id/members/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("users:manage"), rc.CompanyHandler.RemoveMember)

	// ── COA ──────────────────────────────────────────────────────────────────
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

	// ── Journal Entries ───────────────────────────────────────────────────────
	rc.App.Post("/api/v1/journal_entries", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Create)
	rc.App.Get("/api/v1/journal_entries", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.JournalEntryHandler.FindAll)
	rc.App.Get("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.JournalEntryHandler.FindById)
	rc.App.Put("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Update)
	rc.App.Delete("/api/v1/journal_entries/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.JournalEntryHandler.Delete)
	rc.App.Put("/api/v1/journal_entries/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.JournalEntryHandler.Post)
	rc.App.Put("/api/v1/journal_entries/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.JournalEntryHandler.Void)

	// ── Revenues ──────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/revenues", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Create)
	rc.App.Get("/api/v1/revenues", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.RevenueHandler.FindAll)
	rc.App.Get("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.RevenueHandler.FindById)
	rc.App.Put("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Update)
	rc.App.Delete("/api/v1/revenues/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.RevenueHandler.Delete)
	rc.App.Put("/api/v1/revenues/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.RevenueHandler.Post)
	rc.App.Put("/api/v1/revenues/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.RevenueHandler.Void)

	// ── Expenses ──────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/expenses", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Create)
	rc.App.Get("/api/v1/expenses", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.ExpenseHandler.FindAll)
	rc.App.Get("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.ExpenseHandler.FindById)
	rc.App.Put("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Update)
	rc.App.Delete("/api/v1/expenses/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.ExpenseHandler.Delete)
	rc.App.Put("/api/v1/expenses/:id/post", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.ExpenseHandler.Post)
	rc.App.Put("/api/v1/expenses/:id/void", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.ExpenseHandler.Void)

	// ── Reports ───────────────────────────────────────────────────────────────
	rc.App.Get("/api/v1/reports/trial-balance", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.TrialBalance)
	rc.App.Get("/api/v1/reports/profit-loss", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.ProfitLoss)
	rc.App.Get("/api/v1/reports/balance-sheet", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.BalanceSheet)
	rc.App.Get("/api/v1/reports/cash-flow", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.CashFlow)
	rc.App.Get("/api/v1/reports/equity-statement", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.EquityStatement)
	// FIX [BUG-01]: Route General Ledger & Journal Book sebelumnya tidak terdaftar — menyebabkan 404 permanen
	rc.App.Get("/api/v1/reports/general-ledger", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.GeneralLedger)
	rc.App.Get("/api/v1/reports/journal-book", rc.CompanyMiddleware, middleware.NewRequirePermission("reports:read"), rc.ReportHandler.JournalBook)

	// ── Fiscal ────────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/fiscal/years", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.CreateFiscalYear)
	rc.App.Get("/api/v1/fiscal/years", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:read"), rc.FiscalHandler.FindAllFiscalYears)
	// ── Sub-routes must be registered BEFORE the generic /:id routes in Fiber ──
	rc.App.Get("/api/v1/fiscal/years/:id/readiness", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.CheckReadiness)
	rc.App.Put("/api/v1/fiscal/years/:id/close", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.CloseFiscalYear)
	rc.App.Post("/api/v1/fiscal/years/:id/opening-balance", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.GenerateOpeningBalance)
	rc.App.Delete("/api/v1/fiscal/years/:id/opening-balance", rc.CompanyMiddleware, middleware.NewRequirePermission("fiscal:manage"), rc.FiscalHandler.DeleteOpeningBalance)
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

	// ── Company Configuration ─────────────────────────────────────────────────
	rc.App.Get("/api/v1/companies/:id/configuration", rc.CompanyMiddleware, middleware.NewRequirePermission("company:read"), rc.CompanyConfigHandler.Get)
	rc.App.Put("/api/v1/companies/:id/configuration", rc.CompanyMiddleware, middleware.NewRequirePermission("company:update"), rc.CompanyConfigHandler.Upsert)

	// ── Customers ─────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/customers", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.CustomerHandler.Create)
	rc.App.Get("/api/v1/customers", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.CustomerHandler.FindAll)
	rc.App.Get("/api/v1/customers/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.CustomerHandler.FindById)
	rc.App.Put("/api/v1/customers/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.CustomerHandler.Update)
	rc.App.Delete("/api/v1/customers/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.CustomerHandler.Delete)
	rc.App.Get("/api/v1/dropdown/customers", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.CustomerHandler.SelectDropdownList)

	// ── Vendors ───────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/vendors", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.VendorHandler.Create)
	rc.App.Get("/api/v1/vendors", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.VendorHandler.FindAll)
	rc.App.Get("/api/v1/vendors/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.VendorHandler.FindById)
	rc.App.Put("/api/v1/vendors/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.VendorHandler.Update)
	rc.App.Delete("/api/v1/vendors/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.VendorHandler.Delete)
	// FIX [CFG-02/BUG-21]: Route dropdown vendor — sebelumnya hilang, menyebabkan UI vendor kosong
	rc.App.Get("/api/v1/dropdown/vendors", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.VendorHandler.SelectDropdownList)

	// ── Quotations ────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/quotations", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Create)
	rc.App.Get("/api/v1/quotations", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.QuotationHandler.FindAll)
	rc.App.Get("/api/v1/quotations/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.QuotationHandler.FindById)
	rc.App.Put("/api/v1/quotations/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Update)
	rc.App.Delete("/api/v1/quotations/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Delete)
	rc.App.Put("/api/v1/quotations/:id/send", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Send)
	rc.App.Put("/api/v1/quotations/:id/accept", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Accept)
	rc.App.Put("/api/v1/quotations/:id/decline", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.QuotationHandler.Decline)

	// ── Invoices ──────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/invoices", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.InvoiceHandler.Create)
	rc.App.Post("/api/v1/quotations/:id/convert", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.InvoiceHandler.CreateFromQuotation)
	rc.App.Get("/api/v1/invoices", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.InvoiceHandler.FindAll)
	rc.App.Get("/api/v1/invoices/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.InvoiceHandler.FindById)
	rc.App.Put("/api/v1/invoices/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.InvoiceHandler.Update)
	rc.App.Delete("/api/v1/invoices/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.InvoiceHandler.Delete)
	rc.App.Put("/api/v1/invoices/:id/confirm", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.InvoiceHandler.Confirm)
	rc.App.Put("/api/v1/invoices/:id/pay", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.InvoiceHandler.MarkPaid)
	rc.App.Put("/api/v1/invoices/:id/cancel", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.InvoiceHandler.Cancel)

	// ── Purchase Orders ───────────────────────────────────────────────────────
	rc.App.Post("/api/v1/purchase-orders", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Create)
	rc.App.Get("/api/v1/purchase-orders", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.PurchaseOrderHandler.FindAll)
	rc.App.Get("/api/v1/purchase-orders/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.PurchaseOrderHandler.FindById)
	rc.App.Put("/api/v1/purchase-orders/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Update)
	rc.App.Delete("/api/v1/purchase-orders/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Delete)
	rc.App.Put("/api/v1/purchase-orders/:id/send", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Send)
	rc.App.Put("/api/v1/purchase-orders/:id/approve", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Approve)
	rc.App.Put("/api/v1/purchase-orders/:id/decline", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.PurchaseOrderHandler.Decline)
	// NEW: Dropdown PO — hanya PO approved yang belum converted
	rc.App.Get("/api/v1/dropdown/purchase-orders", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.PurchaseOrderHandler.SelectDropdownList)

	// ── Bills ─────────────────────────────────────────────────────────────────
	rc.App.Post("/api/v1/bills", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.BillHandler.Create)
	rc.App.Post("/api/v1/bills/from-purchase-order/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.BillHandler.CreateFromPO)
	rc.App.Get("/api/v1/bills", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.BillHandler.FindAll)
	rc.App.Get("/api/v1/bills/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.BillHandler.FindById)
	rc.App.Put("/api/v1/bills/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.BillHandler.Update)
	rc.App.Delete("/api/v1/bills/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:write"), rc.BillHandler.Delete)
	rc.App.Put("/api/v1/bills/:id/confirm", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.BillHandler.Confirm)
	rc.App.Post("/api/v1/bills/:id/pay", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.BillHandler.AddPayment)
	rc.App.Put("/api/v1/bills/:id/cancel", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:post"), rc.BillHandler.Cancel)
	// NEW: Dropdown Bills
	rc.App.Get("/api/v1/dropdown/bills", rc.CompanyMiddleware, middleware.NewRequirePermission("transactions:read"), rc.BillHandler.SelectDropdownList)
}
