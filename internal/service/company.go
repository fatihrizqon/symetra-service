package service

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICompanyService interface {
	// Company CRUD
	Create(req request.CompanyCreateRequest, createdBy uuid.UUID) (entity.Company, error)
	FindAll(qp *util.QueryParams) ([]response.CompanyResponse, int, error) // superadmin only
	FindMyCompanies(userID uuid.UUID) ([]response.MyCompanyResponse, error)
	FindById(id uuid.UUID) (response.CompanyResponse, error)
	Update(req request.CompanyUpdateRequest) (entity.Company, error)
	Delete(companyID uuid.UUID, callerID uuid.UUID) error

	// Member management
	AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (entity.CompanyMember, error)
	FindMembers(companyID uuid.UUID) ([]response.CompanyMemberResponse, error)
	UpdateMemberRole(req request.UpdateMemberRoleRequest, callerID uuid.UUID) error
	RemoveMember(companyID, targetUserID, callerID uuid.UUID) error

	// Platform guards
	// AssertSuperadmin returns error jika user bukan superadmin.
	// Digunakan oleh handler yang tidak bisa menggunakan CompanyMiddleware.
	AssertSuperadmin(userID uuid.UUID) error
}

type CompanyService struct {
	companyRepo repository.ICompanyRepository
	memberRepo  repository.ICompanyMemberRepository
	userRepo    repository.IUserRepository
	validate    *validator.Validate
}

func NewCompanyService(
	companyRepo repository.ICompanyRepository,
	memberRepo repository.ICompanyMemberRepository,
	userRepo repository.IUserRepository,
	validate *validator.Validate,
) ICompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		validate:    validate,
	}
}

// ── Company CRUD ──────────────────────────────────────────────────────────────

func (s *CompanyService) Create(req request.CompanyCreateRequest, createdBy uuid.UUID) (entity.Company, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Company{}, err
	}

	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	c := entity.Company{
		Name:      req.Name,
		LegalName: req.LegalName,
		TaxID:     req.TaxID,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Industry:  req.Industry,
		Currency:  currency,
		Status:    1,
		CreatedBy: createdBy,
	}

	created, err := s.companyRepo.Create(c)
	if err != nil {
		return created, err
	}

	// Auto-assign creator as Owner
	_, err = s.memberRepo.AddMember(entity.CompanyMember{
		CompanyId: created.Id,
		UserId:    createdBy,
		Role:      entity.RoleOwner,
		InvitedBy: createdBy,
	})
	if err != nil {
		// Company created but owner assignment failed — log but don't block
		return created, fmt.Errorf("company created but owner assignment failed: %w", err)
	}

	return created, nil
}

func (s *CompanyService) FindAll(qp *util.QueryParams) ([]response.CompanyResponse, int, error) {
	entities, totalCount, err := s.companyRepo.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.CompanyResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.CompanyResponse, 0, len(entities))
	for _, c := range entities {
		resps = append(resps, mapCompany(c))
	}
	return resps, totalCount, nil
}

func (s *CompanyService) FindMyCompanies(userID uuid.UUID) ([]response.MyCompanyResponse, error) {
	memberships, err := s.memberRepo.FindCompaniesByUser(userID)
	if err != nil {
		return nil, err
	}

	resps := make([]response.MyCompanyResponse, 0, len(memberships))
	for _, m := range memberships {
		resps = append(resps, response.MyCompanyResponse{
			CompanyResponse: mapCompany(m.Company),
			Role:            string(m.Role),
		})
	}
	return resps, nil
}

func (s *CompanyService) FindById(id uuid.UUID) (response.CompanyResponse, error) {
	c, err := s.companyRepo.FindById(id)
	if err != nil {
		return response.CompanyResponse{}, err
	}
	return mapCompany(c), nil
}

func (s *CompanyService) Update(req request.CompanyUpdateRequest) (entity.Company, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Company{}, err
	}
	c, err := s.companyRepo.FindById(req.Id)
	if err != nil {
		return c, err
	}
	c.Name = req.Name
	c.LegalName = req.LegalName
	c.TaxID = req.TaxID
	c.Address = req.Address
	c.Phone = req.Phone
	c.Email = req.Email
	c.Industry = req.Industry
	if req.Currency != "" {
		c.Currency = req.Currency
	}
	return c, s.companyRepo.Update(c)
}

// Delete soft-deletes a company. Only Owner or Superadmin can do this.
func (s *CompanyService) Delete(companyID uuid.UUID, callerID uuid.UUID) error {
	if _, err := s.companyRepo.FindById(companyID); err != nil {
		return err
	}
	return s.companyRepo.SoftDelete(companyID)
}

// ── Member Management ─────────────────────────────────────────────────────────

