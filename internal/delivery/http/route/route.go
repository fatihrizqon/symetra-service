package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App                *fiber.App
	UserHandler        *handler.UserHandler
	AuthHandler        *handler.AuthHandler
	AuthMiddleware     fiber.Handler
	DashboardHandler   *handler.DashboardHandler
	COAGroupHandler    *handler.COAGroupHandler
	COASubGroupHandler *handler.COASubGroupHandler
	// COAHandler         *handler.COAHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)

	rc.App.Get("/swagger/*", swagger.HandlerDefault)

	rc.App.Post("/api/v1/coa_groups", rc.COAGroupHandler.Create)
	rc.App.Get("/api/v1/coa_groups", rc.COAGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_groups/:id", rc.COAGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_groups/:id", rc.COAGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_groups/:id", rc.COAGroupHandler.Delete)

	rc.App.Post("/api/v1/coa_subgroups", rc.COASubGroupHandler.Create)
	rc.App.Get("/api/v1/coa_subgroups", rc.COASubGroupHandler.FindAll)
	rc.App.Get("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.FindById)
	rc.App.Put("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Update)
	rc.App.Delete("/api/v1/coa_subgroups/:id", rc.COASubGroupHandler.Delete)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)
	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)

	rc.App.Get("/api/v1/dashboard/overview", rc.DashboardHandler.Overview)

	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

}
