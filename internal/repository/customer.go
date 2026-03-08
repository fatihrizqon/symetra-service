package repository

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var customerSortColumns = map[string]string{
	"code":       "customers.code",
	"name":       "customers.name",
	"email":      "customers.email",
	"status":     "customers.status",
	"created_at": "customers.created_at",
	"updated_at": "customers.updated_at",
}

type ICustomerRepository interface {
	Create(c entity.Customer) (entity.Customer, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int, error)
	FindById(companyId, id uuid.UUID) (entity.Customer, error)
	Update(c entity.Customer) error
	Delete(companyId, id uuid.UUID) error
}

type CustomerRepository struct {
	Db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{Db: db}
}

func (r *CustomerRepository) Create(c entity.Customer) (entity.Customer, error) {
	c.Id = uuid.New()
	if err := r.Db.Create(&c).Error; err != nil {
		return c, err
	}
	return r.FindById(c.CompanyId, c.Id)
}

func (r *CustomerRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int, error) {
	var entities []entity.Customer
	var totalCount int64

	query := r.Db.Preload("COA").Model(&entity.Customer{}).
		Where("customers.company_id = ?", companyId)
	query = util.ApplySearch(query, qp)
	query = entity.Customer{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, customerSortColumns, "customers.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *CustomerRepository) FindById(companyId, id uuid.UUID) (entity.Customer, error) {
	var c entity.Customer
	if err := r.Db.Preload("COA").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (r *CustomerRepository) Update(c entity.Customer) error {
	return r.Db.Save(&c).Error
}

func (r *CustomerRepository) Delete(companyId, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyId).
		Delete(&entity.Customer{}).Error
}
