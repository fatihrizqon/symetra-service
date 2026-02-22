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
// @Description Store a new chart of account subgroup record
// @Tags COASubGroups
// @Accept json
// @Produce json
// @Param request body request.COASubGroupCreateRequest true "COASubGroup Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa_subgroups [post]
func (handler *COASubGroupHandler) Create(ctx *fiber.Ctx) error {
	req := request.COASubGroupCreateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOASubGroupService.Create(req)
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

// Find All COASubGroups
// @Summary Get all chart of account subgroups
// @Description Retrieve all chart of account subgroup records with pagination
// @Tags COASubGroups
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/coa_subgroups [get]
func (handler *COASubGroupHandler) FindAll(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COASubGroup{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOASubGroupService.FindAll(page, pageSize, search, options)
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
			Data:    []response.COASubGroupResponse{},
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

// Find COASubGroup by Id
// @Summary Get chart of account subgroup by ID
// @Description Retrieve a single chart of account subgroup by its ID
// @Tags COASubGroups
// @Accept json
// @Produce json
// @Param id path string true "COASubGroup ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COASubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [get]
func (handler *COASubGroupHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOASubGroupService.FindById(parsedId)
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

// Update COA SubGroup by Id
// @Summary Update chart of account subgroup
// @Description Update chart of account subgroup data by ID
// @Tags COASubGroups
// @Accept json
// @Produce json
// @Param id path string true "COA SubGroup ID"
// @Param request body request.COASubGroupUpdateRequest true "COA SubGroup Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA SubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [put]
func (handler *COASubGroupHandler) Update(ctx *fiber.Ctx) error {
	req := request.COASubGroupUpdateRequest{}
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

	entity, err := handler.ICOASubGroupService.Update(req)
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

// Delete COASubGroup by Id
// @Summary Delete chart of account subgroup
// @Description Remove a chart of account subgroup record by ID
// @Tags COASubGroups
// @Accept json
// @Produce json
// @Param id path string true "COASubGroup ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COASubGroup not found"
// @Router /api/v1/coa_subgroups/{id} [delete]
func (handler *COASubGroupHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOASubGroupService.Delete(parsedId)
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

// Select COASubGroup Dropdown List
// @Summary Get COASubGroup dropdown options
// @Description Retrieve a paginated list of chart of account subgroups for dropdown selection. Supports optional search query.
// @Tags Dropdowns
// @Accept json
// @Produce json
// @Param search query string false "Search keyword for filtering subgroups"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Failure 500 {object} response.JSON "Failed to retrieve records"
// @Router /api/v1/dropdown/coa_subgroups/ [get]
func (handler *COASubGroupHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COASubGroup{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOASubGroupService.SelectDropdownList(page, pageSize, search, options)
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
