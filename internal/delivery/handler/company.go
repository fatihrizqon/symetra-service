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

type CompanyHandler struct {
	ICompanyService service.ICompanyService
}

func NewCompanyHandler(svc service.ICompanyService) *CompanyHandler {
	return &CompanyHandler{ICompanyService: svc}
}

// ── Company CRUD ──────────────────────────────────────────────────────────────

// CreateCompany godoc
// @Summary Create a new company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CompanyCreateRequest true "Company Create Request"
// @Success 201 {object} response.JSON
// @Router /api/v1/companies [post]
func (h *CompanyHandler) CreateCompany(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	var req request.CompanyCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.Create(req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  400,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  201,
		Message: "Company created successfully.",
		Data:    result,
	})
}

// FindAllCompanies godoc — Superadmin only
// @Summary List all companies (superadmin only)
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies [get]
func (h *CompanyHandler) FindAllCompanies(ctx *fiber.Ctx) error {
	// Guard: endpoint ini tidak pakai CompanyMiddleware (tidak ada X-Company-ID),
	// jadi superadmin check dilakukan langsung di handler via ICompanyService.
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	if err := h.ICompanyService.AssertSuperadmin(callerID); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(response.JSON{
			Status:  403,
			Message: "access denied: superadmin only",
		})
	}

	qp := util.ParseQueryParams(ctx, entity.Company{}.SearchableFields())

	entities, totalCount, err := h.ICompanyService.FindAll(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: 500, Message: "Failed to retrieve companies", Errors: err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status: 200, Message: "No records found.", Data: []response.CompanyResponse{},
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Successfully retrieved companies.", Data: entities, Meta: &meta,
	})
}

// FindMyCompanies godoc
// @Summary Get all companies the authenticated user belongs to
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/mine [get]
func (h *CompanyHandler) FindMyCompanies(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	result, err := h.ICompanyService.FindMyCompanies(callerID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: 500, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Successfully retrieved your companies.", Data: result,
	})
}

// FindCompanyById godoc
// @Summary Get company by ID
// @Tags Company
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Router /api/v1/companies/{id} [get]
func (h *CompanyHandler) FindCompanyById(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.FindById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status: 404, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Successfully retrieved company.", Data: result,
	})
}

// UpdateCompany godoc
// @Summary Update company details
// @Tags Company
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Router /api/v1/companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var req request.CompanyUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.Id = id

	result, err := h.ICompanyService.Update(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status: 400, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Company updated.", Data: result,
	})
}

// DeleteCompany godoc
// @Summary Soft-delete a company (owner only)
// @Tags Company
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Router /api/v1/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.ICompanyService.Delete(id, callerID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status: 400, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Company deleted.",
	})
}

// ── Member Management ─────────────────────────────────────────────────────────

// AssignMember godoc
// @Summary Add a user to a company
// @Tags Company Members
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Router /api/v1/companies/{id}/members [post]
func (h *CompanyHandler) AssignMember(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	companyID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var req request.AssignMemberRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.CompanyId = companyID

	result, err := h.ICompanyService.AssignMember(req, callerID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status: 400, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status: 201, Message: "Member added successfully.", Data: result,
	})
}

// FindMembers godoc
// @Summary List all members of a company
// @Tags Company Members
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Router /api/v1/companies/{id}/members [get]
func (h *CompanyHandler) FindMembers(ctx *fiber.Ctx) error {
	companyID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.FindMembers(companyID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: 500, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Successfully retrieved members.", Data: result,
	})
}

// UpdateMemberRole godoc
// @Summary Change a member's role within a company
// @Tags Company Members
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Param user_id path string true "User ID"
// @Router /api/v1/companies/{id}/members/{user_id}/role [put]
func (h *CompanyHandler) UpdateMemberRole(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	companyID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	targetUserID, err := uuid.Parse(ctx.Params("user_id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var req request.UpdateMemberRoleRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.CompanyId = companyID
	req.UserID = targetUserID

	if err := h.ICompanyService.UpdateMemberRole(req, callerID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status: 400, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Member role updated.",
	})
}

// RemoveMember godoc
// @Summary Remove a user from a company
// @Tags Company Members
// @Security BearerAuth
// @Param X-Company-ID header string true "Company ID"
// @Param id path string true "Company ID"
// @Param user_id path string true "User ID to remove"
// @Router /api/v1/companies/{id}/members/{user_id} [delete]
func (h *CompanyHandler) RemoveMember(ctx *fiber.Ctx) error {
	callerID, err := util.GetCallerID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	companyID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	targetUserID, err := uuid.Parse(ctx.Params("user_id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.ICompanyService.RemoveMember(companyID, targetUserID, callerID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status: 400, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: 200, Message: "Member removed.",
	})
}
