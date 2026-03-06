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

type CustomerHandler struct {
	ICustomerService service.ICustomerService
}

func NewCustomerHandler(svc service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{ICustomerService: svc}
}

// Create a New Customer
// @Summary Create customer
// @Description Store a new customer record
// @Tags Customers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CustomerCreateRequest true "Customer Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/customers [post]
func (h *CustomerHandler) Create(ctx *fiber.Ctx) error {
	req := request.CustomerCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICustomerService.Create(req)
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

// Find All Customers
// @Summary Get all customers
// @Description Retrieve all customer records with pagination
// @Tags Customers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort column (code, name, email, status, created_at, updated_at)"
// @Param order query string false "Sort direction (asc, desc)"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/customers [get]
func (h *CustomerHandler) FindAll(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.Customer{}.SearchableFields())

	items, totalCount, err := h.ICustomerService.FindAll(qp)
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
			Data:    []response.CustomerResponse{},
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved all records.",
		Data:    items,
		Meta:    &meta,
	})
}

// Find Customer by Id
// @Summary Get customer by ID
// @Description Retrieve a single customer by its ID
// @Tags Customers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Customer not found"
// @Router /api/v1/customers/{id} [get]
func (h *CustomerHandler) FindById(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICustomerService.FindById(parsedId)
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

// Update Customer by Id
// @Summary Update customer
// @Description Update customer data by ID
// @Tags Customers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body request.CustomerUpdateRequest true "Customer Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Customer not found"
// @Router /api/v1/customers/{id} [put]
func (h *CustomerHandler) Update(ctx *fiber.Ctx) error {
	req := request.CustomerUpdateRequest{}
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

	result, err := h.ICustomerService.Update(req)
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

// Delete Customer by Id
// @Summary Delete customer
// @Description Remove a customer record by ID
// @Tags Customers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Customer not found"
// @Router /api/v1/customers/{id} [delete]
func (h *CustomerHandler) Delete(ctx *fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	_, err = h.ICustomerService.Delete(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  404,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Selected record has been deleted.",
	})
}

// Select Customer Dropdown List
// @Summary Get customer dropdown options
// @Description Retrieve customers for dropdown selection
// @Tags Dropdowns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/customers [get]
func (h *CustomerHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.Customer{}.SearchableFields())

	items, totalCount, err := h.ICustomerService.SelectDropdownList(qp)
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

	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}
