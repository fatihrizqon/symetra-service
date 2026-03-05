package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler

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
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)
	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)
	rc.App.Get("/api/v1/dashboard/overview", rc.DashboardHandler.Overview)

	// Users
	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

	// COA
	rc.App.Post("/api/v1/coa_groups", rc.COAGroupHandler.Create)
	rc.App.Get("/api/v1/coa_groups", rc.COAGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_groups/:id", rc.COAGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_groups/:id", rc.COAGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_groups/:id", rc.COAGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_groups", rc.COAGroupHandler.SelectDropdownList)

	rc.App.Post("/api/v1/coa_subgroups", rc.COASubGroupHandler.Create)
	rc.App.Get("/api/v1/coa_subgroups", rc.COASubGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa_subgroups", rc.COASubGroupHandler.SelectDropdownList)

	rc.App.Post("/api/v1/coa", rc.COAHandler.Create)
	rc.App.Get("/api/v1/coa", rc.COAHandler.FindAll)
	rc.App.Get("/api/v1/coa/:id", rc.COAHandler.FindById)
	rc.App.Put("/api/v1/coa/:id", rc.COAHandler.Update)
	rc.App.Delete("/api/v1/coa/:id", rc.COAHandler.Delete)
	rc.App.Get("/api/v1/dropdown/coa", rc.COAHandler.SelectDropdownList)

	// Journal
	rc.App.Post("/api/v1/journal_entries", rc.JournalEntryHandler.Create)
	rc.App.Get("/api/v1/journal_entries", rc.JournalEntryHandler.FindAll)
	rc.App.Get("/api/v1/journal_entries/:id", rc.JournalEntryHandler.FindById)
	rc.App.Put("/api/v1/journal_entries/:id", rc.JournalEntryHandler.Update)
	rc.App.Delete("/api/v1/journal_entries/:id", rc.JournalEntryHandler.Delete)
	rc.App.Put("/api/v1/journal_entries/:id/post", rc.JournalEntryHandler.Post)
	rc.App.Put("/api/v1/journal_entries/:id/void", rc.JournalEntryHandler.Void)

	rc.App.Post("/api/v1/revenues", rc.RevenueHandler.Create)
	rc.App.Get("/api/v1/revenues", rc.RevenueHandler.FindAll)
	rc.App.Get("/api/v1/revenues/:id", rc.RevenueHandler.FindById)
	rc.App.Put("/api/v1/revenues/:id", rc.RevenueHandler.Update)
	rc.App.Delete("/api/v1/revenues/:id", rc.RevenueHandler.Delete)
	rc.App.Put("/api/v1/revenues/:id/post", rc.RevenueHandler.Post)
	rc.App.Put("/api/v1/revenues/:id/void", rc.RevenueHandler.Void)

	rc.App.Post("/api/v1/expenses", rc.ExpenseHandler.Create)
	rc.App.Get("/api/v1/expenses", rc.ExpenseHandler.FindAll)
	rc.App.Get("/api/v1/expenses/:id", rc.ExpenseHandler.FindById)
	rc.App.Put("/api/v1/expenses/:id", rc.ExpenseHandler.Update)
	rc.App.Delete("/api/v1/expenses/:id", rc.ExpenseHandler.Delete)
	rc.App.Put("/api/v1/expenses/:id/post", rc.ExpenseHandler.Post)
	rc.App.Put("/api/v1/expenses/:id/void", rc.ExpenseHandler.Void)

	// Reports
	rc.App.Get("/api/v1/reports/trial-balance", rc.ReportHandler.TrialBalance)
	rc.App.Get("/api/v1/reports/profit-loss", rc.ReportHandler.ProfitLoss)
	rc.App.Get("/api/v1/reports/balance-sheet", rc.ReportHandler.BalanceSheet)
	rc.App.Get("/api/v1/reports/cash-flow", rc.ReportHandler.CashFlow)
	rc.App.Get("/api/v1/reports/equity-statement", rc.ReportHandler.EquityStatement)
}
