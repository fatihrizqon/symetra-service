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

type JournalEntryHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewJournalEntryHandler(svc service.IJournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{IJournalEntryService: svc}
}

func (h *JournalEntryHandler) Create(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	callerID, _ := util.GetCallerID(ctx)

	var req request.JournalEntryCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IJournalEntryService.Create(companyID, req, callerID, entity.JournalTypeGeneral)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(201).JSON(response.JSON{Status: 201, Message: "Journal entry created.", Data: result})
}

func (h *JournalEntryHandler) FindAll(ctx *fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, entity.JournalEntry{}.SearchableFields())
	entities, totalCount, err := h.IJournalEntryService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(500).JSON(response.JSON{Status: 500, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "No records found.", Data: []response.JournalEntryResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Successfully retrieved records.", Data: entities, Meta: &meta})
}

func (h *JournalEntryHandler) FindById(ctx *fiber.Ctx) error {
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
	result, err := h.IJournalEntryService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(404).JSON(response.JSON{Status: 404, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record found.", Data: result})
}

func (h *JournalEntryHandler) Update(ctx *fiber.Ctx) error {
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
	var req request.JournalEntryUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id
	result, err := h.IJournalEntryService.Update(companyID, req)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Journal entry updated.", Data: result})
}

func (h *JournalEntryHandler) Delete(ctx *fiber.Ctx) error {
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
	if err := h.IJournalEntryService.Delete(companyID, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Journal entry deleted."})
}

func (h *JournalEntryHandler) Post(ctx *fiber.Ctx) error {
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
	if err := h.IJournalEntryService.Post(companyID, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Journal entry posted."})
}

func (h *JournalEntryHandler) Void(ctx *fiber.Ctx) error {
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
	if err := h.IJournalEntryService.Void(companyID, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Journal entry voided."})
}
