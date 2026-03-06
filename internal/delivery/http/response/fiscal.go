package response

import (
	"time"

	"github.com/google/uuid"
)

// ─── Fiscal Year ──────────────────────────────────────────────────────────────

type FiscalYearResponse struct {
	Id         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	StartDate  time.Time              `json:"start_date"`
	EndDate    time.Time              `json:"end_date"`
	PeriodType string                 `json:"period_type"`
	Status     string                 `json:"status"`
	ClosedAt   *time.Time             `json:"closed_at,omitempty"`
	ClosedBy   *uuid.UUID             `json:"closed_by,omitempty"`
	CreatedBy  uuid.UUID              `json:"created_by"`
	Periods    []FiscalPeriodResponse `json:"periods,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// ─── Fiscal Period ────────────────────────────────────────────────────────────

type FiscalPeriodResponse struct {
	Id           uuid.UUID  `json:"id"`
	FiscalYearId uuid.UUID  `json:"fiscal_year_id"`
	FiscalYear   string     `json:"fiscal_year_name,omitempty"`
	Name         string     `json:"name"`
	PeriodNumber int        `json:"period_number"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      time.Time  `json:"end_date"`
	IsStub       bool       `json:"is_stub"`
	Status       string     `json:"status"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	ClosedBy     *uuid.UUID `json:"closed_by,omitempty"`
	LockedAt     *time.Time `json:"locked_at,omitempty"`
	LockedBy     *uuid.UUID `json:"locked_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ─── Fiscal Period Log ────────────────────────────────────────────────────────

type FiscalPeriodLogResponse struct {
	Id             uuid.UUID `json:"id"`
	FiscalPeriodId uuid.UUID `json:"fiscal_period_id"`
	Action         string    `json:"action"`
	FromStatus     string    `json:"from_status"`
	ToStatus       string    `json:"to_status"`
	Reason         string    `json:"reason,omitempty"`
	PerformedBy    uuid.UUID `json:"performed_by"`
	PerformedAt    time.Time `json:"performed_at"`
}
