package handler

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FiscalHandler struct {
	IFiscalService service.IFiscalService
}

func NewFiscalHandler(svc service.IFiscalService) *FiscalHandler {
	return &FiscalHandler{IFiscalService: svc}
}

// getCallerID extracts the authenticated user UUID from Fiber locals (set by auth middleware).
func getCallerID(ctx *fiber.Ctx) (uuid.UUID, error) {
	claims, ok := ctx.Locals("auth").(*util.Claims)
	if !ok || claims == nil {
		return uuid.Nil, fiber.ErrUnauthorized
	}
	return claims.UserID, nil
}

// ── Fiscal Year ───────────────────────────────────────────────────────────────

// CreateFiscalYear godoc
// @Summary Create a new fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.FiscalYearCreateRequest true "Fiscal Year Create Request"
// @Success 201 {object} response.JSON
// @Router /api/v1/fiscal/years [post]
func (h *FiscalHandler) CreateFiscalYear(ctx *fiber.Ctx) error {
	callerID, err := getCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	var req request.FiscalYearCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.CreateFiscalYear(req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "Fiscal year created successfully.",
		Data:    result,
	})
}

// FindAllFiscalYears godoc
// @Summary Get all fiscal years
// @Tags Fiscal
// @Security BearerAuth
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (draft, active, closed)"
// @Router /api/v1/fiscal/years [get]
func (h *FiscalHandler) FindAllFiscalYears(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.FiscalYear{}.SearchableFields())

	entities, totalCount, err := h.IFiscalService.FindAllFiscalYears(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve fiscal years",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  200,
			Message: "No records found.",
			Data:    []response.FiscalYearResponse{},
			Meta:    nil,
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved fiscal years.",
		Data:    entities,
		Meta:    &meta,
	})
}

// FindFiscalYearById godoc
// @Summary Get fiscal year by ID (includes periods)
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Fiscal Year ID"
// @Router /api/v1/fiscal/years/{id} [get]
func (h *FiscalHandler) FindFiscalYearById(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.FindFiscalYearById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved fiscal year.",
		Data:    result,
	})
}

// UpdateFiscalYear godoc
// @Summary Rename a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Fiscal Year ID"
// @Param request body request.FiscalYearUpdateRequest true "Update Request"
// @Router /api/v1/fiscal/years/{id} [put]
func (h *FiscalHandler) UpdateFiscalYear(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var req request.FiscalYearUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id

	result, err := h.IFiscalService.UpdateFiscalYear(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal year updated.",
		Data:    result,
	})
}

// ActivateFiscalYear godoc
// @Summary Activate a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Fiscal Year ID"
// @Router /api/v1/fiscal/years/{id}/activate [put]
func (h *FiscalHandler) ActivateFiscalYear(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.ActivateFiscalYear(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal year is now active.",
		Data:    result,
	})
}

// DeleteFiscalYear godoc
// @Summary Delete a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Fiscal Year ID"
// @Router /api/v1/fiscal/years/{id} [delete]
func (h *FiscalHandler) DeleteFiscalYear(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IFiscalService.DeleteFiscalYear(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal year deleted.",
	})
}

// ── Fiscal Periods ────────────────────────────────────────────────────────────

// FindAllPeriods godoc
// @Summary Get all fiscal periods
// @Tags Fiscal
// @Security BearerAuth
// @Param fiscal_year_id query string false "Filter by fiscal year ID"
// @Param status query string false "Filter by status (open, closed, locked)"
// @Router /api/v1/fiscal/periods [get]
func (h *FiscalHandler) FindAllPeriods(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.FiscalPeriod{}.SearchableFields())

	entities, totalCount, err := h.IFiscalService.FindAllPeriods(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve periods",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  200,
			Message: "No records found.",
			Data:    []response.FiscalPeriodResponse{},
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved fiscal periods.",
		Data:    entities,
		Meta:    &meta,
	})
}

// FindPeriodById godoc
// @Summary Get a fiscal period by ID
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Period ID"
// @Router /api/v1/fiscal/periods/{id} [get]
func (h *FiscalHandler) FindPeriodById(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.FindPeriodById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved fiscal period.",
		Data:    result,
	})
}

// ClosePeriod godoc
// @Summary Close an open fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Period ID"
// @Router /api/v1/fiscal/periods/{id}/close [put]
func (h *FiscalHandler) ClosePeriod(ctx *fiber.Ctx) error {
	callerID, err := getCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.ClosePeriod(id, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal period closed successfully.",
		Data:    result,
	})
}

// ReopenPeriod godoc
// @Summary Reopen a closed fiscal period (requires reason)
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Period ID"
// @Param request body request.FiscalPeriodReopenRequest true "Reopen Request"
// @Router /api/v1/fiscal/periods/{id}/reopen [put]
func (h *FiscalHandler) ReopenPeriod(ctx *fiber.Ctx) error {
	callerID, err := getCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var req request.FiscalPeriodReopenRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id

	result, err := h.IFiscalService.ReopenPeriod(req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal period reopened successfully.",
		Data:    result,
	})
}

// LockPeriod godoc
// @Summary Permanently lock a closed fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Period ID"
// @Router /api/v1/fiscal/periods/{id}/lock [put]
func (h *FiscalHandler) LockPeriod(ctx *fiber.Ctx) error {
	callerID, err := getCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IFiscalService.LockPeriod(id, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Fiscal period permanently locked.",
		Data:    result,
	})
}

// FindPeriodLogs godoc
// @Summary Get audit log for a fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Param id path string true "Period ID"
// @Router /api/v1/fiscal/periods/{id}/logs [get]
func (h *FiscalHandler) FindPeriodLogs(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	logs, err := h.IFiscalService.FindPeriodLogs(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved period logs.",
		Data:    logs,
	})
}
