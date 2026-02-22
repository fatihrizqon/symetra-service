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

type VendorHandler struct {
	IVendorService service.IVendorService
}

func NewVendorHandler(svc service.IVendorService) *VendorHandler {
	return &VendorHandler{IVendorService: svc}
}

// Create a New Vendor
// @Summary Create vendor
// @Description Store a new vendor record
// @Tags Vendors
// @Accept json
// @Produce json
// @Param request body request.VendorCreateRequest true "Vendor Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/vendors [post]
func (h *VendorHandler) Create(ctx *fiber.Ctx) error {
	req := request.VendorCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IVendorService.Create(req)
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

// Find All Vendors
// @Summary Get all vendors
// @Description Retrieve all vendor records with pagination
// @Tags Vendors
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/vendors [get]
func (h *VendorHandler) FindAll(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)
	search := ctx.Query("search")
	options := util.SearchOptions{Fields: entity.Vendor{}.SearchableFields()}

	filters := entity.VendorFilters{}
	if status := ctx.Query("status"); status != "" {
		filters.Status = &status
	}

	items, totalCount, err := h.IVendorService.FindAll(page, pageSize, search, options, filters)
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
			Data:    []response.VendorResponse{},
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, search, page, pageSize, totalCount, nil)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  200,
		Message: "Successfully retrieved all records.",
		Data:    items,
		Meta:    &meta,
	})
}

// Find Vendor by Id
// @Summary Get vendor by ID
// @Description Retrieve a single vendor by its ID
// @Tags Vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [get]
func (h *VendorHandler) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IVendorService.FindById(parsedId)
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

// Update Vendor by Id
// @Summary Update vendor
// @Description Update vendor data by ID
// @Tags Vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Param request body request.VendorUpdateRequest true "Vendor Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [put]
func (h *VendorHandler) Update(ctx *fiber.Ctx) error {
	req := request.VendorUpdateRequest{}
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

	result, err := h.IVendorService.Update(req)
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

// Delete Vendor by Id
// @Summary Delete vendor
// @Description Remove a vendor record by ID
// @Tags Vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [delete]
func (h *VendorHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	_, err = h.IVendorService.Delete(parsedId)
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

// Select Vendor Dropdown List
// @Summary Get vendor dropdown options
// @Description Retrieve vendors for dropdown selection
// @Tags Dropdowns
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/vendors [get]
func (h *VendorHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	page, pageSize, _ := util.ParsePaginationParams(ctx)
	search := ctx.Query("search")
	options := util.SearchOptions{Fields: entity.Vendor{}.SearchableFields()}

	items, totalCount, err := h.IVendorService.SelectDropdownList(page, pageSize, search, options)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  500,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (page-1)*pageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{
			Data: []response.SelectDropdownListResponse{},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}
