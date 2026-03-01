package repository

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaSubGroupSortColumns = map[string]string{
	"code":       "coa_subgroups.code",
	"name":       "coa_subgroups.name",
	"status":     "coa_subgroups.status",
	"created_at": "coa_subgroups.created_at",
	"updated_at": "coa_subgroups.updated_at",
}

type ICOASubGroupRepository interface {
	Create(entity.COASubGroup) (entity.COASubGroup, error)
	FindAll(qp *util.QueryParams) ([]entity.COASubGroup, int, error)
	FindById(entityId uuid.UUID) (entity.COASubGroup, error)
	Update(entity.COASubGroup) error
	Delete(entityId uuid.UUID) error
}

type COASubGroupRepository struct {
	Db *gorm.DB
}

func NewCOASubGroupRepository(Db *gorm.DB) ICOASubGroupRepository {
	return &COASubGroupRepository{Db: Db}
}

func (e *COASubGroupRepository) Create(sg entity.COASubGroup) (entity.COASubGroup, error) {
	tx := e.Db.Preload("Group").Begin()
	if err := tx.Create(&sg).Error; err != nil {
		tx.Rollback()
		return sg, err
	}
	tx.Commit()
	return sg, nil
}

func (e *COASubGroupRepository) FindAll(qp *util.QueryParams) ([]entity.COASubGroup, int, error) {
	var entities []entity.COASubGroup
	var totalCount int64

	query := e.Db.Preload("Group").Model(&entity.COASubGroup{})
	query = util.ApplySearch(query, qp)
	query = entity.COASubGroup{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSubGroupSortColumns, "coa_subgroups.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (e *COASubGroupRepository) FindById(entityId uuid.UUID) (entity.COASubGroup, error) {
	var sg entity.COASubGroup
	if err := e.Db.Preload("Group").Where("id = ?", entityId).First(&sg).Error; err != nil {
		return sg, err
	}
	return sg, nil
}

func (e *COASubGroupRepository) Update(sg entity.COASubGroup) error {
	tx := e.Db.Preload("Group").Begin()
	if err := tx.Model(&sg).Updates(sg).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (e *COASubGroupRepository) Delete(entityId uuid.UUID) error {
	tx := e.Db.Begin()
	if err := tx.Where("id = ?", entityId).Delete(&entity.COASubGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
