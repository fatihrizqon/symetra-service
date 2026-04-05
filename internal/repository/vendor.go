package repository

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var vendorSortColumns = map[string]string{
	"code":       "vendors.code",
	"name":       "vendors.name",
	"email":      "vendors.email",
	"status":     "vendors.status",
	"created_at": "vendors.created_at",
	"updated_at": "vendors.updated_at",
}

type IVendorRepository interface {
	Create(entity.Vendor) (entity.Vendor, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int, error)
	FindById(companyID, id uuid.UUID) (entity.Vendor, error)
	Update(entity.Vendor) error
	Delete(companyID, id uuid.UUID) error
}

type VendorRepository struct {
	Db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) IVendorRepository {
	return &VendorRepository{Db: db}
}

func (r *VendorRepository) Create(v entity.Vendor) (entity.Vendor, error) {
	v.Id = uuid.New()
	if err := r.Db.Create(&v).Error; err != nil {
		return v, err
	}
	return r.FindById(v.CompanyId, v.Id)
}

func (r *VendorRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int, error) {
	var entities []entity.Vendor
	var totalCount int64

	query := r.Db.Preload("COA").Model(&entity.Vendor{}).
		Where("vendors.company_id = ?", companyID)
	query = util.ApplySearch(query, qp)
	query = entity.Vendor{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, vendorSortColumns, "vendors.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *VendorRepository) FindById(companyID, id uuid.UUID) (entity.Vendor, error) {
	var v entity.Vendor
	if err := r.Db.Preload("COA").
		Where("id = ? AND company_id = ?", id, companyID).
		First(&v).Error; err != nil {
		return v, err
	}
	return v, nil
}

func (r *VendorRepository) Update(v entity.Vendor) error {
	return r.Db.Save(&v).Error
}

func (r *VendorRepository) Delete(companyID, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyID).
		Delete(&entity.Vendor{}).Error
}
