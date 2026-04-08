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

// FIX [ARSITEKTUR-03]: Hapus fungsi lokal getCallerID — sudah ada util.GetCallerID yang identik.
// Duplikasi logika berisiko divergence jika satu diupdate dan lainnya tidak.
// Semua pemanggilan getCallerID di bawah ini sudah diganti ke util.GetCallerID.

// ── Fiscal Year ───────────────────────────────────────────────────────────────

// @Summary Create a new fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years [post]
func (h *FiscalHandler) CreateFiscalYear(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	var req request.FiscalYearCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	result, err := h.IFiscalService.CreateFiscalYear(companyID, req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{Status: fiber.StatusCreated, Message: "Fiscal year created successfully.", Data: result})
}

// @Summary Get all fiscal years
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years [get]
func (h *FiscalHandler) FindAllFiscalYears(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.FiscalYear{}.SearchableFields())
	entities, totalCount, err := h.IFiscalService.FindAllFiscalYears(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "No records found.", Data: []response.FiscalYearResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved fiscal years.", Data: entities, Meta: &meta})
}

// @Summary Get fiscal year by ID (includes periods)
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id} [get]
func (h *FiscalHandler) FindFiscalYearById(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.FindFiscalYearById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved fiscal year.", Data: result})
}

// @Summary Rename a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id} [put]
func (h *FiscalHandler) UpdateFiscalYear(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	var req request.FiscalYearUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	req.Id = id
	result, err := h.IFiscalService.UpdateFiscalYear(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal year updated.", Data: result})
}

// @Summary Activate a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id}/activate [put]
func (h *FiscalHandler) ActivateFiscalYear(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.ActivateFiscalYear(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal year is now active.", Data: result})
}

// @Summary Delete a draft fiscal year
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id} [delete]
func (h *FiscalHandler) DeleteFiscalYear(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	if err := h.IFiscalService.DeleteFiscalYear(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal year deleted."})
}

// ── Fiscal Periods ────────────────────────────────────────────────────────────

// @Summary Get all fiscal periods
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods [get]
func (h *FiscalHandler) FindAllPeriods(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.FiscalPeriod{}.SearchableFields())
	entities, totalCount, err := h.IFiscalService.FindAllPeriods(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "No records found.", Data: []response.FiscalPeriodResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved fiscal periods.", Data: entities, Meta: &meta})
}

// @Summary Get a fiscal period by ID
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods/{id} [get]
func (h *FiscalHandler) FindPeriodById(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.FindPeriodById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved fiscal period.", Data: result})
}

// @Summary Close an open fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods/{id}/close [put]
func (h *FiscalHandler) ClosePeriod(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.ClosePeriod(companyID, id, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal period closed successfully.", Data: result})
}

// @Summary Reopen a closed fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods/{id}/reopen [put]
func (h *FiscalHandler) ReopenPeriod(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	var req request.FiscalPeriodReopenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	req.Id = id
	result, err := h.IFiscalService.ReopenPeriod(companyID, req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal period reopened successfully.", Data: result})
}

// @Summary Permanently lock a closed fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods/{id}/lock [put]
func (h *FiscalHandler) LockPeriod(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.LockPeriod(companyID, id, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal period permanently locked.", Data: result})
}

// @Summary Get audit log for a fiscal period
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/periods/{id}/logs [get]
func (h *FiscalHandler) FindPeriodLogs(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	logs, err := h.IFiscalService.FindPeriodLogs(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved period logs.", Data: logs})
}

// ── Year-End Closing — NEW ────────────────────────────────────────────────────

// @Summary Check if a fiscal year is ready to close
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id}/readiness [get]
func (h *FiscalHandler) CheckReadiness(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.ReadyToClose(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Readiness check complete.", Data: result})
}

// @Summary Close a fiscal year (simple or formal mode)
// @Tags Fiscal
// @Security BearerAuth
// @Accept json
// @Param request body request.FiscalYearCloseRequest true "mode: simple|formal"
// @Router /api/v1/fiscal/years/{id}/close [put]
func (h *FiscalHandler) CloseFiscalYear(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	var req request.FiscalYearCloseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	result, err := h.IFiscalService.CloseFiscalYear(companyID, id, req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Fiscal year closed successfully.", Data: result})
}

// ── Opening Balance — NEW ─────────────────────────────────────────────────────

// @Summary Generate opening balance from previous period's ledger
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id}/opening-balance [post]
func (h *FiscalHandler) GenerateOpeningBalance(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.GenerateOpeningBalance(companyID, id, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Opening balance generated successfully.", Data: result})
}

// @Summary Delete opening balance (to allow regeneration)
// @Tags Fiscal
// @Security BearerAuth
// @Router /api/v1/fiscal/years/{id}/opening-balance [delete]
func (h *FiscalHandler) DeleteOpeningBalance(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IFiscalService.DeleteOpeningBalance(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Opening balance deleted. You may now regenerate it.", Data: result})
}
