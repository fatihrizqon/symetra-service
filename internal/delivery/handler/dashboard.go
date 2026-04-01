package handler

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	IDashboardService service.IDashboardService
}

func NewDashboardHandler(service service.IDashboardService) *DashboardHandler {
	return &DashboardHandler{
		IDashboardService: service,
	}
}

func (h *DashboardHandler) Overview(ctx *fiber.Ctx) error {

	module := ctx.Query("module")
	period := ctx.Query("period", "30d")

	if module == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"module is required",
		)
	}

	result, err := h.IDashboardService.Overview(
		ctx.Context(),
		module,
		period,
	)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved all records.",
		Data:    result,
		Meta:    nil,
	})
}
