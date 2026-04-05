package handler

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/fatihrizqon/symetra-service/internal/util"
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
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

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
		companyId,
		module,
		period,
	)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK,
		Message: "Successfully retrieved all records.",
		Data:    result,
		Meta:    nil,
	})
}
