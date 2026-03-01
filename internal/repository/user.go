package repository

import (
	"context"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userSortColumns is the whitelist of sortable columns for the users table.
// Key = public-facing ?sort= value, Value = safe SQL column expression.
var userSortColumns = map[string]string{
	"name":       "users.name",
	"username":   "users.username",
	"email":      "users.email",
	"status":     "users.status",
	"created_at": "users.created_at",
	"updated_at": "users.updated_at",
}

type IUserRepository interface {
	Create(entity.User) (entity.User, error)
	FindAll(qp *util.QueryParams) ([]entity.User, int, error)
	FindById(id uuid.UUID) (entity.User, error)
	Update(entity.User) error
	Delete(id uuid.UUID) error
	CountAll(ctx context.Context) (int64, error)
	CountBetween(ctx context.Context, from time.Time, to time.Time) (int64, error)
	CountVerified(ctx context.Context) (int64, error)
}

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) Create(u entity.User) (entity.User, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return u, err
	}
	tx.Commit()
	return u, nil
}

func (r *UserRepository) FindAll(qp *util.QueryParams) ([]entity.User, int, error) {
	var entities []entity.User
	var totalCount int64

	query := r.Db.Model(&entity.User{})
	query = util.ApplySearch(query, qp)
	query = entity.User{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, userSortColumns, "users.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (r *UserRepository) FindById(id uuid.UUID) (entity.User, error) {
	var u entity.User
	if err := r.Db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (r *UserRepository) Update(u entity.User) error {
	tx := r.Db.Begin()
	if err := tx.Model(&u).Updates(u).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
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
