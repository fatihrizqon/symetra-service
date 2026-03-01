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

// Create a New Journal Entry
// @Summary Create journal entry
// @Description Store a new journal entry. Total debit must equal total credit.
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param request body request.JournalEntryCreateRequest true "Journal Entry Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request or unbalanced entry"
// @Router /api/v1/journal_entries [post]
func (h *JournalEntryHandler) Create(ctx *fiber.Ctx) error {
	req := request.JournalEntryCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	// Extract authenticated user ID from context (set by auth middleware)
	var createdBy uuid.UUID
	if claims, ok := ctx.Locals("auth").(*util.Claims); ok {
		createdBy = claims.UserID
	}

	result, err := h.IJournalEntryService.Create(req, createdBy)
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

// Find All Journal Entries
// @Summary Get all journal entries
// @Description Retrieve all journal entries with pagination and optional status filter
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (draft, posted, void)"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/journal_entries [get]
func (h *JournalEntryHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.User{}.SearchableFields())

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
			Data:    []response.UserResponse{},
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
// @Description Retrieve a single journal entry by its ID including all lines
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Journal entry not found"
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
// @Description Update a draft journal entry. Posted entries cannot be updated.
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Param request body request.JournalEntryUpdateRequest true "Journal Entry Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 400 {object} response.JSON "Bad request or business rule violation"
// @Router /api/v1/journal_entries/{id} [put]
func (h *JournalEntryHandler) Update(ctx *fiber.Ctx) error {
	req := request.JournalEntryUpdateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
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
// @Description Delete a draft journal entry. Posted or void entries cannot be deleted.
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 400 {object} response.JSON "Business rule violation"
// @Failure 404 {object} response.JSON "Journal entry not found"
// @Router /api/v1/journal_entries/{id} [delete]
func (h *JournalEntryHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
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
// @Description Transition a draft journal entry to posted status. Entry must be balanced.
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Journal entry has been posted."
// @Failure 400 {object} response.JSON "Business rule violation"
// @Router /api/v1/journal_entries/{id}/post [patch]
func (h *JournalEntryHandler) Post(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
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
// @Description Transition a posted journal entry to void status.
// @Tags JournalEntries
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Journal entry has been voided."
// @Failure 400 {object} response.JSON "Business rule violation"
// @Router /api/v1/journal_entries/{id}/void [patch]
func (h *JournalEntryHandler) Void(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
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
