package route

import (
	_ "github.com/fatihrizqon/symetra-service/docs"
	"github.com/fatihrizqon/symetra-service/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App              *fiber.App
	UserHandler      *handler.UserHandler
	AuthHandler      *handler.AuthHandler
	AuthMiddleware   fiber.Handler
	DashboardHandler *handler.DashboardHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)

	rc.App.Get("/api/v1/dashboard/overview", rc.DashboardHandler.Overview)

	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)
	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)
	// rc.App.Post("/api/v1/auth/me", rc.AuthHandler.Me)

	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/api/v1/users", rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.UserHandler.Delete)

}
