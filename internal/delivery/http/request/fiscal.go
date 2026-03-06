package request

import "github.com/google/uuid"

// ─── Fiscal Year ──────────────────────────────────────────────────────────────

type FiscalYearCreateRequest struct {
	Name        string `validate:"required,min=1,max=100" json:"name"`
	StartDate   string `validate:"required" json:"start_date"` // YYYY-MM-DD
	EndDate     string `validate:"required" json:"end_date"`   // YYYY-MM-DD
	PeriodType  string `validate:"required,oneof=monthly quarterly custom" json:"period_type"`
	// StubStartDate is optional. If provided, the first period will start from this
	// date instead of the fiscal year start_date. Used when onboarding mid-period.
	StubStartDate string `json:"stub_start_date,omitempty"` // YYYY-MM-DD
}

type FiscalYearUpdateRequest struct {
	Id   uuid.UUID
	Name string `validate:"required,min=1,max=100" json:"name"`
}

type FiscalYearActivateRequest struct {
	Id uuid.UUID
}

// ─── Fiscal Period ────────────────────────────────────────────────────────────

// FiscalPeriodCloseRequest closes a single period.
type FiscalPeriodCloseRequest struct {
	Id uuid.UUID
}

// FiscalPeriodReopenRequest reopens a closed (not locked) period.
// Reason is mandatory for audit trail.
type FiscalPeriodReopenRequest struct {
	Id     uuid.UUID
	Reason string `validate:"required,min=10" json:"reason"`
}

// FiscalPeriodLockRequest permanently locks a closed period.
type FiscalPeriodLockRequest struct {
	Id uuid.UUID
}

// FiscalCustomPeriodRequest is used when PeriodType = "custom" to define
// each period manually.
type FiscalCustomPeriodRequest struct {
	Name      string `validate:"required,min=1,max=100" json:"name"`
	StartDate string `validate:"required" json:"start_date"`
	EndDate   string `validate:"required" json:"end_date"`
}

// FiscalYearCreateCustomRequest is the full request when PeriodType = "custom".
type FiscalYearCreateCustomRequest struct {
	FiscalYearCreateRequest
	Periods []FiscalCustomPeriodRequest `validate:"required,min=1,dive" json:"periods"`
}
