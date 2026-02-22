package repository

import (
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICustomerRepository interface {
	Create(entity.Customer) (entity.Customer, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.CustomerFilters) ([]entity.Customer, int, error)
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

func (r *CustomerRepository) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.CustomerFilters) ([]entity.Customer, int, error) {
	var entities []entity.Customer
	var totalCount int64

	query := r.Db.Preload("COA").Model(&entity.Customer{})

	if search != "" && len(options.Fields) > 0 {
		var conditions []string
		var values []interface{}
		for _, term := range strings.Split(search, ";") {
			term = strings.TrimSpace(term)
			for _, field := range options.Fields {
				conditions = append(conditions, "LOWER("+field+") LIKE LOWER(?)")
				values = append(values, "%"+term+"%")
			}
		}
		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return entities, 0, nil
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Limit(pageSize).Offset(offset).Find(&entities).Error; err != nil {
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
