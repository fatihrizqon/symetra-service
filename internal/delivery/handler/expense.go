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

// ExpenseHandler is a dedicated handler for expense journal entries.
// It reuses IJournalEntryService and scopes all operations to type=expense.
type ExpenseHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewExpenseHandler(svc service.IJournalEntryService) *ExpenseHandler {
	return &ExpenseHandler{IJournalEntryService: svc}
}

// Create a New Expense Entry
// @Summary Create expense entry
// @Description Store a new expense journal entry (type=expense, prefix EX-).
// @Tags Expenses
// @Accept json
// @Produce json
// @Param request body request.JournalEntryCreateRequest true "Expense Create Request"
// @Success 201 {object} response.JSON
// @Router /api/v1/expenses [post]
func (h *ExpenseHandler) Create(ctx *fiber.Ctx) error {
	req := request.JournalEntryCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var createdBy uuid.UUID
	if claims, ok := ctx.Locals("auth").(*util.Claims); ok {
		createdBy = claims.UserID
	}

	result, err := h.IJournalEntryService.Create(req, createdBy, entity.JournalTypeExpense)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "A new expense entry has been stored.",
		Data:    result,
	})
}

// Find All Expense Entries
// @Summary Get all expense entries
// @Tags Expenses
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses [get]
func (h *ExpenseHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.JournalEntry{}.SearchableFields())
	// Scope to expense type only
	qp.Filters["type"] = []string{string(entity.JournalTypeExpense)}

	entities, totalCount, err := h.IJournalEntryService.FindAll(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  200,
			Message: "No records found.",
			Data:    []response.JournalEntryResponse{},
			Meta:    nil,
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved all records.",
		Data:    entities,
		Meta:    &meta,
	})
}

// Find Expense Entry by Id
// @Summary Get expense entry by ID
// @Tags Expenses
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses/{id} [get]
func (h *ExpenseHandler) FindById(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entry, err := h.IJournalEntryService.FindById(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved selected record.",
		Data:    entry,
	})
}

// Update Expense Entry by Id
// @Summary Update expense entry
// @Tags Expenses
// @Accept json
// @Param id path string true "Entry ID"
// @Param request body request.JournalEntryUpdateRequest true "Expense Update Request"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses/{id} [put]
func (h *ExpenseHandler) Update(ctx *fiber.Ctx) error {
	req := request.JournalEntryUpdateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = parsedId

	result, err := h.IJournalEntryService.Update(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been updated.",
		Data:    result,
	})
}

// Delete Expense Entry by Id
// @Summary Delete expense entry
// @Tags Expenses
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses/{id} [delete]
func (h *ExpenseHandler) Delete(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IJournalEntryService.Delete(parsedId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been deleted.",
	})
}

// Post Expense Entry
// @Summary Post expense entry
// @Tags Expenses
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses/{id}/post [put]
func (h *ExpenseHandler) Post(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IJournalEntryService.Post(parsedId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Expense entry has been posted.",
	})
}

// Void Expense Entry
// @Summary Void expense entry
// @Tags Expenses
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/expenses/{id}/void [put]
func (h *ExpenseHandler) Void(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IJournalEntryService.Void(parsedId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Expense entry has been voided.",
	})
}
