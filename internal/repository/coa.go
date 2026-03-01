package repository

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaSortColumns = map[string]string{
	"code":       "chart_of_accounts.code",
	"name":       "chart_of_accounts.name",
	"status":     "chart_of_accounts.status",
	"created_at": "chart_of_accounts.created_at",
	"updated_at": "chart_of_accounts.updated_at",
}

type ICOARepository interface {
	Create(entity.COA) (entity.COA, error)
	FindAll(qp *util.QueryParams) ([]entity.COA, int, error)
	FindById(entityId uuid.UUID) (entity.COA, error)
	Update(entity.COA) error
	Delete(entityId uuid.UUID) error
}

type COARepository struct {
	Db *gorm.DB
}

func NewCOARepository(Db *gorm.DB) ICOARepository {
	return &COARepository{Db: Db}
}

func (e *COARepository) Create(c entity.COA) (entity.COA, error) {
	tx := e.Db.Preload("SubGroup").Begin()
	if err := tx.Create(&c).Error; err != nil {
		tx.Rollback()
		return c, err
	}
	tx.Commit()
	return c, nil
}

func (e *COARepository) FindAll(qp *util.QueryParams) ([]entity.COA, int, error) {
	var entities []entity.COA
	var totalCount int64

	query := e.Db.Model(&entity.COA{})
	query = util.ApplySearch(query, qp)
	query = entity.COA{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSortColumns, "chart_of_accounts.created_at")
	query = util.ApplyPagination(query, qp)

	err := query.
		Preload("SubGroup", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name", "status", "group_id", "created_at", "updated_at").
				Preload("Group", func(db2 *gorm.DB) *gorm.DB {
					return db2.Select("id", "code", "name", "normal_balance", "status", "created_at", "updated_at")
				})
		}).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (e *COARepository) FindById(entityId uuid.UUID) (entity.COA, error) {
	var c entity.COA
	if err := e.Db.Preload("SubGroup").Where("id = ?", entityId).First(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (e *COARepository) Update(c entity.COA) error {
	tx := e.Db.Preload("SubGroup").Begin()
	if err := tx.Model(&c).Updates(c).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (e *COARepository) Delete(entityId uuid.UUID) error {
	tx := e.Db.Begin()
	if err := tx.Where("id = ?", entityId).Delete(&entity.COA{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
