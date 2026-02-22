package repository

import (
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICOAGroupRepository interface {
	Create(entity.COAGroup) (entity.COAGroup, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions) ([]entity.COAGroup, int, error)
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

// Create implements ICOAGroupRepository.
func (e *COAGroupRepository) Create(entity entity.COAGroup) (entity.COAGroup, error) {
	tx := e.Db.Begin()

	if err := tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	tx.Commit()
	return entity, nil
}

// FindAll implements ICOAGroupRepository with pagination.
func (e *COAGroupRepository) FindAll(page, pageSize int, search string, options util.SearchOptions) ([]entity.COAGroup, int, error) {
	var entities []entity.COAGroup
	var totalCount int64

	query := e.Db.Model(&entity.COAGroup{})

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

// FindById implements ICOAGroupRepository.
func (e *COAGroupRepository) FindById(entityId uuid.UUID) (entity.COAGroup, error) {
	var entity entity.COAGroup
	if err := e.Db.Where("id = ?", entityId).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update implements ICOAGroupRepository.
func (e *COAGroupRepository) Update(entity entity.COAGroup) error {
	tx := e.Db.Begin()

	if err := tx.Model(&entity).Updates(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete implements ICOAGroupRepository.
func (e *COAGroupRepository) Delete(entityId uuid.UUID) error {
	var entity entity.COAGroup
	tx := e.Db.Begin()

	if err := tx.Where("id = ?", entityId).Delete(&entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