// AssignMember adds a registered user to a company with a given role.
// The target user must already have an account in the system.
func (s *CompanyService) AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (entity.CompanyMember, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.CompanyMember{}, err
	}

	// Verify target user exists
	if _, err := s.userRepo.FindById(req.UserID); err != nil {
		return entity.CompanyMember{}, errors.New("target user not found")
	}

	// Verify company exists
	if _, err := s.companyRepo.FindById(req.CompanyId); err != nil {
		return entity.CompanyMember{}, err
	}

	// Check not already a member
	if _, err := s.memberRepo.FindMembership(req.UserID, req.CompanyId); err == nil {
		return entity.CompanyMember{}, errors.New("user is already a member of this company")
	}

	// Prevent assigning superadmin via API
	if entity.CompanyRole(req.Role) == entity.RoleSuperadmin {
		return entity.CompanyMember{}, errors.New("superadmin role cannot be assigned via API")
	}

	return s.memberRepo.AddMember(entity.CompanyMember{
		CompanyId: req.CompanyId,
		UserId:    req.UserID,
		Role:      entity.CompanyRole(req.Role),
		InvitedBy: invitedBy,
	})
}

func (s *CompanyService) FindMembers(companyID uuid.UUID) ([]response.CompanyMemberResponse, error) {
	members, err := s.memberRepo.FindMembersByCompany(companyID)
	if err != nil {
		return nil, err
	}

	resps := make([]response.CompanyMemberResponse, 0, len(members))
	for _, m := range members {
		resps = append(resps, response.CompanyMemberResponse{
			Id:        m.Id,
			CompanyId: m.CompanyId,
			UserId:    m.UserId,
			Username:  m.User.Username,
			Name:      m.User.Name,
			Email:     m.User.Email,
			Role:      string(m.Role),
			JoinedAt:  m.JoinedAt,
		})
	}
	return resps, nil
}

// UpdateMemberRole changes a member's role.
// Rules:
//   - Cannot demote/promote yourself
//   - Cannot change another owner's role unless you are superadmin
//   - Cannot assign superadmin via API
func (s *CompanyService) UpdateMemberRole(req request.UpdateMemberRoleRequest, callerID uuid.UUID) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}
	if req.UserID == callerID {
		return errors.New("you cannot change your own role")
	}
	if entity.CompanyRole(req.Role) == entity.RoleSuperadmin {
		return errors.New("superadmin role cannot be assigned via API")
	}

	// Check target is a member
	existing, err := s.memberRepo.FindMembership(req.UserID, req.CompanyId)
	if err != nil {
		return errors.New("target user is not a member of this company")
	}

	// Only superadmin can change another owner's role
	callerMembership, _ := s.memberRepo.FindMembership(callerID, req.CompanyId)
	if existing.Role == entity.RoleOwner && callerMembership.Role != entity.RoleSuperadmin {
		return errors.New("only superadmin can change an owner's role")
	}

	return s.memberRepo.UpdateRole(req.UserID, req.CompanyId, entity.CompanyRole(req.Role))
}

// RemoveMember removes a user from a company.
// Rules:
//   - Cannot remove yourself
//   - Cannot remove the last owner
func (s *CompanyService) RemoveMember(companyID, targetUserID, callerID uuid.UUID) error {
	if targetUserID == callerID {
		return errors.New("you cannot remove yourself from a company")
	}

	existing, err := s.memberRepo.FindMembership(targetUserID, companyID)
	if err != nil {
		return errors.New("target user is not a member of this company")
	}

	// Protect last owner
	if existing.Role == entity.RoleOwner {
		count, err := s.memberRepo.CountOwners(companyID)
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot remove the last owner of a company")
		}
	}

	return s.memberRepo.RemoveMember(targetUserID, companyID)
}

// ── Mapper ────────────────────────────────────────────────────────────────────

func mapCompany(c entity.Company) response.CompanyResponse {
	return response.CompanyResponse{
		Id:        c.Id,
		Name:      c.Name,
		LegalName: c.LegalName,
		TaxID:     c.TaxID,
		Address:   c.Address,
		Phone:     c.Phone,
		Email:     c.Email,
		Industry:  c.Industry,
		Currency:  c.Currency,
		Status:    c.Status,
		CreatedBy: c.CreatedBy,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// AssertSuperadmin memverifikasi bahwa userID adalah superadmin platform.
// Digunakan di handler yang tidak bisa menggunakan CompanyMiddleware.
func (s *CompanyService) AssertSuperadmin(userID uuid.UUID) error {
	ok, err := s.memberRepo.IsSuperadmin(userID)
	if err != nil {
		return errors.New("failed to verify superadmin status")
	}
	if !ok {
		return errors.New("superadmin access required")
	}
	return nil
}
