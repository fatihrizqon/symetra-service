package repository

import (
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IVendorRepository interface {
	Create(entity.Vendor) (entity.Vendor, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.VendorFilters) ([]entity.Vendor, int, error)
	FindById(id uuid.UUID) (entity.Vendor, error)
	Update(entity.Vendor) error
	Delete(id uuid.UUID) error
}

type VendorRepository struct {
	Db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) IVendorRepository {
	return &VendorRepository{Db: db}
}

func (r *VendorRepository) Create(v entity.Vendor) (entity.Vendor, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&v).Error; err != nil {
		tx.Rollback()
		return v, err
	}
	tx.Commit()
	return v, nil
}

func (r *VendorRepository) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.VendorFilters) ([]entity.Vendor, int, error) {
	var entities []entity.Vendor
	var totalCount int64

	query := r.Db.Preload("COA").Model(&entity.Vendor{})

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

func (r *VendorRepository) FindById(id uuid.UUID) (entity.Vendor, error) {
	var v entity.Vendor
	if err := r.Db.Preload("COA").Where("id = ?", id).First(&v).Error; err != nil {
		return v, err
	}
	return v, nil
}

func (r *VendorRepository) Update(v entity.Vendor) error {
	tx := r.Db.Begin()
	if err := tx.Model(&v).Updates(v).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *VendorRepository) Delete(id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Vendor{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
