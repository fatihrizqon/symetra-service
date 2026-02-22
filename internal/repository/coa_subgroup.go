package repository

import (
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICOASubGroupRepository interface {
	Create(entity.COASubGroup) (entity.COASubGroup, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COASubGroupFilters) ([]entity.COASubGroup, int, error)
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

// Create implements ICOASubGroupRepository.
func (e *COASubGroupRepository) Create(entity entity.COASubGroup) (entity.COASubGroup, error) {
	tx := e.Db.Preload("Group").Begin()

	if err := tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	tx.Commit()
	return entity, nil
}

// FindAll implements ICOASubGroupRepository with pagination.
func (e *COASubGroupRepository) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COASubGroupFilters) ([]entity.COASubGroup, int, error) {
	var entities []entity.COASubGroup
	var totalCount int64

	query := e.Db.Preload("Group").Model(&entity.COASubGroup{})

	if search != "" && len(options.Fields) > 0 {
		orConditions := []string{}
		values := []interface{}{}
		for _, term := range strings.Split(search, ";") {
			term = strings.TrimSpace(term)
			for _, field := range options.Fields {
				orConditions = append(orConditions, "LOWER("+field+") LIKE LOWER(?)")
				values = append(values, "%"+term+"%")
			}
		}
		query = query.Where(strings.Join(orConditions, " OR "), values...)
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

// FindById implements ICOASubGroupRepository.
func (e *COASubGroupRepository) FindById(entityId uuid.UUID) (entity.COASubGroup, error) {
	var entity entity.COASubGroup
	if err := e.Db.Preload("Group").Where("id = ?", entityId).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update implements ICOASubGroupRepository.
func (e *COASubGroupRepository) Update(entity entity.COASubGroup) error {
	tx := e.Db.Preload("Group").Begin()

	if err := tx.Model(&entity).Updates(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete implements ICOASubGroupRepository.
func (e *COASubGroupRepository) Delete(entityId uuid.UUID) error {
	var entity entity.COASubGroup
	tx := e.Db.Begin()

	if err := tx.Where("id = ?", entityId).Delete(&entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
