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
// @Description Store a new chart of account group record
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param request body request.COAGroupCreateRequest true "COAGroup Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa_groups [post]
func (handler *COAGroupHandler) Create(ctx *fiber.Ctx) error {
	req := request.COAGroupCreateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAGroupService.Create(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "A new record has been stored.",
		Data:    entity,
	})
}

// Find All COAGroups
// @Summary Get all chart of account groups
// @Description Retrieve all chart of account group records with pagination
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/coa_groups [get]
func (handler *COAGroupHandler) FindAll(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COAGroup{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOAGroupService.FindAll(page, pageSize, search, options)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (page-1)*pageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  200,
			Message: "No records found.",
			Data:    []response.COAGroupResponse{},
			Meta:    nil,
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, search, page, pageSize, totalCount, nil)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved all records.",
		Data:    entities,
		Meta:    &meta,
	})
}

// Find COAGroup by Id
// @Summary Get chart of account group by ID
// @Description Retrieve a single chart of account group by its ID
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COAGroup ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COAGroup not found"
// @Router /api/v1/coa_groups/{id} [get]
func (handler *COAGroupHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAGroupService.FindById(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved selected record.",
		Data:    entity,
	})
}

// Update COA Group by Id
// @Summary Update chart of account group
// @Description Update chart of account group data by ID
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COA Group ID"
// @Param request body request.COAGroupUpdateRequest true "COA Group Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA Group not found"
// @Router /api/v1/coa_groups/{id} [put]
func (handler *COAGroupHandler) Update(ctx *fiber.Ctx) error {
	req := request.COAGroupUpdateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
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

	entity, err := handler.ICOAGroupService.Update(req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been updated.",
		Data:    entity,
	})
}

// Delete COAGroup by Id
// @Summary Delete chart of account group
// @Description Remove a chart of account group record by ID
// @Tags COAGroups
// @Accept json
// @Produce json
// @Param id path string true "COAGroup ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COAGroup not found"
// @Router /api/v1/coa_groups/{id} [delete]
func (handler *COAGroupHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAGroupService.Delete(parsedId)
	if err != nil {
		resp := response.JSON{
			Status:  404,
			Message: err.Error(),
		}
		if entity.Id == uuid.Nil {
			return ctx.Status(fiber.StatusNotFound).JSON(resp)
		}
		return ctx.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := response.JSON{
		Status:  200,
		Message: "Selected record has been deleted.",
		Data:    nil,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (handler *COAGroupHandler) setCOAGroupFilters(ctx *fiber.Ctx) entity.COAGroupFilters {
	filters := entity.COAGroupFilters{}

	if status := ctx.Query("status"); status != "" {
		filters.Status = &status
	}

	return filters
}

// Select COAGroup Dropdown List
// @Summary Get COAGroup dropdown options
// @Description Retrieve a paginated list of chart of account groups for dropdown selection. Supports optional search query.
// @Tags Dropdowns
// @Accept json
// @Produce json
// @Param search query string false "Search keyword for filtering groups"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Failure 500 {object} response.JSON "Failed to retrieve records"
// @Router /api/v1/dropdown/coa_groups [get]
func (handler *COAGroupHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COAGroup{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOAGroupService.SelectDropdownList(page, pageSize, search, options)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (page-1)*pageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  200,
			Message: "No records found.",
			Data:    []response.COAGroupResponse{},
			Meta:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{
		Data: entities,
	})
}
