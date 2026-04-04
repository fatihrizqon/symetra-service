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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.VendorCreateRequest true "Vendor Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/vendors [post]
func (h *VendorHandler) Create(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	req := request.VendorCreateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	result, err := h.IVendorService.Create(companyId, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{Status: fiber.StatusCreated, Message: "A new record has been stored.", Data: result})
}

// Find All Vendors
// @Summary Get all vendors
// @Description Retrieve all vendor records with pagination
// @Tags Vendors
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
// @Router /api/v1/vendors [get]
func (h *VendorHandler) FindAll(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.Vendor{}.SearchableFields())
	items, totalCount, err := h.IVendorService.FindAll(companyId, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "No records found.", Data: []response.VendorResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved all records.", Data: items, Meta: &meta})
}

// Find Vendor by Id
// @Summary Get vendor by ID
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [get]
func (h *VendorHandler) FindById(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.IVendorService.FindById(companyId, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved selected record.", Data: result})
}

// Update Vendor by Id
// @Summary Update vendor
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Param request body request.VendorUpdateRequest true "Vendor Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [put]
func (h *VendorHandler) Update(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	req := request.VendorUpdateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	req.Id = id
	result, err := h.IVendorService.Update(companyId, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Selected record has been updated.", Data: result})
}

// Delete Vendor by Id
// @Summary Delete vendor
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [delete]
func (h *VendorHandler) Delete(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	if _, err := h.IVendorService.Delete(companyId, id); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Selected record has been deleted."})
}

// Select Vendor Dropdown List
// @Summary Get vendor dropdown options
// @Tags Dropdowns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.SelectJSON "Successfully retrieved dropdown options"
// @Router /api/v1/dropdown/vendors [get]
func (h *VendorHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.Vendor{}.SearchableFields())
	items, totalCount, err := h.IVendorService.SelectDropdownList(companyId, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	if totalCount == 0 {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: []response.SelectDropdownListResponse{}})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}
