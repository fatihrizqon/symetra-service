package repository

import (
	"errors"
	"fmt"
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
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error)
	FindById(companyID, id uuid.UUID) (entity.FiscalYear, error)
	FindActive(companyID uuid.UUID) (entity.FiscalYear, error)
	FindLastClosed(companyID uuid.UUID) (entity.FiscalYear, error)
	Update(fy entity.FiscalYear) error
	Delete(companyID, id uuid.UUID) error
	// Closing/Opening helpers
	CountOpenOrClosedPeriods(fyID uuid.UUID) (int64, error)
	CountDraftJournalEntries(companyID, fyID uuid.UUID) (int64, error)
	HasClosingJE(fyID uuid.UUID) (bool, error)
	HasOpeningJE(fyID uuid.UUID) (bool, error)
}

type IFiscalPeriodRepository interface {
	CreateBatch(periods []entity.FiscalPeriod) error
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalPeriod, int, error)
	FindById(companyID, id uuid.UUID) (entity.FiscalPeriod, error)
	FindByDate(companyID uuid.UUID, date time.Time) (entity.FiscalPeriod, error)
	FindCurrentOpen(companyID uuid.UUID) (entity.FiscalPeriod, error)
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
	fmt.Println("COMPANY ID", fy.CompanyId)
	tx := r.Db.Begin()
	if err := tx.Create(&fy).Error; err != nil {
		tx.Rollback()
		return fy, err
	}
	tx.Commit()
	return r.FindById(fy.CompanyId, fy.Id)
}

func (r *FiscalYearRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error) {
	var entities []entity.FiscalYear
	var totalCount int64

	query := r.Db.Model(&entity.FiscalYear{}).Where("fiscal_years.company_id = ?", companyID)
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

func (r *FiscalYearRepository) FindById(companyID, id uuid.UUID) (entity.FiscalYear, error) {
	var fy entity.FiscalYear
	if err := r.Db.Preload("Periods").
		Where("id = ? AND company_id = ?", id, companyID).
		First(&fy).Error; err != nil {
		return fy, errors.New("fiscal year not found")
	}
	return fy, nil
}

func (r *FiscalYearRepository) FindActive(companyID uuid.UUID) (entity.FiscalYear, error) {
	var fy entity.FiscalYear
	if err := r.Db.Where("company_id = ? AND status = ?", companyID, entity.FiscalYearActive).
		First(&fy).Error; err != nil {
		return fy, errors.New("no active fiscal year found")
	}
	return fy, nil
}

func (r *FiscalYearRepository) Update(fy entity.FiscalYear) error {
	tx := r.Db.Begin()
	// Use Save() instead of Updates() so that nil pointer fields (e.g. OpeningJEId = nil
	// when deleting an opening balance) are also persisted to the database.
	// Updates() skips zero-value / nil fields by design in GORM.
	if err := tx.Save(&fy).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *FiscalYearRepository) Delete(companyID, id uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ? AND company_id = ?", id, companyID).
		Delete(&entity.FiscalYear{}).Error; err != nil {
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

func (r *FiscalPeriodRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalPeriod, int, error) {
	var entities []entity.FiscalPeriod
	var totalCount int64

	query := r.Db.Model(&entity.FiscalPeriod{}).Preload("FiscalYear").
		Joins("INNER JOIN fiscal_years fy ON fy.id = fiscal_periods.fiscal_year_id").
		Where("fy.company_id = ?", companyID)
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

func (r *FiscalPeriodRepository) FindById(companyID, id uuid.UUID) (entity.FiscalPeriod, error) {
	var p entity.FiscalPeriod
	err := r.Db.Preload("FiscalYear").
		Joins("INNER JOIN fiscal_years fy ON fy.id = fiscal_periods.fiscal_year_id").
		Where("fiscal_periods.id = ? AND fy.company_id = ?", id, companyID).
		First(&p).Error
	if err != nil {
		return p, fmt.Errorf("fiscal period not found")
	}
	return p, nil
}

// FindByDate finds the open or closed (but not locked) period covering the given date
// for the specified company.
func (r *FiscalPeriodRepository) FindByDate(companyID uuid.UUID, date time.Time) (entity.FiscalPeriod, error) {
	var p entity.FiscalPeriod
	err := r.Db.
		Joins("INNER JOIN fiscal_years fy ON fy.id = fiscal_periods.fiscal_year_id").
		Where("fy.company_id = ? AND fiscal_periods.start_date <= ? AND fiscal_periods.end_date >= ? AND fiscal_periods.status != ?",
			companyID, date, date, entity.FiscalPeriodLocked).
		First(&p).Error
	if err != nil {
		return p, errors.New("no editable fiscal period found for this date")
	}
	return p, nil
}

func (r *FiscalPeriodRepository) FindCurrentOpen(companyID uuid.UUID) (entity.FiscalPeriod, error) {
	var p entity.FiscalPeriod
	now := time.Now()
	err := r.Db.
		Joins("INNER JOIN fiscal_years fy ON fy.id = fiscal_periods.fiscal_year_id").
		Where("fy.company_id = ? AND fiscal_periods.start_date <= ? AND fiscal_periods.end_date >= ? AND fiscal_periods.status = ?",
			companyID, now, now, entity.FiscalPeriodOpen).
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

// ─── Additional FiscalYear Repository Methods ─────────────────────────────────

func (r *FiscalYearRepository) FindLastClosed(companyID uuid.UUID) (entity.FiscalYear, error) {
	var fy entity.FiscalYear
	err := r.Db.Where("company_id = ? AND status = ?", companyID, entity.FiscalYearClosed).
		Order("end_date DESC").First(&fy).Error
	if err != nil {
		return fy, fmt.Errorf("no closed fiscal year found")
	}
	return fy, nil
}

// CountOpenOrClosedPeriods returns number of periods NOT yet locked in a FY.
// Used to validate readiness for closing_review state.
func (r *FiscalYearRepository) CountOpenOrClosedPeriods(fyID uuid.UUID) (int64, error) {
	var count int64
	err := r.Db.Model(&entity.FiscalPeriod{}).
		Where("fiscal_year_id = ? AND status IN ?", fyID, []string{"open", "closed"}).
		Count(&count).Error
	return count, err
}

// CountDraftJournalEntries returns number of draft JEs in the FY date range.
// A FY cannot be closed while drafts exist.
func (r *FiscalYearRepository) CountDraftJournalEntries(companyID, fyID uuid.UUID) (int64, error) {
	var count int64
	var fy entity.FiscalYear
	if err := r.Db.Where("id = ? AND company_id = ?", fyID, companyID).First(&fy).Error; err != nil {
		return 0, err
	}
	err := r.Db.Model(&entity.JournalEntry{}).
		Where("company_id = ? AND status = ? AND date >= ? AND date <= ?",
			companyID, "draft", fy.StartDate, fy.EndDate).
		Count(&count).Error
	return count, err
}

// HasClosingJE checks if a closing journal entry already exists for this FY.
func (r *FiscalYearRepository) HasClosingJE(fyID uuid.UUID) (bool, error) {
	var fy entity.FiscalYear
	if err := r.Db.Where("id = ?", fyID).First(&fy).Error; err != nil {
		return false, err
	}
	return fy.ClosingJEId != nil, nil
}

// HasOpeningJE checks if an opening balance JE already exists for this FY.
func (r *FiscalYearRepository) HasOpeningJE(fyID uuid.UUID) (bool, error) {
	var fy entity.FiscalYear
	if err := r.Db.Where("id = ?", fyID).First(&fy).Error; err != nil {
		return false, err
	}
	return fy.OpeningJEId != nil, nil
}
