package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App                 *fiber.App
	UserHandler         *handler.UserHandler
	AuthHandler         *handler.AuthHandler
	AuthMiddleware      fiber.Handler
	DashboardHandler    *handler.DashboardHandler
	COAGroupHandler     *handler.COAGroupHandler
	COASubGroupHandler  *handler.COASubGroupHandler
	COAHandler          *handler.COAHandler
	JournalEntryHandler *handler.JournalEntryHandler
	CustomerHandler     *handler.CustomerHandler
	VendorHandler       *handler.VendorHandler
	ReportHandler       *handler.ReportHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)

	// COA Groups
	rc.App.Post("/api/v1/coa_groups", rc.COAGroupHandler.Create)
	rc.App.Get("/api/v1/coa_groups", rc.COAGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_groups/:id", rc.COAGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_groups/:id", rc.COAGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_groups/:id", rc.COAGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_groups", rc.COAGroupHandler.SelectDropdownList)

	// COA Subgroups
	rc.App.Post("/api/v1/coa_subgroups", rc.COASubGroupHandler.Create)
	rc.App.Get("/api/v1/coa_subgroups", rc.COASubGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_subgroups", rc.COASubGroupHandler.SelectDropdownList)

	// COA
	rc.App.Post("/api/v1/coa", rc.COAHandler.Create)
	rc.App.Get("/api/v1/coa", rc.COAHandler.FindAll)
	rc.App.Get("/api/v1/coa/:id", rc.COAHandler.FindById)
	rc.App.Put("/api/v1/coa/:id", rc.COAHandler.Update)
	rc.App.Delete("/api/v1/coa/:id", rc.COAHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa", rc.COAHandler.SelectDropdownList)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)
	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)

	// Dashboard
	rc.App.Get("/api/v1/dashboard/overview", rc.DashboardHandler.Overview)

	// Users
	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

	// Journal Entries
	rc.App.Post("/api/v1/journal_entries", rc.JournalEntryHandler.Create)
	rc.App.Get("/api/v1/journal_entries", rc.JournalEntryHandler.FindAll)
	rc.App.Get("/api/v1/journal_entries/:id", rc.JournalEntryHandler.FindById)
	rc.App.Put("/api/v1/journal_entries/:id", rc.JournalEntryHandler.Update)
	rc.App.Delete("/api/v1/journal_entries/:id", rc.JournalEntryHandler.Delete)
	rc.App.Put("/api/v1/journal_entries/:id/post", rc.JournalEntryHandler.Post)
	rc.App.Put("/api/v1/journal_entries/:id/void", rc.JournalEntryHandler.Void)

	// Customers
	rc.App.Post("/api/v1/customers", rc.CustomerHandler.Create)
	rc.App.Get("/api/v1/customers", rc.CustomerHandler.FindAll)
	rc.App.Get("/api/v1/customers/:id", rc.CustomerHandler.FindById)
	rc.App.Put("/api/v1/customers/:id", rc.CustomerHandler.Update)
	rc.App.Delete("/api/v1/customers/:id", rc.CustomerHandler.Delete)
	rc.App.Get("/api/v1/dropdown/customers", rc.CustomerHandler.SelectDropdownList)

	// Vendors
	rc.App.Post("/api/v1/vendors", rc.VendorHandler.Create)
	rc.App.Get("/api/v1/vendors", rc.VendorHandler.FindAll)
	rc.App.Get("/api/v1/vendors/:id", rc.VendorHandler.FindById)
	rc.App.Put("/api/v1/vendors/:id", rc.VendorHandler.Update)
	rc.App.Delete("/api/v1/vendors/:id", rc.VendorHandler.Delete)
	rc.App.Get("/api/v1/dropdown/vendors", rc.VendorHandler.SelectDropdownList)

	// Reports
	rc.App.Get("/api/v1/reports/trial-balance", rc.ReportHandler.TrialBalance)
	rc.App.Get("/api/v1/reports/profit-loss", rc.ReportHandler.ProfitLoss)
	rc.App.Get("/api/v1/reports/balance-sheet", rc.ReportHandler.BalanceSheet)
	rc.App.Get("/api/v1/reports/cash-flow", rc.ReportHandler.CashFlow)
	rc.App.Get("/api/v1/reports/equity-statement", rc.ReportHandler.EquityStatement)
}
