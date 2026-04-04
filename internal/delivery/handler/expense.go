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

type ExpenseHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewExpenseHandler(svc service.IJournalEntryService) *ExpenseHandler {
	return &ExpenseHandler{IJournalEntryService: svc}
}

func (h *ExpenseHandler) Create(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	callerID, _ := util.GetCallerID(ctx)
	var req request.JournalEntryCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.IJournalEntryService.Create(companyID, req, callerID, entity.JournalTypeExpense)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{Status: fiber.StatusCreated, Message: "Expense entry created.", Data: result})
}

func (h *ExpenseHandler) FindAll(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	qp := util.ParseQueryParams(ctx, entity.JournalEntry{}.SearchableFields())
	qp.Filters["type"] = []string{string(entity.JournalTypeExpense)}
	entities, totalCount, err := h.IJournalEntryService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "No records found.", Data: []response.JournalEntryResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved expense records.", Data: entities, Meta: &meta})
}

func (h *ExpenseHandler) FindById(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.IJournalEntryService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Record found.", Data: result})
}

func (h *ExpenseHandler) Update(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	var req request.JournalEntryUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id
	result, err := h.IJournalEntryService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Expense updated.", Data: result})
}

func (h *ExpenseHandler) Delete(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IJournalEntryService.Delete(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Expense deleted."})
}

func (h *ExpenseHandler) Post(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IJournalEntryService.Post(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Expense posted."})
}

func (h *ExpenseHandler) Void(ctx *fiber.Ctx) error {
	companyID, _ := util.GetCompanyID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IJournalEntryService.Void(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Expense voided."})
}
