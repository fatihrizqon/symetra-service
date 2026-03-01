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

type COAGroupHandler struct {
	ICOAGroupService service.ICOAGroupService
}

func NewCOAGroupHandler(serv service.ICOAGroupService) *COAGroupHandler {
	return &COAGroupHandler{ICOAGroupService: serv}
}

// Create a New COAGroup
// @Summary Create chart of account group
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param request body request.COAGroupCreateRequest true "COAGroup Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa_groups [post]
func (h *COAGroupHandler) Create(ctx *fiber.Ctx) error {
	req := request.COAGroupCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAGroupService.Create(req)
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

// Find All COAGroups
// @Summary Get all chart of account groups
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort column (code, name, normal_balance, status, created_at, updated_at)"
// @Param order query string false "Sort direction (asc, desc)"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/coa_groups [get]
func (h *COAGroupHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COAGroup{}.SearchableFields())

	entities, totalCount, err := h.ICOAGroupService.FindAll(qp)
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
			Data:    []response.COAGroupResponse{},
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

// Find COAGroup by Id
// @Summary Get chart of account group by ID
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COAGroup ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COAGroup not found"
// @Router /api/v1/coa_groups/{id} [get]
func (h *COAGroupHandler) FindById(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAGroupService.FindById(parsedId)
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

// Update COAGroup by Id
// @Summary Update chart of account group
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COAGroup ID"
// @Param request body request.COAGroupUpdateRequest true "COAGroup Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COAGroup not found"
// @Router /api/v1/coa_groups/{id} [put]
func (h *COAGroupHandler) Update(ctx *fiber.Ctx) error {
	req := request.COAGroupUpdateRequest{}
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

	result, err := h.ICOAGroupService.Update(req)
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

// Delete COAGroup by Id
// @Summary Delete chart of account group
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COAGroup ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COAGroup not found"
// @Router /api/v1/coa_groups/{id} [delete]
func (h *COAGroupHandler) Delete(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAGroupService.Delete(parsedId)
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

// Select COAGroup Dropdown List
// @Summary Get COAGroup dropdown options
// @Tags Dropdowns
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/coa_groups [get]
func (h *COAGroupHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COAGroup{}.SearchableFields())

	entities, totalCount, err := h.ICOAGroupService.SelectDropdownList(qp)
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
