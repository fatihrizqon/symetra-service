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

type COASubGroupHandler struct {
	ICOASubGroupService service.ICOASubGroupService
}

func NewCOASubGroupHandler(serv service.ICOASubGroupService) *COASubGroupHandler {
	return &COASubGroupHandler{ICOASubGroupService: serv}
}

// Create a New COASubGroup
// @Summary Create chart of account subgroup
// @Tags #ChartofAccounts SubGroups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.COASubGroupCreateRequest true "COASubGroup Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa_subgroups [post]
func (h *COASubGroupHandler) Create(ctx *fiber.Ctx) error {
	req := request.COASubGroupCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOASubGroupService.Create(req)
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

// Find All COASubGroups
// @Summary Get all chart of account subgroups
// @Tags #ChartofAccounts SubGroups
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
// @Router /api/v1/coa_subgroups [get]
func (h *COASubGroupHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COASubGroup{}.SearchableFields())

	entities, totalCount, err := h.ICOASubGroupService.FindAll(qp)
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
			Data:    []response.COASubGroupResponse{},
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

// Find COASubGroup by Id
// @Summary Get chart of account subgroup by ID
// @Tags #ChartofAccounts SubGroups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COASubGroup ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COASubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [get]
func (h *COASubGroupHandler) FindById(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOASubGroupService.FindById(parsedId)
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

// Update COASubGroup by Id
// @Summary Update chart of account subgroup
// @Tags #ChartofAccounts SubGroups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COASubGroup ID"
// @Param request body request.COASubGroupUpdateRequest true "COASubGroup Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COASubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [put]
func (h *COASubGroupHandler) Update(ctx *fiber.Ctx) error {
	req := request.COASubGroupUpdateRequest{}
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

	result, err := h.ICOASubGroupService.Update(req)
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

// Delete COASubGroup by Id
// @Summary Delete chart of account subgroup
// @Tags #ChartofAccounts SubGroups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COASubGroup ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COASubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [delete]
func (h *COASubGroupHandler) Delete(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOASubGroupService.Delete(parsedId)
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

// Select COASubGroup Dropdown List
// @Summary Get COASubGroup dropdown options
// @Tags Dropdowns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/coa_subgroups [get]
func (h *COASubGroupHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.COASubGroup{}.SearchableFields())

	entities, totalCount, err := h.ICOASubGroupService.SelectDropdownList(qp)
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
