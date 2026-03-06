package repository

import (
	"errors"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var companySortColumns = map[string]string{
	"name":       "companies.name",
	"status":     "companies.status",
	"created_at": "companies.created_at",
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

type ICompanyRepository interface {
	Create(c entity.Company) (entity.Company, error)
	FindAll(qp *util.QueryParams) ([]entity.Company, int, error)
	FindById(id uuid.UUID) (entity.Company, error)
	FindByUserId(userID uuid.UUID, qp *util.QueryParams) ([]entity.Company, int, error)
	Update(c entity.Company) error
	SoftDelete(id uuid.UUID) error
}

type ICompanyMemberRepository interface {
	AddMember(m entity.CompanyMember) (entity.CompanyMember, error)
	FindMembership(userID, companyID uuid.UUID) (entity.CompanyMember, error)
	FindMembersByCompany(companyID uuid.UUID) ([]entity.CompanyMember, error)
	FindCompaniesByUser(userID uuid.UUID) ([]entity.CompanyMember, error)
	UpdateRole(userID, companyID uuid.UUID, role entity.CompanyRole) error
	RemoveMember(userID, companyID uuid.UUID) error
	IsSuperadmin(userID uuid.UUID) (bool, error)
	CountOwners(companyID uuid.UUID) (int64, error)
}

// ─── Company Repository ───────────────────────────────────────────────────────

type CompanyRepository struct {
	Db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ICompanyRepository {
	return &CompanyRepository{Db: db}
}

func (r *CompanyRepository) Create(c entity.Company) (entity.Company, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&c).Error; err != nil {
		tx.Rollback()
		return c, err
	}
	tx.Commit()
	return c, nil
}

func (r *CompanyRepository) FindAll(qp *util.QueryParams) ([]entity.Company, int, error) {
	var entities []entity.Company
	var totalCount int64

	query := r.Db.Model(&entity.Company{}).Where("deleted_at IS NULL")
	query = util.ApplySearch(query, qp)
	query = entity.Company{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, companySortColumns, "companies.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *CompanyRepository) FindById(id uuid.UUID) (entity.Company, error) {
	var c entity.Company
	if err := r.Db.Where("id = ? AND deleted_at IS NULL", id).First(&c).Error; err != nil {
		return c, errors.New("company not found")
	}
	return c, nil
}

// FindByUserId returns all companies a specific user belongs to.
func (r *CompanyRepository) FindByUserId(userID uuid.UUID, qp *util.QueryParams) ([]entity.Company, int, error) {
	var entities []entity.Company
	var totalCount int64

	query := r.Db.Model(&entity.Company{}).
		Joins("JOIN company_members ON company_members.company_id = companies.id").
		Where("company_members.user_id = ? AND companies.deleted_at IS NULL", userID)

	query = util.ApplySearch(query, qp)
	query = entity.Company{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, companySortColumns, "companies.name")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *CompanyRepository) Update(c entity.Company) error {
	tx := r.Db.Begin()
	if err := tx.Model(&c).Updates(map[string]interface{}{
		"name":       c.Name,
		"legal_name": c.LegalName,
		"tax_id":     c.TaxID,
		"address":    c.Address,
		"phone":      c.Phone,
		"email":      c.Email,
		"industry":   c.Industry,
		"currency":   c.Currency,
		"status":     c.Status,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *CompanyRepository) SoftDelete(id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Model(&entity.Company{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// ─── Company Member Repository ────────────────────────────────────────────────

type CompanyMemberRepository struct {
	Db *gorm.DB
}

func NewCompanyMemberRepository(db *gorm.DB) ICompanyMemberRepository {
	return &CompanyMemberRepository{Db: db}
}

func (r *CompanyMemberRepository) AddMember(m entity.CompanyMember) (entity.CompanyMember, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return m, err
	}
	tx.Commit()
	return m, nil
}

func (r *CompanyMemberRepository) FindMembership(userID, companyID uuid.UUID) (entity.CompanyMember, error) {
	var m entity.CompanyMember
	err := r.Db.Where("user_id = ? AND company_id = ?", userID, companyID).First(&m).Error
	if err != nil {
		return m, errors.New("membership not found")
	}
	return m, nil
}

func (r *CompanyMemberRepository) FindMembersByCompany(companyID uuid.UUID) ([]entity.CompanyMember, error) {
	var members []entity.CompanyMember
	err := r.Db.Preload("User").
		Where("company_id = ?", companyID).
		Order("joined_at ASC").
		Find(&members).Error
	return members, err
}

func (r *CompanyMemberRepository) FindCompaniesByUser(userID uuid.UUID) ([]entity.CompanyMember, error) {
	var members []entity.CompanyMember
	err := r.Db.Preload("Company").
		Where("user_id = ?", userID).
		Order("joined_at ASC").
		Find(&members).Error
	return members, err
}

func (r *CompanyMemberRepository) UpdateRole(userID, companyID uuid.UUID, role entity.CompanyRole) error {
	return r.Db.Model(&entity.CompanyMember{}).
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Update("role", role).Error
}

func (r *CompanyMemberRepository) RemoveMember(userID, companyID uuid.UUID) error {
	return r.Db.
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Delete(&entity.CompanyMember{}).Error
}

// IsSuperadmin checks if a user has the superadmin platform role.
// Superadmin is stored as a special company_member record with company_id = uuid.Nil
// OR as a boolean field on the user — we use a dedicated platform membership row:
// company_id = '00000000-0000-0000-0000-000000000000', role = 'superadmin'
func (r *CompanyMemberRepository) IsSuperadmin(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.Db.Model(&entity.CompanyMember{}).
		Where("user_id = ? AND role = ?", userID, entity.RoleSuperadmin).
		Count(&count).Error
	return count > 0, err
}

func (r *CompanyMemberRepository) CountOwners(companyID uuid.UUID) (int64, error) {
	var count int64
	err := r.Db.Model(&entity.CompanyMember{}).
		Where("company_id = ? AND role = ?", companyID, entity.RoleOwner).
		Count(&count).Error
	return count, err
}
