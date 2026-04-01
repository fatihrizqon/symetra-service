package repository

import (
	"errors"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Sort column whitelists ───────────────────────────────────────────────────

var fiscalYearSortColumns = map[string]string{
	"name":       "fiscal_years.name",
	"start_date": "fiscal_years.start_date",
	"end_date":   "fiscal_years.end_date",
	"status":     "fiscal_years.status",
	"created_at": "fiscal_years.created_at",
}

var fiscalPeriodSortColumns = map[string]string{
	"name":          "fiscal_periods.name",
	"period_number": "fiscal_periods.period_number",
	"start_date":    "fiscal_periods.start_date",
	"end_date":      "fiscal_periods.end_date",
	"status":        "fiscal_periods.status",
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

type IFiscalYearRepository interface {
	Create(fy entity.FiscalYear) (entity.FiscalYear, error)
	FindAll(qp *util.QueryParams) ([]entity.FiscalYear, int, error)
	FindById(id uuid.UUID) (entity.FiscalYear, error)
	FindActive() (entity.FiscalYear, error)
	Update(fy entity.FiscalYear) error
	Delete(id uuid.UUID) error
}

type IFiscalPeriodRepository interface {
	CreateBatch(periods []entity.FiscalPeriod) error
	FindAll(qp *util.QueryParams) ([]entity.FiscalPeriod, int, error)
	FindById(id uuid.UUID) (entity.FiscalPeriod, error)
	FindByDate(date time.Time) (entity.FiscalPeriod, error)
	FindCurrentOpen() (entity.FiscalPeriod, error)
	Update(p entity.FiscalPeriod) error
	CreateLog(log entity.FiscalPeriodLog) error
	FindLogs(periodId uuid.UUID) ([]entity.FiscalPeriodLog, error)
}

// ─── FiscalYear Repository ────────────────────────────────────────────────────

type FiscalYearRepository struct {
	Db *gorm.DB
}

func NewFiscalYearRepository(db *gorm.DB) IFiscalYearRepository {
	return &FiscalYearRepository{Db: db}
}

func (r *FiscalYearRepository) Create(fy entity.FiscalYear) (entity.FiscalYear, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&fy).Error; err != nil {
		tx.Rollback()
		return fy, err
	}
	tx.Commit()
	return fy, nil
}

func (r *FiscalYearRepository) FindAll(qp *util.QueryParams) ([]entity.FiscalYear, int, error) {
	var entities []entity.FiscalYear
	var totalCount int64

	query := r.Db.Model(&entity.FiscalYear{})
	query = util.ApplySearch(query, qp)
	query = entity.FiscalYear{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, fiscalYearSortColumns, "fiscal_years.start_date")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *FiscalYearRepository) FindById(id uuid.UUID) (entity.FiscalYear, error) {
	var fy entity.FiscalYear
	if err := r.Db.Preload("Periods").Where("id = ?", id).First(&fy).Error; err != nil {
		return fy, errors.New("fiscal year not found")
	}
	return fy, nil
}

func (r *FiscalYearRepository) FindActive() (entity.FiscalYear, error) {
	var fy entity.FiscalYear
	if err := r.Db.Where("status = ?", entity.FiscalYearActive).First(&fy).Error; err != nil {
		return fy, errors.New("no active fiscal year found")
	}
	return fy, nil
}

func (r *FiscalYearRepository) Update(fy entity.FiscalYear) error {
	tx := r.Db.Begin()
	if err := tx.Model(&fy).Updates(fy).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *FiscalYearRepository) Delete(id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.FiscalYear{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// ─── FiscalPeriod Repository ──────────────────────────────────────────────────

type FiscalPeriodRepository struct {
	Db *gorm.DB
}

func NewFiscalPeriodRepository(db *gorm.DB) IFiscalPeriodRepository {
	return &FiscalPeriodRepository{Db: db}
}

func (r *FiscalPeriodRepository) CreateBatch(periods []entity.FiscalPeriod) error {
	tx := r.Db.Begin()
	if err := tx.Create(&periods).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *FiscalPeriodRepository) FindAll(qp *util.QueryParams) ([]entity.FiscalPeriod, int, error) {
	var entities []entity.FiscalPeriod
	var totalCount int64

	query := r.Db.Model(&entity.FiscalPeriod{}).Preload("FiscalYear")
	query = util.ApplySearch(query, qp)
	query = entity.FiscalPeriod{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, fiscalPeriodSortColumns, "fiscal_periods.period_number")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *FiscalPeriodRepository) FindById(id uuid.UUID) (entity.FiscalPeriod, error) {
	// BUG NOTE 06032026: sempat terjadi saat close fiscal period, uuid not found. tapi sesaat itu bisa.
	var p entity.FiscalPeriod
	if err := r.Db.Preload("FiscalYear").Where("id = ?", id).First(&p).Error; err != nil {
		return p, errors.New("fiscal period not found")
	}
	return p, nil
}

// FindByDate finds the open or closed (but not locked) period that covers the given date.
// Used by transaction guards to validate whether a date is in an editable period.
func (r *FiscalPeriodRepository) FindByDate(date time.Time) (entity.FiscalPeriod, error) {
	var p entity.FiscalPeriod
	err := r.Db.
		Where("start_date <= ? AND end_date >= ? AND status != ?",
			date, date, entity.FiscalPeriodLocked).
		First(&p).Error
	if err != nil {
		return p, errors.New("no editable fiscal period found for this date")
	}
	return p, nil
}

func (r *FiscalPeriodRepository) FindCurrentOpen() (entity.FiscalPeriod, error) {
	var p entity.FiscalPeriod
	now := time.Now()
	err := r.Db.
		Where("start_date <= ? AND end_date >= ? AND status = ?",
			now, now, entity.FiscalPeriodOpen).
		First(&p).Error
	if err != nil {
		return p, errors.New("no open fiscal period found for today")
	}
	return p, nil
}

func (r *FiscalPeriodRepository) Update(p entity.FiscalPeriod) error {
	tx := r.Db.Begin()
	if err := tx.Model(&p).Updates(map[string]interface{}{
		"status":    p.Status,
		"closed_at": p.ClosedAt,
		"closed_by": p.ClosedBy,
		"locked_at": p.LockedAt,
		"locked_by": p.LockedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *FiscalPeriodRepository) CreateLog(log entity.FiscalPeriodLog) error {
	return r.Db.Create(&log).Error
}

func (r *FiscalPeriodRepository) FindLogs(periodId uuid.UUID) ([]entity.FiscalPeriodLog, error) {
	var logs []entity.FiscalPeriodLog
	if err := r.Db.Where("fiscal_period_id = ?", periodId).
		Order("performed_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
