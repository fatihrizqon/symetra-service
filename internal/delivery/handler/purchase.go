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

// ─── Purchase Order Handler ───────────────────────────────────────────────────

type PurchaseOrderHandler struct {
	svc service.IPurchaseOrderService
}

func NewPurchaseOrderHandler(svc service.IPurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{svc: svc}
}

func (h *PurchaseOrderHandler) Create(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	var req request.POCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	result, err := h.svc.Create(companyId, req, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(201).JSON(response.JSON{Status: 201, Message: "Purchase order created.", Data: result})
}

func (h *PurchaseOrderHandler) FindAll(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.PurchaseOrder{}.SearchableFields())
	entities, totalCount, err := h.svc.FindAllTyped(companyId, qp)
	if err != nil {
		return ctx.Status(500).JSON(response.JSON{Status: 500, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "No records found.", Data: []response.POResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Successfully retrieved records.", Data: entities, Meta: &meta})
}

func (h *PurchaseOrderHandler) FindById(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	result, err := h.svc.FindById(companyId, id)
	if err != nil {
		return ctx.Status(404).JSON(response.JSON{Status: 404, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record found.", Data: result})
}

func (h *PurchaseOrderHandler) Update(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	var req request.POUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	req.Id = id
	result, err := h.svc.Update(companyId, req)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Purchase order updated.", Data: result})
}

func (h *PurchaseOrderHandler) Delete(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	if err := h.svc.Delete(companyId, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Purchase order deleted."})
}

func (h *PurchaseOrderHandler) Send(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	if err := h.svc.Send(companyId, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Purchase order sent."})
}

func (h *PurchaseOrderHandler) Approve(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	if err := h.svc.Approve(companyId, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Purchase order approved."})
}

func (h *PurchaseOrderHandler) Decline(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	if err := h.svc.Decline(companyId, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Purchase order declined."})
}

// ─── Bill Handler ─────────────────────────────────────────────────────────────

type BillHandler struct {
	svc service.IBillService
}

func NewBillHandler(svc service.IBillService) *BillHandler {
	return &BillHandler{svc: svc}
}

func (h *BillHandler) Create(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	var req request.BillCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	result, err := h.svc.Create(companyId, req, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(201).JSON(response.JSON{Status: 201, Message: "Bill created.", Data: result})
}

func (h *BillHandler) CreateFromPO(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	poId, err := uuid.Parse(ctx.Params("po_id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid po_id"})
	}
	result, err := h.svc.CreateFromPO(companyId, poId, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(201).JSON(response.JSON{Status: 201, Message: "Bill created from purchase order.", Data: result})
}

func (h *BillHandler) FindAll(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.Bill{}.SearchableFields())
	entities, totalCount, err := h.svc.FindAllTyped(companyId, qp)
	if err != nil {
		return ctx.Status(500).JSON(response.JSON{Status: 500, Message: err.Error()})
	}
	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "No records found.", Data: []response.BillResponse{}})
	}
	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Successfully retrieved records.", Data: entities, Meta: &meta})
}

func (h *BillHandler) FindById(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	result, err := h.svc.FindById(companyId, id)
	if err != nil {
		return ctx.Status(404).JSON(response.JSON{Status: 404, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Record found.", Data: result})
}

func (h *BillHandler) Update(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	var req request.BillUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	req.Id = id
	result, err := h.svc.Update(companyId, req)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Bill updated.", Data: result})
}

func (h *BillHandler) Delete(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	if err := h.svc.Delete(companyId, id); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Bill deleted."})
}

func (h *BillHandler) Confirm(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	result, err := h.svc.Confirm(companyId, id, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Bill confirmed. Expense journal created.", Data: result})
}

func (h *BillHandler) AddPayment(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	var req request.BillPaymentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	result, err := h.svc.AddPayment(companyId, id, req, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Payment recorded.", Data: result})
}

func (h *BillHandler) Cancel(ctx *fiber.Ctx) error {
	companyId, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	callerId, _ := util.GetCallerID(ctx)
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: "invalid id"})
	}
	result, err := h.svc.Cancel(companyId, id, callerId)
	if err != nil {
		return ctx.Status(400).JSON(response.JSON{Status: 400, Message: err.Error()})
	}
	return ctx.Status(200).JSON(response.JSON{Status: 200, Message: "Bill cancelled.", Data: result})
}
