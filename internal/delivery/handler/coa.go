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
// @Description Store a new chart of account  record
// @Tags COAs
// @Accept json
// @Produce json
// @Param request body request.COACreateRequest true "COA Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa [post]
func (handler *COAHandler) Create(ctx *fiber.Ctx) error {
	req := request.COACreateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAService.Create(req)
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

// Find All COAs
// @Summary Get all chart of account s
// @Description Retrieve all chart of account  records with pagination
// @Tags COAs
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/coa [get]
func (handler *COAHandler) FindAll(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COA{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOAService.FindAll(page, pageSize, search, options)
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
			Data:    []response.COAResponse{},
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

// Find COA by Id
// @Summary Get chart of account  by ID
// @Description Retrieve a single chart of account  by its ID
// @Tags COAs
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COA not found"
// @Router /api/v1/coa/{id} [get]
func (handler *COAHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAService.FindById(parsedId)
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

// Update COA  by Id
// @Summary Update chart of account
// @Description Update chart of account  data by ID
// @Tags COAs
// @Accept json
// @Produce json
// @Param id path string true "COA  ID"
// @Param request body request.COAUpdateRequest true "COA  Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA  not found"
// @Router /api/v1/coa/{id} [put]
func (handler *COAHandler) Update(ctx *fiber.Ctx) error {
	req := request.COAUpdateRequest{}
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

	entity, err := handler.ICOAService.Update(req)
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

// Delete COA by Id
// @Summary Delete chart of account
// @Description Remove a chart of account  record by ID
// @Tags COAs
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COA not found"
// @Router /api/v1/coa/{id} [delete]
func (handler *COAHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.ICOAService.Delete(parsedId)
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

func (handler *COAHandler) setCOAFilters(ctx *fiber.Ctx) entity.COAFilters {
	filters := entity.COAFilters{}

	if status := ctx.Query("status"); status != "" {
		filters.Status = &status
	}

	return filters
}

// Select COA Dropdown List
// @Summary Get COA dropdown options
// @Description Retrieve a paginated list of chart of account s for dropdown selection. Supports optional search query.
// @Tags Dropdowns
// @Accept json
// @Produce json
// @Param search query string false "Search keyword for filtering s"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Failure 500 {object} response.JSON "Failed to retrieve records"
// @Router /api/v1/dropdown/coa [get]
func (handler *COAHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)

	search := ctx.Query("search")
	entity := entity.COA{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.ICOAService.SelectDropdownList(page, pageSize, search, options)
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
			Data:    []response.COAResponse{},
			Meta:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{
		Data: entities,
	})
}
