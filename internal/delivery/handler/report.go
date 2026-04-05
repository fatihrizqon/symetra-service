package handler

import (
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	IReportService service.IReportService
}

func NewReportHandler(svc service.IReportService) *ReportHandler {
	return &ReportHandler{IReportService: svc}
}

func parseDateRange(ctx *fiber.Ctx) (time.Time, time.Time, error) {
	now := time.Now()
	defaultStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	defaultEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	startStr := ctx.Query("start_date", defaultStart.Format("2006-01-02"))
	endStr := ctx.Query("end_date", defaultEnd.Format("2006-01-02"))

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
	}
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	return start, end, nil
}

// @Summary Trial Balance
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date   query string false "End date (YYYY-MM-DD)"
// @Router /api/v1/reports/trial-balance [get]
func (h *ReportHandler) TrialBalance(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	result, err := h.IReportService.TrialBalance(companyId, start, end)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Trial balance generated.", Data: result})
}

// @Summary Profit & Loss Statement
// @Tags Reports
// @Security BearerAuth
// @Router /api/v1/reports/profit-loss [get]
func (h *ReportHandler) ProfitLoss(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	result, err := h.IReportService.ProfitLoss(companyId, start, end)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Profit & loss statement generated.", Data: result})
}

// @Summary Balance Sheet
// @Tags Reports
// @Security BearerAuth
// @Param as_of_date query string false "As of date (YYYY-MM-DD)"
// @Router /api/v1/reports/balance-sheet [get]
func (h *ReportHandler) BalanceSheet(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	now := time.Now()
	asOfStr := ctx.Query("as_of_date", now.Format("2006-01-02"))
	asOf, err := time.Parse("2006-01-02", asOfStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid as_of_date format, use YYYY-MM-DD")
	}
	asOf = time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 23, 59, 59, 0, time.UTC)
	result, err := h.IReportService.BalanceSheet(companyId, asOf)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Balance sheet generated.", Data: result})
}

// @Summary Cash Flow Statement
// @Tags Reports
// @Security BearerAuth
// @Router /api/v1/reports/cash-flow [get]
func (h *ReportHandler) CashFlow(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	result, err := h.IReportService.CashFlow(companyId, start, end)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Cash flow statement generated.", Data: result})
}

// @Summary Equity Statement
// @Tags Reports
// @Security BearerAuth
// @Router /api/v1/reports/equity-statement [get]
func (h *ReportHandler) EquityStatement(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	result, err := h.IReportService.EquityStatement(companyId, start, end)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Equity statement generated.", Data: result})
}

// @Summary General Ledger (Buku Besar)
// @Tags Reports
// @Security BearerAuth
// @Param coa_id query string false "Filter by specific COA account UUID (optional)"
// @Router /api/v1/reports/general-ledger [get]
func (h *ReportHandler) GeneralLedger(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	coaID := ctx.Query("coa_id", "")
	result, err := h.IReportService.GeneralLedger(companyId, start, end, coaID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "General ledger generated.", Data: result})
}

// @Summary Journal Book (Jurnal Umum)
// @Tags Reports
// @Security BearerAuth
// @Router /api/v1/reports/journal-book [get]
func (h *ReportHandler) JournalBook(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	start, end, err := parseDateRange(ctx)
	if err != nil {
		return err
	}
	result, err := h.IReportService.JournalBook(companyId, start, end)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return ctx.JSON(response.JSON{Status: fiber.StatusOK, Message: "Journal book generated.", Data: result})
}
