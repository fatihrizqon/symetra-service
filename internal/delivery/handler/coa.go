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

type COAHandler struct {
	ICOAService service.ICOAService
}

func NewCOAHandler(serv service.ICOAService) *COAHandler {
	return &COAHandler{ICOAService: serv}
}

// Create a New COA
// @Summary Create chart of account
// @Tags #ChartofAccounts COA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.COACreateRequest true "COA Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa [post]
func (h *COAHandler) Create(ctx *fiber.Ctx) error {
	req := request.COACreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAService.Create(req)
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

// Find All COAs
// @Summary Get all chart of accounts
// @Tags #ChartofAccounts COA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort column (code, name, status, created_at, updated_at)"
// @Param order query string false "Sort direction (asc, desc)"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/coa [get]
func (h *COAHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COA{}.SearchableFields())

	entities, totalCount, err := h.ICOAService.FindAll(qp)
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
			Data:    []response.COAResponse{},
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

// Find COA by Id
// @Summary Get chart of account by ID
// @Tags #ChartofAccounts COA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COA not found"
// @Router /api/v1/coa/{id} [get]
func (h *COAHandler) FindById(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAService.FindById(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved selected record.",
		Data:    result,
	})
}

// Update COA by Id
// @Summary Update chart of account
// @Tags #ChartofAccounts COA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Param request body request.COAUpdateRequest true "COA Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA not found"
// @Router /api/v1/coa/{id} [put]
func (h *COAHandler) Update(ctx *fiber.Ctx) error {
	req := request.COAUpdateRequest{}
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

	result, err := h.ICOAService.Update(req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been updated.",
		Data:    result,
	})
}

// Delete COA by Id
// @Summary Delete chart of account
// @Tags #ChartofAccounts COA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COA not found"
// @Router /api/v1/coa/{id} [delete]
func (h *COAHandler) Delete(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAService.Delete(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been deleted.",
		Data:    result,
	})
}

// Select COA Dropdown List
// @Summary Get COA dropdown options
// @Tags Dropdowns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/coa [get]
func (h *COAHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COA{}.SearchableFields())

	entities, totalCount, err := h.ICOAService.SelectDropdownList(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{
			Data: []response.SelectDropdownListResponse{},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: entities})
}
