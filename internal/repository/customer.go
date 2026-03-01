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
	Create(entity.Customer) (entity.Customer, error)
	FindAll(qp *util.QueryParams) ([]entity.Customer, int, error)
	FindById(id uuid.UUID) (entity.Customer, error)
	Update(entity.Customer) error
	Delete(id uuid.UUID) error
}

type CustomerRepository struct {
	Db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{Db: db}
}

func (r *CustomerRepository) Create(c entity.Customer) (entity.Customer, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&c).Error; err != nil {
		tx.Rollback()
		return c, err
	}
	tx.Commit()
	return c, nil
}

func (r *CustomerRepository) FindAll(qp *util.QueryParams) ([]entity.Customer, int, error) {
	var entities []entity.Customer
	var totalCount int64

	query := r.Db.Preload("COA").Model(&entity.Customer{})
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

func (r *CustomerRepository) FindById(id uuid.UUID) (entity.Customer, error) {
	var c entity.Customer
	if err := r.Db.Preload("COA").Where("id = ?", id).First(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (r *CustomerRepository) Update(c entity.Customer) error {
	tx := r.Db.Begin()
	if err := tx.Model(&c).Updates(c).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *CustomerRepository) Delete(id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Customer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
