package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ─── Interface ────────────────────────────────────────────────────────────────

type IFiscalService interface {
	// Fiscal Year
	CreateFiscalYear(req request.FiscalYearCreateRequest, createdBy uuid.UUID) (entity.FiscalYear, error)
	FindAllFiscalYears(qp *util.QueryParams) ([]response.FiscalYearResponse, int, error)
	FindFiscalYearById(id uuid.UUID) (response.FiscalYearResponse, error)
	UpdateFiscalYear(req request.FiscalYearUpdateRequest) (entity.FiscalYear, error)
	ActivateFiscalYear(id uuid.UUID) (entity.FiscalYear, error)
	DeleteFiscalYear(id uuid.UUID) error

	// Fiscal Period
	FindAllPeriods(qp *util.QueryParams) ([]response.FiscalPeriodResponse, int, error)
	FindPeriodById(id uuid.UUID) (response.FiscalPeriodResponse, error)
	ClosePeriod(id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	ReopenPeriod(req request.FiscalPeriodReopenRequest, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	LockPeriod(id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	FindPeriodLogs(id uuid.UUID) ([]response.FiscalPeriodLogResponse, error)

	// Guard — used by other services (journal entries, revenues, expenses, etc.)
	ValidatePeriodOpen(date time.Time) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type FiscalService struct {
	fyRepo   repository.IFiscalYearRepository
	fpRepo   repository.IFiscalPeriodRepository
	validate *validator.Validate
}

func NewFiscalService(
	fyRepo repository.IFiscalYearRepository,
	fpRepo repository.IFiscalPeriodRepository,
	validate *validator.Validate,
) IFiscalService {
	return &FiscalService{fyRepo: fyRepo, fpRepo: fpRepo, validate: validate}
}

// ── Fiscal Year ───────────────────────────────────────────────────────────────

func (s *FiscalService) CreateFiscalYear(req request.FiscalYearCreateRequest, createdBy uuid.UUID) (entity.FiscalYear, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.FiscalYear{}, err
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return entity.FiscalYear{}, errors.New("invalid start_date format, use YYYY-MM-DD")
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return entity.FiscalYear{}, errors.New("invalid end_date format, use YYYY-MM-DD")
	}
	if !endDate.After(startDate) {
		return entity.FiscalYear{}, errors.New("end_date must be after start_date")
	}

	fy := entity.FiscalYear{
		Name:       req.Name,
		StartDate:  startDate,
		EndDate:    endDate,
		PeriodType: entity.PeriodType(req.PeriodType),
		Status:     entity.FiscalYearDraft,
		CreatedBy:  createdBy,
	}

	created, err := s.fyRepo.Create(fy)
	if err != nil {
		return created, err
	}

	// Auto-generate periods (unless custom — those are handled separately)
	if req.PeriodType != string(entity.PeriodTypeCustom) {
		actualStart := startDate
		if req.StubStartDate != "" {
			stubStart, err := parseDate(req.StubStartDate)
			if err == nil && stubStart.After(startDate) && stubStart.Before(endDate) {
				actualStart = stubStart
			}
		}
		periods := generatePeriods(created.Id, startDate, actualStart, endDate, entity.PeriodType(req.PeriodType))
		if err := s.fpRepo.CreateBatch(periods); err != nil {
			return created, fmt.Errorf("fiscal year created but periods failed: %w", err)
		}
	}

	return s.fyRepo.FindById(created.Id)
}

func (s *FiscalService) FindAllFiscalYears(qp *util.QueryParams) ([]response.FiscalYearResponse, int, error) {
	entities, totalCount, err := s.fyRepo.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.FiscalYearResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.FiscalYearResponse, 0, len(entities))
	for _, fy := range entities {
		resps = append(resps, mapFiscalYear(fy))
	}
	return resps, totalCount, nil
}

func (s *FiscalService) FindFiscalYearById(id uuid.UUID) (response.FiscalYearResponse, error) {
	fy, err := s.fyRepo.FindById(id)
	if err != nil {
		return response.FiscalYearResponse{}, err
	}
	res := mapFiscalYear(fy)
	periods := make([]response.FiscalPeriodResponse, 0, len(fy.Periods))
	for _, p := range fy.Periods {
		periods = append(periods, mapFiscalPeriod(p))
	}
	res.Periods = periods
	return res, nil
}

func (s *FiscalService) UpdateFiscalYear(req request.FiscalYearUpdateRequest) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(req.Id)
	if err != nil {
		return fy, err
	}
	if fy.Status != entity.FiscalYearDraft {
		return fy, errors.New("only draft fiscal years can be renamed")
	}
	fy.Name = req.Name
	return fy, s.fyRepo.Update(fy)
}

// ActivateFiscalYear sets the fiscal year status to active.
// Business rule: only one fiscal year can be active at a time.
func (s *FiscalService) ActivateFiscalYear(id uuid.UUID) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(id)
	if err != nil {
		return fy, err
	}
	if fy.Status != entity.FiscalYearDraft {
		return fy, errors.New("only draft fiscal years can be activated")
	}

	// Ensure no other active fiscal year exists
	existing, err := s.fyRepo.FindActive()
	if err == nil && existing.Id != id {
		return fy, fmt.Errorf("fiscal year '%s' is already active — close it first", existing.Name)
	}

	fy.Status = entity.FiscalYearActive
	return fy, s.fyRepo.Update(fy)
}

