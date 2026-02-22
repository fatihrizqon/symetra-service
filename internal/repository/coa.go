package repository

import (
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICOARepository interface {
	Create(entity.COA) (entity.COA, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions) ([]entity.COA, int, error)
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

// Create implements ICOARepository.
func (e *COARepository) Create(entity entity.COA) (entity.COA, error) {
	tx := e.Db.Preload("SubGroup").Begin()

	if err := tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	tx.Commit()
	return entity, nil
}

// FindAll implements ICOARepository with pagination.
func (e *COARepository) FindAll(page, pageSize int, search string, options util.SearchOptions) ([]entity.COA, int, error) {
	var entities []entity.COA
	var totalCount int64

	query := e.Db.Model(&entity.COA{})

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

	if err := query.
		Preload("SubGroup", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name", "status", "group_id", "created_at", "updated_at").
				Preload("Group", func(db2 *gorm.DB) *gorm.DB {
					return db2.Select("id", "code", "name", "normal_balance", "status", "created_at", "updated_at")
				})
		}).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

// FindById implements ICOARepository.
func (e *COARepository) FindById(entityId uuid.UUID) (entity.COA, error) {
	var entity entity.COA
	if err := e.Db.Preload("SubGroup").Where("id = ?", entityId).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update implements ICOARepository.
func (e *COARepository) Update(entity entity.COA) error {
	tx := e.Db.Preload("SubGroup").Begin()

	if err := tx.Model(&entity).Updates(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete implements ICOARepository.
func (e *COARepository) Delete(entityId uuid.UUID) error {
	var entity entity.COA
	tx := e.Db.Begin()

	if err := tx.Where("id = ?", entityId).Delete(&entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
