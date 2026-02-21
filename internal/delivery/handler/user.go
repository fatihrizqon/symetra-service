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

type UserHandler struct {
	IUserService service.IUserService
}

func NewUserHandler(serv service.IUserService) *UserHandler {
	return &UserHandler{IUserService: serv}
}

// Create a New User
// @Summary Create user
// @Description Store a new user record
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.UserCreateRequest true "User Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/users [post]
func (handler *UserHandler) Create(ctx *fiber.Ctx) error {
	req := request.UserCreateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.IUserService.Create(req)
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

// Find All Users
// @Summary Get all users
// @Description Retrieve all user records with pagination
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/users [get]
func (handler *UserHandler) FindAll(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)
	userFilters := handler.setUserFilters(ctx)

	search := ctx.Query("search")
	entity := entity.User{}
	options := util.SearchOptions{
		Fields: entity.SearchableFields(),
	}

	entities, totalCount, err := handler.IUserService.FindAll(page, pageSize, search, options, userFilters)
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
			Data:    []response.UserResponse{},
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

// Find User by Id
// @Summary Get user by ID
// @Description Retrieve a single user by its ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "User not found"
// @Router /api/v1/users/{id} [get]
func (handler *UserHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.IUserService.FindById(parsedId)
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

// Update User by Id
// @Summary Update user
// @Description Update user data by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UserUpdateRequest true "User Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "User not found"
// @Router /api/v1/users/{id} [put]
func (handler *UserHandler) Update(ctx *fiber.Ctx) error {
	req := request.UserUpdateRequest{}
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

	entity, err := handler.IUserService.Update(req)
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

// Delete User by Id
// @Summary Delete user
// @Description Remove a user record by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "User not found"
// @Router /api/v1/users/{id} [delete]
func (handler *UserHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	entity, err := handler.IUserService.Delete(parsedId)
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

func (handler *UserHandler) setUserFilters(ctx *fiber.Ctx) entity.UserFilters {
	filters := entity.UserFilters{}

	if status := ctx.Query("status"); status != "" {
		filters.Status = &status
	}

	if verified := ctx.Query("verified"); verified != "" {
		filters.Verified = &verified
	}

	return filters
}