// DeleteFiscalYear only allows deleting draft fiscal years with no transactions.
func (s *FiscalService) DeleteFiscalYear(id uuid.UUID) error {
	fy, err := s.fyRepo.FindById(id)
	if err != nil {
		return err
	}
	if fy.Status != entity.FiscalYearDraft {
		return errors.New("only draft fiscal years can be deleted")
	}
	return s.fyRepo.Delete(id)
}

// ── Fiscal Period ─────────────────────────────────────────────────────────────

func (s *FiscalService) FindAllPeriods(qp *util.QueryParams) ([]response.FiscalPeriodResponse, int, error) {
	entities, totalCount, err := s.fpRepo.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.FiscalPeriodResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.FiscalPeriodResponse, 0, len(entities))
	for _, p := range entities {
		resps = append(resps, mapFiscalPeriod(p))
	}
	return resps, totalCount, nil
}

func (s *FiscalService) FindPeriodById(id uuid.UUID) (response.FiscalPeriodResponse, error) {
	p, err := s.fpRepo.FindById(id)
	if err != nil {
		return response.FiscalPeriodResponse{}, err
	}
	return mapFiscalPeriod(p), nil
}

// ClosePeriod transitions a period from OPEN → CLOSED.
func (s *FiscalService) ClosePeriod(id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	p, err := s.fpRepo.FindById(id)
	if err != nil {
		return p, err
	}
	if p.Status != entity.FiscalPeriodOpen {
		return p, fmt.Errorf("period '%s' is not open (current status: %s)", p.Name, p.Status)
	}

	now := time.Now()
	p.Status = entity.FiscalPeriodClosed
	p.ClosedAt = &now
	p.ClosedBy = &performedBy

	if err := s.fpRepo.Update(p); err != nil {
		return p, err
	}

	_ = s.fpRepo.CreateLog(entity.FiscalPeriodLog{
		FiscalPeriodId: id,
		Action:         "CLOSED",
		FromStatus:     string(entity.FiscalPeriodOpen),
		ToStatus:       string(entity.FiscalPeriodClosed),
		PerformedBy:    performedBy,
	})

	return p, nil
}

