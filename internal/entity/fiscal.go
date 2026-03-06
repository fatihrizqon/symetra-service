package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Fiscal Year ──────────────────────────────────────────────────────────────

func (FiscalYear) TableName() string { return "fiscal_years" }

// FiscalYearStatus represents the lifecycle state of a fiscal year.
type FiscalYearStatus string

const (
	FiscalYearDraft  FiscalYearStatus = "draft"
	FiscalYearActive FiscalYearStatus = "active"
	FiscalYearClosed FiscalYearStatus = "closed"
)

// PeriodType defines how periods are auto-generated within a fiscal year.
type PeriodType string

const (
	PeriodTypeMonthly   PeriodType = "monthly"   // 12 periods
	PeriodTypeQuarterly PeriodType = "quarterly"  // 4 periods
	PeriodTypeCustom    PeriodType = "custom"     // user-defined
)

type FiscalYear struct {
	Id          uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	Name        string           `gorm:"type:character varying(100);not null;" json:"name"`
	StartDate   time.Time        `gorm:"type:date;not null;" json:"start_date"`
	EndDate     time.Time        `gorm:"type:date;not null;" json:"end_date"`
	PeriodType  PeriodType       `gorm:"type:character varying(20);not null;default:'monthly';" json:"period_type"`
	Status      FiscalYearStatus `gorm:"type:character varying(20);not null;default:'draft';" json:"status"`
	ClosedAt    *time.Time       `gorm:"default:null;" json:"closed_at,omitempty"`
	ClosedBy    *uuid.UUID       `gorm:"type:uuid;default:null;" json:"closed_by,omitempty"`
	CreatedBy   uuid.UUID        `gorm:"type:uuid;not null;" json:"created_by"`
	Periods     []FiscalPeriod   `gorm:"foreignKey:FiscalYearId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"periods,omitempty"`
	CreatedAt   time.Time        `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (FiscalYear) SearchableFields() []string {
	return []string{"name"}
}

func (FiscalYear) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("fiscal_years.status IN ?", values)
	}
	if values, ok := filters["period_type"]; ok {
		db = db.Where("fiscal_years.period_type IN ?", values)
	}
	return db
}

// ─── Fiscal Period ────────────────────────────────────────────────────────────

func (FiscalPeriod) TableName() string { return "fiscal_periods" }

// FiscalPeriodStatus represents the lifecycle state of a single period.
// Transitions: open → closed → locked
// Reopen: closed → open (with reason, admin only)
// Locked is permanent — no reopen.
type FiscalPeriodStatus string

const (
	FiscalPeriodOpen   FiscalPeriodStatus = "open"
	FiscalPeriodClosed FiscalPeriodStatus = "closed"
	FiscalPeriodLocked FiscalPeriodStatus = "locked"
)

type FiscalPeriod struct {
	Id           uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	FiscalYearId uuid.UUID          `gorm:"type:uuid;not null;index;" json:"fiscal_year_id"`
	FiscalYear   FiscalYear         `gorm:"foreignKey:FiscalYearId;" json:"fiscal_year,omitempty"`
	Name         string             `gorm:"type:character varying(100);not null;" json:"name"`
	PeriodNumber int                `gorm:"type:int;not null;" json:"period_number"`
	StartDate    time.Time          `gorm:"type:date;not null;" json:"start_date"`
	EndDate      time.Time          `gorm:"type:date;not null;" json:"end_date"`
	IsStub       bool               `gorm:"type:boolean;not null;default:false;" json:"is_stub"`
	Status       FiscalPeriodStatus `gorm:"type:character varying(20);not null;default:'open';" json:"status"`
	ClosedAt     *time.Time         `gorm:"default:null;" json:"closed_at,omitempty"`
	ClosedBy     *uuid.UUID         `gorm:"type:uuid;default:null;" json:"closed_by,omitempty"`
	LockedAt     *time.Time         `gorm:"default:null;" json:"locked_at,omitempty"`
	LockedBy     *uuid.UUID         `gorm:"type:uuid;default:null;" json:"locked_by,omitempty"`
	CreatedAt    time.Time          `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (FiscalPeriod) SearchableFields() []string {
	return []string{"name"}
}

func (FiscalPeriod) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("fiscal_periods.status IN ?", values)
	}
	if values, ok := filters["fiscal_year_id"]; ok {
		db = db.Where("fiscal_periods.fiscal_year_id IN ?", values)
	}
	return db
}

// ─── Fiscal Period Log ────────────────────────────────────────────────────────

func (FiscalPeriodLog) TableName() string { return "fiscal_period_logs" }

type FiscalPeriodLog struct {
	Id             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	FiscalPeriodId uuid.UUID  `gorm:"type:uuid;not null;index;" json:"fiscal_period_id"`
	Action         string     `gorm:"type:character varying(50);not null;" json:"action"` // CLOSED, REOPENED, LOCKED
	FromStatus     string     `gorm:"type:character varying(20);not null;" json:"from_status"`
	ToStatus       string     `gorm:"type:character varying(20);not null;" json:"to_status"`
	Reason         string     `gorm:"type:text;" json:"reason"`
	PerformedBy    uuid.UUID  `gorm:"type:uuid;not null;" json:"performed_by"`
	PerformedAt    time.Time  `gorm:"autoCreateTime;" json:"performed_at"`
}
