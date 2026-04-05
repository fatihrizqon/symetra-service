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
	svc service.ICustomerService
}

func NewCustomerHandler(svc service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) Create(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	var req request.CustomerCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	result, err := h.svc.Create(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{Status: fiber.StatusCreated, Message: "A new record has been stored.", Data: result})
}

func (h *CustomerHandler) FindAll(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.Customer{}.SearchableFields())
	items, totalCount, err := h.svc.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "No records found.", Data: []response.CustomerResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved all records.", Data: items, Meta: &meta})
}

func (h *CustomerHandler) FindById(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	result, err := h.svc.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Successfully retrieved selected record.", Data: result})
}

func (h *CustomerHandler) Update(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	var req request.CustomerUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	req.Id = id
	result, err := h.svc.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Selected record has been updated.", Data: result})
}

func (h *CustomerHandler) Delete(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid id"})
	}
	if err := h.svc.Delete(companyID, id); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{Status: fiber.StatusNotFound, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Selected record has been deleted."})
}

func (h *CustomerHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.Customer{}.SearchableFields())
	items, totalCount, err := h.svc.SelectDropdownList(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}
	if totalCount == 0 {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: []response.SelectDropdownListResponse{}})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}
