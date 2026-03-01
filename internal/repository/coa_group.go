package repository

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaGroupSortColumns = map[string]string{
	"code":           "coa_groups.code",
	"name":           "coa_groups.name",
	"normal_balance": "coa_groups.normal_balance",
	"status":         "coa_groups.status",
	"created_at":     "coa_groups.created_at",
	"updated_at":     "coa_groups.updated_at",
}

type ICOAGroupRepository interface {
	Create(entity.COAGroup) (entity.COAGroup, error)
	FindAll(qp *util.QueryParams) ([]entity.COAGroup, int, error)
	FindById(entityId uuid.UUID) (entity.COAGroup, error)
	Update(entity.COAGroup) error
	Delete(entityId uuid.UUID) error
}

type COAGroupRepository struct {
	Db *gorm.DB
}

func NewCOAGroupRepository(Db *gorm.DB) ICOAGroupRepository {
	return &COAGroupRepository{Db: Db}
}

func (e *COAGroupRepository) Create(g entity.COAGroup) (entity.COAGroup, error) {
	tx := e.Db.Begin()
	if err := tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return g, err
	}
	tx.Commit()
	return g, nil
}

func (e *COAGroupRepository) FindAll(qp *util.QueryParams) ([]entity.COAGroup, int, error) {
	var entities []entity.COAGroup
	var totalCount int64

	query := e.Db.Model(&entity.COAGroup{})
	query = util.ApplySearch(query, qp)
	query = entity.COAGroup{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaGroupSortColumns, "coa_groups.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (e *COAGroupRepository) FindById(entityId uuid.UUID) (entity.COAGroup, error) {
	var g entity.COAGroup
	if err := e.Db.Where("id = ?", entityId).First(&g).Error; err != nil {
		return g, err
	}
	return g, nil
}

func (e *COAGroupRepository) Update(g entity.COAGroup) error {
	tx := e.Db.Begin()
	if err := tx.Model(&g).Updates(g).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (e *COAGroupRepository) Delete(entityId uuid.UUID) error {
	tx := e.Db.Begin()
	if err := tx.Where("id = ?", entityId).Delete(&entity.COAGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
