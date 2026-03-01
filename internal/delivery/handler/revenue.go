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

// RevenueHandler is a dedicated handler for revenue journal entries.
// It reuses IJournalEntryService and scopes all operations to type=revenue.
type RevenueHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewRevenueHandler(svc service.IJournalEntryService) *RevenueHandler {
	return &RevenueHandler{IJournalEntryService: svc}
}

// Create a New Revenue Entry
// @Summary Create revenue entry
// @Description Store a new revenue journal entry (type=revenue, prefix RV-).
// @Tags Revenue
// @Accept json
// @Produce json
// @Param request body request.JournalEntryCreateRequest true "Revenue Create Request"
// @Success 201 {object} response.JSON
// @Router /api/v1/revenues [post]
func (h *RevenueHandler) Create(ctx *fiber.Ctx) error {
	req := request.JournalEntryCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var createdBy uuid.UUID
	if claims, ok := ctx.Locals("auth").(*util.Claims); ok {
		createdBy = claims.UserID
	}

	result, err := h.IJournalEntryService.Create(req, createdBy, entity.JournalTypeRevenue)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "A new revenue entry has been stored.",
		Data:    result,
	})
}

// Find All Revenue Entries
// @Summary Get all revenue entries
// @Tags Revenue
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues [get]
func (h *RevenueHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.JournalEntry{}.SearchableFields())
	// Scope to revenue type only — override any client-supplied type filter
	qp.Filters["type"] = []string{string(entity.JournalTypeRevenue)}

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

// Find Revenue Entry by Id
// @Summary Get revenue entry by ID
// @Tags Revenue
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues/{id} [get]
func (h *RevenueHandler) FindById(ctx *fiber.Ctx) error {
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

// Update Revenue Entry by Id
// @Summary Update revenue entry
// @Tags Revenue
// @Accept json
// @Param id path string true "Entry ID"
// @Param request body request.JournalEntryUpdateRequest true "Revenue Update Request"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues/{id} [put]
func (h *RevenueHandler) Update(ctx *fiber.Ctx) error {
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

// Delete Revenue Entry by Id
// @Summary Delete revenue entry
// @Tags Revenue
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues/{id} [delete]
func (h *RevenueHandler) Delete(ctx *fiber.Ctx) error {
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

// Post Revenue Entry
// @Summary Post revenue entry
// @Tags Revenue
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues/{id}/post [put]
func (h *RevenueHandler) Post(ctx *fiber.Ctx) error {
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
		Message: "Revenue entry has been posted.",
	})
}

// Void Revenue Entry
// @Summary Void revenue entry
// @Tags Revenue
// @Param id path string true "Entry ID"
// @Success 200 {object} response.JSON
// @Router /api/v1/revenues/{id}/void [put]
func (h *RevenueHandler) Void(ctx *fiber.Ctx) error {
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
		Message: "Revenue entry has been voided.",
	})
}