// ReopenPeriod transitions CLOSED → OPEN. Reason is required for audit.
// Locked periods cannot be reopened — ever.
func (s *FiscalService) ReopenPeriod(req request.FiscalPeriodReopenRequest, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.FiscalPeriod{}, err
	}

	p, err := s.fpRepo.FindById(req.Id)
	if err != nil {
		return p, err
	}
	if p.Status == entity.FiscalPeriodLocked {
		return p, fmt.Errorf("period '%s' is permanently locked and cannot be reopened", p.Name)
	}
	if p.Status != entity.FiscalPeriodClosed {
		return p, fmt.Errorf("period '%s' is not closed (current status: %s)", p.Name, p.Status)
	}

	p.Status = entity.FiscalPeriodOpen
	p.ClosedAt = nil
	p.ClosedBy = nil

	if err := s.fpRepo.Update(p); err != nil {
		return p, err
	}

	_ = s.fpRepo.CreateLog(entity.FiscalPeriodLog{
		FiscalPeriodId: req.Id,
		Action:         "REOPENED",
		FromStatus:     string(entity.FiscalPeriodClosed),
		ToStatus:       string(entity.FiscalPeriodOpen),
		Reason:         req.Reason,
		PerformedBy:    performedBy,
	})

	return p, nil
}

// LockPeriod transitions CLOSED → LOCKED. This is permanent.
func (s *FiscalService) LockPeriod(id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	p, err := s.fpRepo.FindById(id)
	if err != nil {
		return p, err
	}
	if p.Status != entity.FiscalPeriodClosed {
		return p, fmt.Errorf("period '%s' must be closed before locking (current status: %s)", p.Name, p.Status)
	}

	now := time.Now()
	p.Status = entity.FiscalPeriodLocked
	p.LockedAt = &now
	p.LockedBy = &performedBy

	if err := s.fpRepo.Update(p); err != nil {
		return p, err
	}

	_ = s.fpRepo.CreateLog(entity.FiscalPeriodLog{
		FiscalPeriodId: id,
		Action:         "LOCKED",
		FromStatus:     string(entity.FiscalPeriodClosed),
		ToStatus:       string(entity.FiscalPeriodLocked),
		PerformedBy:    performedBy,
	})

	return p, nil
}

func (s *FiscalService) FindPeriodLogs(id uuid.UUID) ([]response.FiscalPeriodLogResponse, error) {
	logs, err := s.fpRepo.FindLogs(id)
	if err != nil {
		return nil, err
	}
	resps := make([]response.FiscalPeriodLogResponse, 0, len(logs))
	for _, l := range logs {
		resps = append(resps, response.FiscalPeriodLogResponse{
			Id:             l.Id,
			FiscalPeriodId: l.FiscalPeriodId,
			Action:         l.Action,
			FromStatus:     l.FromStatus,
			ToStatus:       l.ToStatus,
			Reason:         l.Reason,
			PerformedBy:    l.PerformedBy,
			PerformedAt:    l.PerformedAt,
		})
	}
	return resps, nil
}

// ── Guard ─────────────────────────────────────────────────────────────────────

// ValidatePeriodOpen is called by other services (journal entries, revenues, expenses)
// before creating or editing a transaction. It rejects the operation if the period
// covering the given date is closed or locked.
func (s *FiscalService) ValidatePeriodOpen(date time.Time) error {
	p, err := s.fpRepo.FindByDate(date)
	if err != nil {
		// No period at all for this date — might be before fiscal year starts
		return fmt.Errorf("no fiscal period found for date %s — check your fiscal year setup", date.Format("2006-01-02"))
	}
	switch p.Status {
	case entity.FiscalPeriodLocked:
		return fmt.Errorf("fiscal period '%s' is permanently locked — transactions cannot be modified", p.Name)
	case entity.FiscalPeriodClosed:
		return fmt.Errorf("fiscal period '%s' is closed — ask your administrator to reopen it", p.Name)
	case entity.FiscalPeriodOpen:
		return nil
	default:
		return fmt.Errorf("fiscal period '%s' has unknown status: %s", p.Name, p.Status)
	}
}

