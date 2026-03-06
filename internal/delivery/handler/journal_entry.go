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

type JournalEntryHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewJournalEntryHandler(svc service.IJournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{IJournalEntryService: svc}
}

// Create a New General Journal Entry
// @Summary Create general journal entry
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.JournalEntryCreateRequest true "Journal Entry Create Request"
// @Success 201 {object} response.JSON
// @Router /api/v1/journal_entries [post]
func (h *JournalEntryHandler) Create(ctx *fiber.Ctx) error {
	req := request.JournalEntryCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var createdBy uuid.UUID
	if claims, ok := ctx.Locals("auth").(*util.Claims); ok {
		createdBy = claims.UserID
	}

	result, err := h.IJournalEntryService.Create(req, createdBy, entity.JournalTypeGeneral)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "A new record has been stored.",
		Data:    result,
	})
}

// Find All General Journal Entries
// @Summary Get all general journal entries
// @Tags JournalEntries
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (draft, posted, void)"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries [get]
func (h *JournalEntryHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.JournalEntry{}.SearchableFields())
	// Scope to general type only
	qp.Filters["type"] = []string{string(entity.JournalTypeGeneral)}

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

// Find Journal Entry by Id
// @Summary Get journal entry by ID
// @Tags JournalEntries
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries/{id} [get]
func (h *JournalEntryHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
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

// Update Journal Entry by Id
// @Summary Update journal entry
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Param request body request.JournalEntryUpdateRequest true "Journal Entry Update Request"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries/{id} [put]
func (h *JournalEntryHandler) Update(ctx *fiber.Ctx) error {
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

// Delete Journal Entry by Id
// @Summary Delete journal entry
// @Tags JournalEntries
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries/{id} [delete]
func (h *JournalEntryHandler) Delete(ctx *fiber.Ctx) error {
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

// Post Journal Entry
// @Summary Post journal entry
// @Tags JournalEntries
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries/{id}/post [put]
func (h *JournalEntryHandler) Post(ctx *fiber.Ctx) error {
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
		Message: "Journal entry has been posted.",
	})
}

// Void Journal Entry
// @Summary Void journal entry
// @Tags JournalEntries
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/journal_entries/{id}/void [put]
func (h *JournalEntryHandler) Void(ctx *fiber.Ctx) error {
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
		Message: "Journal entry has been voided.",
	})
}
