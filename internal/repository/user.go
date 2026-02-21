package repository

import (
	"context"
	"strings"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(entity.User) (entity.User, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.UserFilters) ([]entity.User, int, error)
	FindById(entityId uuid.UUID) (entity.User, error)
	Update(entity.User) error
	Delete(entityId uuid.UUID) error

	CountAll(ctx context.Context) (int64, error)
	CountBetween(ctx context.Context, from time.Time, to time.Time) (int64, error)
	CountVerified(ctx context.Context) (int64, error)
}

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) IUserRepository {
	return &UserRepository{Db: Db}
}

// Create implements IUserRepository.
func (e *UserRepository) Create(entity entity.User) (entity.User, error) {
	tx := e.Db.Begin()

	if err := tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	tx.Commit()
	return entity, nil
}

// FindAll implements IUserRepository with pagination.
func (e *UserRepository) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.UserFilters) ([]entity.User, int, error) {
	var entities []entity.User
	var totalCount int64

	query := e.Db.Model(&entity.User{})

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

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Verified != nil {
		verified := strings.ToLower(strings.TrimSpace(*filters.Verified))
		if verified == "true" {
			query = query.Where("email_verified_at IS NOT NULL")
		} else if verified == "false" {
			query = query.Where("email_verified_at IS NULL")
		}
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

// FindById implements IUserRepository.
func (e *UserRepository) FindById(entityId uuid.UUID) (entity.User, error) {
	var entity entity.User
	if err := e.Db.Where("id = ?", entityId).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update implements IUserRepository.
func (e *UserRepository) Update(entity entity.User) error {
	tx := e.Db.Begin()

	if err := tx.Model(&entity).Updates(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete implements IUserRepository.
func (e *UserRepository) Delete(entityId uuid.UUID) error {
	var entity entity.User
	tx := e.Db.Begin()

	if err := tx.Where("id = ?", entityId).Delete(&entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// CountAll implements IUserRepository.
func (e *UserRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := e.Db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountBetween(ctx, from time.Time, to time.Time) implements IUserRepository.
func (e *UserRepository) CountBetween(ctx context.Context, from time.Time, to time.Time) (int64, error) {
	var count int64
	if err := e.Db.WithContext(ctx).Model(&entity.User{}).Where("created_at BETWEEN ? AND ?", from, to).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountVerified implements IUserRepository.
func (e *UserRepository) CountVerified(ctx context.Context) (int64, error) {
	var count int64
	if err := e.Db.WithContext(ctx).Model(&entity.User{}).Where("email_verified_at IS NOT NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