// generatePeriods auto-generates fiscal periods for monthly or quarterly types.
// If actualStart > startDate, the first period is a stub period.
func generatePeriods(fyId uuid.UUID, fyStart, actualStart, fyEnd time.Time, pType entity.PeriodType) []entity.FiscalPeriod {
	var periods []entity.FiscalPeriod
	periodNumber := 1

	// If there is a stub (onboarding mid-period)
	if actualStart.After(fyStart) {
		stubEnd := firstDayOfNextMonth(fyStart).Add(-time.Second)
		if stubEnd.After(fyEnd) {
			stubEnd = fyEnd
		}
		periods = append(periods, entity.FiscalPeriod{
			FiscalYearId: fyId,
			Name:         fmt.Sprintf("Stub – %s", fyStart.Format("Jan 2006")),
			PeriodNumber: periodNumber,
			StartDate:    actualStart,
			EndDate:      stubEnd,
			IsStub:       true,
			Status:       entity.FiscalPeriodOpen,
		})
		periodNumber++
		fyStart = firstDayOfNextMonth(fyStart)
	}

	cursor := fyStart
	for cursor.Before(fyEnd) || cursor.Equal(fyStart) {
		var periodEnd time.Time
		switch pType {
		case entity.PeriodTypeQuarterly:
			periodEnd = cursor.AddDate(0, 3, 0).Add(-time.Second)
		default: // monthly
			periodEnd = firstDayOfNextMonth(cursor).Add(-time.Second)
		}
		if periodEnd.After(fyEnd) {
			periodEnd = fyEnd
		}

		periods = append(periods, entity.FiscalPeriod{
			FiscalYearId: fyId,
			Name:         periodName(cursor, pType),
			PeriodNumber: periodNumber,
			StartDate:    cursor,
			EndDate:      periodEnd,
			IsStub:       false,
			Status:       entity.FiscalPeriodOpen,
		})

		periodNumber++
		switch pType {
		case entity.PeriodTypeQuarterly:
			cursor = cursor.AddDate(0, 3, 0)
		default:
			cursor = firstDayOfNextMonth(cursor)
		}

		if !cursor.Before(fyEnd) && !cursor.Equal(fyEnd) {
			break
		}
	}

	return periods
}

func firstDayOfNextMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}

func periodName(t time.Time, pType entity.PeriodType) string {
	switch pType {
	case entity.PeriodTypeQuarterly:
		q := (int(t.Month())-1)/3 + 1
		return fmt.Sprintf("Q%d %d", q, t.Year())
	default:
		return t.Format("January 2006")
	}
}

// ── Mappers ───────────────────────────────────────────────────────────────────

func mapFiscalYear(fy entity.FiscalYear) response.FiscalYearResponse {
	return response.FiscalYearResponse{
		Id:         fy.Id,
		Name:       fy.Name,
		StartDate:  fy.StartDate,
		EndDate:    fy.EndDate,
		PeriodType: string(fy.PeriodType),
		Status:     string(fy.Status),
		ClosedAt:   fy.ClosedAt,
		ClosedBy:   fy.ClosedBy,
		CreatedBy:  fy.CreatedBy,
		CreatedAt:  fy.CreatedAt,
		UpdatedAt:  fy.UpdatedAt,
	}
}

func mapFiscalPeriod(p entity.FiscalPeriod) response.FiscalPeriodResponse {
	fyName := ""
	if p.FiscalYear.Name != "" {
		fyName = p.FiscalYear.Name
	}
	return response.FiscalPeriodResponse{
		Id:           p.Id,
		FiscalYearId: p.FiscalYearId,
		FiscalYear:   fyName,
		Name:         p.Name,
		PeriodNumber: p.PeriodNumber,
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		IsStub:       p.IsStub,
		Status:       string(p.Status),
		ClosedAt:     p.ClosedAt,
		ClosedBy:     p.ClosedBy,
		LockedAt:     p.LockedAt,
		LockedBy:     p.LockedBy,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
