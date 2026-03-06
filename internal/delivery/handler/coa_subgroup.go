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

func (h *COASubGroupHandler) Create(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	var req request.COASubGroupCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.ICOASubGroupService.Create(companyID, req)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(201).JSON(response.JSON{Status: 201, Message: "A new record has been stored.", Data: result})
}

func (h *COASubGroupHandler) FindAll(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, entity.COASubGroup{}.SearchableFields())
	entities, totalCount, err := h.ICOASubGroupService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(500).JSON(response.JSON{Status: 500, Message: "Failed to retrieve records.", Errors: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "No records found.", Data: []response.COASubGroupResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Successfully retrieved records.", Data: entities, Meta: &meta})
}

func (h *COASubGroupHandler) FindById(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.ICOASubGroupService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(404).JSON(response.JSON{Status: 404, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record found.", Data: result})
}

func (h *COASubGroupHandler) Update(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	var req request.COASubGroupUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id
	result, err := h.ICOASubGroupService.Update(companyID, req)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record updated.", Data: result})
}

func (h *COASubGroupHandler) Delete(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.ICOASubGroupService.Delete(companyID, id)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record deleted.", Data: result})
}

func (h *COASubGroupHandler) SelectDropdownList(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	result, err := h.ICOASubGroupService.SelectDropdownList(companyID)
	if err != nil {
		return ctx.Status(500).JSON(response.JSON{Status: 500, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Dropdown list retrieved.", Data: result})
}
