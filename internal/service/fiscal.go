package service

import (
	"errors"
	"fmt"
	"strings"
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
	// Fiscal Year CRUD
	CreateFiscalYear(companyID uuid.UUID, req request.FiscalYearCreateRequest, createdBy uuid.UUID) (entity.FiscalYear, error)
	FindAllFiscalYears(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalYearResponse, int, error)
	FindFiscalYearById(companyID, id uuid.UUID) (response.FiscalYearResponse, error)
	UpdateFiscalYear(companyID uuid.UUID, req request.FiscalYearUpdateRequest) (entity.FiscalYear, error)
	ActivateFiscalYear(companyID, id uuid.UUID) (entity.FiscalYear, error)
	DeleteFiscalYear(companyID, id uuid.UUID) error

	// Fiscal Period
	FindAllPeriods(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalPeriodResponse, int, error)
	FindPeriodById(companyID, id uuid.UUID) (response.FiscalPeriodResponse, error)
	ClosePeriod(companyID, id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	ReopenPeriod(companyID uuid.UUID, req request.FiscalPeriodReopenRequest, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	LockPeriod(companyID, id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error)
	FindPeriodLogs(id uuid.UUID) ([]response.FiscalPeriodLogResponse, error)

	// Guard — used by other services
	ValidatePeriodOpen(companyID uuid.UUID, date time.Time) error

	// Year-End Closing
	ReadyToClose(companyID, fyID uuid.UUID) (response.FiscalReadinessResponse, error)
	CloseFiscalYear(companyID, fyID uuid.UUID, req request.FiscalYearCloseRequest, performedBy uuid.UUID) (entity.FiscalYear, error)

	// Opening Balance
	GenerateOpeningBalance(companyID, fyID uuid.UUID, performedBy uuid.UUID) (entity.FiscalYear, error)
	DeleteOpeningBalance(companyID, fyID uuid.UUID) (entity.FiscalYear, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type FiscalService struct {
	fyRepo       repository.IFiscalYearRepository
	fpRepo       repository.IFiscalPeriodRepository
	jeRepo       repository.IJournalEntryRepository
	coaRepo      repository.ICOARepository
	coaGroupRepo repository.ICOAGroupRepository
	reportRepo   repository.IReportRepository
	cfgRepo      repository.ICompanyConfigurationRepository
	validate     *validator.Validate
}

func NewFiscalService(
	fyRepo repository.IFiscalYearRepository,
	fpRepo repository.IFiscalPeriodRepository,
	jeRepo repository.IJournalEntryRepository,
	coaRepo repository.ICOARepository,
	coaGroupRepo repository.ICOAGroupRepository,
	reportRepo repository.IReportRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	validate *validator.Validate,
) IFiscalService {
	return &FiscalService{
		fyRepo:       fyRepo,
		fpRepo:       fpRepo,
		jeRepo:       jeRepo,
		coaRepo:      coaRepo,
		coaGroupRepo: coaGroupRepo,
		reportRepo:   reportRepo,
		cfgRepo:      cfgRepo,
		validate:     validate,
	}
}

// ── Fiscal Year CRUD ──────────────────────────────────────────────────────────

func (s *FiscalService) CreateFiscalYear(companyID uuid.UUID, req request.FiscalYearCreateRequest, createdBy uuid.UUID) (entity.FiscalYear, error) {
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
		CompanyId:  companyID,
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

	return s.fyRepo.FindById(companyID, created.Id)
}

func (s *FiscalService) FindAllFiscalYears(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalYearResponse, int, error) {
	entities, totalCount, err := s.fyRepo.FindAll(companyID, qp)
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

func (s *FiscalService) FindFiscalYearById(companyID, id uuid.UUID) (response.FiscalYearResponse, error) {
	fy, err := s.fyRepo.FindById(companyID, id)
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

func (s *FiscalService) UpdateFiscalYear(companyID uuid.UUID, req request.FiscalYearUpdateRequest) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(companyID, req.Id)
	if err != nil {
		return fy, err
	}
	if fy.Status != entity.FiscalYearDraft {
		return fy, errors.New("only draft fiscal years can be renamed")
	}
	fy.Name = req.Name
	return fy, s.fyRepo.Update(fy)
}

func (s *FiscalService) ActivateFiscalYear(companyID, id uuid.UUID) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(companyID, id)
	if err != nil {
		return fy, err
	}
	if fy.Status != entity.FiscalYearDraft {
		return fy, errors.New("only draft fiscal years can be activated")
	}

	existing, err := s.fyRepo.FindActive(companyID)
	if err == nil && existing.Id != id {
		return fy, fmt.Errorf("fiscal year '%s' is already active — close it first", existing.Name)
	}

	fy.Status = entity.FiscalYearActive
	return fy, s.fyRepo.Update(fy)
}

func (s *FiscalService) DeleteFiscalYear(companyID, id uuid.UUID) error {
	fy, err := s.fyRepo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if fy.Status != entity.FiscalYearDraft {
		return errors.New("only draft fiscal years can be deleted")
	}
	return s.fyRepo.Delete(companyID, id)
}

// ── Fiscal Period ─────────────────────────────────────────────────────────────

func (s *FiscalService) FindAllPeriods(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalPeriodResponse, int, error) {
	entities, totalCount, err := s.fpRepo.FindAll(companyID, qp)
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

func (s *FiscalService) FindPeriodById(companyID, id uuid.UUID) (response.FiscalPeriodResponse, error) {
	p, err := s.fpRepo.FindById(companyID, id)
	if err != nil {
		return response.FiscalPeriodResponse{}, err
	}
	return mapFiscalPeriod(p), nil
}

func (s *FiscalService) ClosePeriod(companyID, id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	p, err := s.fpRepo.FindById(companyID, id)
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

func (s *FiscalService) ReopenPeriod(companyID uuid.UUID, req request.FiscalPeriodReopenRequest, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.FiscalPeriod{}, err
	}

	p, err := s.fpRepo.FindById(companyID, req.Id)
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

func (s *FiscalService) LockPeriod(companyID, id uuid.UUID, performedBy uuid.UUID) (entity.FiscalPeriod, error) {
	p, err := s.fpRepo.FindById(companyID, id)
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

func (s *FiscalService) ValidatePeriodOpen(companyID uuid.UUID, date time.Time) error {
	p, err := s.fpRepo.FindByDate(companyID, date)
	if err != nil {
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

// ── Year-End Closing ──────────────────────────────────────────────────────────

// ReadyToClose checks all preconditions and returns a structured readiness response.
func (s *FiscalService) ReadyToClose(companyID, fyID uuid.UUID) (response.FiscalReadinessResponse, error) {
	fy, err := s.fyRepo.FindById(companyID, fyID)
	if err != nil {
		return response.FiscalReadinessResponse{}, err
	}

	res := response.FiscalReadinessResponse{FiscalYearId: fyID, FiscalYearName: fy.Name}

	unlocked, err := s.fyRepo.CountOpenOrClosedPeriods(fyID)
	if err != nil {
		return res, err
	}
	res.UnlockedPeriods = int(unlocked)

	drafts, err := s.fyRepo.CountDraftJournalEntries(companyID, fyID)
	if err != nil {
		return res, err
	}
	res.DraftJournalEntries = int(drafts)

	res.ReadyToClose = unlocked == 0 && drafts == 0
	if !res.ReadyToClose {
		var issues []string
		if unlocked > 0 {
			issues = append(issues, fmt.Sprintf("%d period(s) not yet locked", unlocked))
		}
		if drafts > 0 {
			issues = append(issues, fmt.Sprintf("%d draft journal entry/entries exist", drafts))
		}
		res.Issues = issues
	}

	return res, nil
}

// CloseFiscalYear performs year-end close in simple or formal mode.
func (s *FiscalService) CloseFiscalYear(companyID, fyID uuid.UUID, req request.FiscalYearCloseRequest, performedBy uuid.UUID) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(companyID, fyID)
	if err != nil {
		return fy, err
	}

	if fy.Status != entity.FiscalYearActive && fy.Status != entity.FiscalYearClosingReview {
		return fy, fmt.Errorf("fiscal year must be active or in closing_review to close (current: %s)", fy.Status)
	}

	// Readiness check
	readiness, err := s.ReadyToClose(companyID, fyID)
	if err != nil {
		return fy, err
	}
	if !readiness.ReadyToClose {
		return fy, fmt.Errorf("fiscal year is not ready to close: %s", strings.Join(readiness.Issues, "; "))
	}

	// Idempotency guard
	if has, _ := s.fyRepo.HasClosingJE(fyID); has {
		return fy, fmt.Errorf("a closing journal entry already exists for this fiscal year")
	}

	mode := entity.ClosingMode(req.Mode)
	if mode != entity.ClosingModeSimple && mode != entity.ClosingModeFormal {
		return fy, fmt.Errorf("invalid closing mode '%s': must be 'simple' or 'formal'", req.Mode)
	}

	if mode == entity.ClosingModeFormal {
		if err := s.generateClosingJE(companyID, &fy, performedBy); err != nil {
			return fy, fmt.Errorf("failed to generate closing journal entry: %w", err)
		}
	}

	now := time.Now()
	fy.Status = entity.FiscalYearClosed
	fy.ClosingMode = mode
	fy.ClosedAt = &now
	fy.ClosedBy = &performedBy

	return fy, s.fyRepo.Update(fy)
}

// generateClosingJE creates the year-end P&L closing journal entry.
func (s *FiscalService) generateClosingJE(companyID uuid.UUID, fy *entity.FiscalYear, performedBy uuid.UUID) error {
	cfg, err := s.cfgRepo.FindByCompanyId(companyID)
	if err != nil || cfg.RetainedEarningsCoaId == nil {
		return fmt.Errorf("retained_earnings_coa_id is not configured — set it in Company Configuration before using formal close")
	}
	retainedEarningsCoaID := *cfg.RetainedEarningsCoaId

	rows, err := s.reportRepo.GetLedger(companyID, fy.StartDate, fy.EndDate)
	if err != nil {
		return err
	}

	type closingLine struct {
		coaID  uuid.UUID
		debit  float64
		credit float64
		desc   string
	}
	var lines []closingLine
	var netIncome float64

	for _, row := range rows {
		coa, err := s.coaRepo.FindByCode(companyID, row.AccountCode)
		if err != nil {
			continue
		}
		normalBalance := s.getGroupNormalBalance(companyID, row.GroupName)

		switch normalBalance {
		case "credit":
			// Income account — close by debiting
			creditBalance := row.TotalCredit - row.TotalDebit
			if creditBalance == 0 {
				continue
			}
			lines = append(lines, closingLine{coaID: coa.Id, debit: creditBalance, credit: 0,
				desc: fmt.Sprintf("Closing - %s", row.AccountName)})
			netIncome += creditBalance

		case "debit":
			if !s.isExpenseGroup(row.GroupName) {
				continue
			}
			debitBalance := row.TotalDebit - row.TotalCredit
			if debitBalance == 0 {
				continue
			}
			lines = append(lines, closingLine{coaID: coa.Id, debit: 0, credit: debitBalance,
				desc: fmt.Sprintf("Closing - %s", row.AccountName)})
			netIncome -= debitBalance
		}
	}

	if len(lines) == 0 {
		return nil // no P&L activity — skip
	}

	if netIncome > 0 {
		lines = append(lines, closingLine{coaID: retainedEarningsCoaID, debit: 0, credit: netIncome,
			desc: "Closing - Transfer net profit to Retained Earnings"})
	} else if netIncome < 0 {
		lines = append(lines, closingLine{coaID: retainedEarningsCoaID, debit: -netIncome, credit: 0,
			desc: "Closing - Transfer net loss to Retained Earnings"})
	}

	jeNumber := fmt.Sprintf("JE-CLOSE-%s", fy.EndDate.Format("2006"))
	je := entity.JournalEntry{
		CompanyId:     companyID,
		JournalNumber: jeNumber,
		Type:          entity.JournalTypeGeneral,
		Date:          fy.EndDate,
		Description:   fmt.Sprintf("Year-End Closing Entry — %s", fy.Name),
		Status:        entity.JournalStatusPosted,
		CreatedBy:     performedBy,
	}

	var totalDebit, totalCredit float64
	jlLines := make([]entity.JournalLine, 0, len(lines))
	for _, l := range lines {
		totalDebit += l.debit
		totalCredit += l.credit
		jlLines = append(jlLines, entity.JournalLine{
			CoaId: l.coaID, Description: l.desc, Debit: l.debit, Credit: l.credit,
		})
	}
	je.TotalDebit = totalDebit
	je.TotalCredit = totalCredit

	createdJE, err := s.jeRepo.Create(je, jlLines)
	if err != nil {
		return fmt.Errorf("failed to save closing journal entry: %w", err)
	}

	fy.ClosingJEId = &createdJE.Id
	return nil
}

// ── Opening Balance ───────────────────────────────────────────────────────────

// GenerateOpeningBalance generates an OB journal entry for fyID from previous period's ledger.
func (s *FiscalService) GenerateOpeningBalance(companyID, fyID uuid.UUID, performedBy uuid.UUID) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(companyID, fyID)
	if err != nil {
		return fy, err
	}
	if fy.Status == entity.FiscalYearClosed {
		return fy, fmt.Errorf("cannot generate opening balance for a closed fiscal year")
	}

	// Idempotency guard
	if has, _ := s.fyRepo.HasOpeningJE(fyID); has {
		return fy, fmt.Errorf("opening balance already exists for this fiscal year — delete it first to regenerate")
	}

	// Get balances up to the day before this FY starts
	asOf := fy.StartDate.AddDate(0, 0, -1)
	rows, err := s.reportRepo.GetLedgerUpTo(companyID, asOf)
	if err != nil {
		return fy, err
	}
	if len(rows) == 0 {
		return fy, fmt.Errorf("no posted journal entries found before %s — nothing to carry forward", asOf.Format("2006-01-02"))
	}

	// Check previous FY closing mode
	prevFY, prevErr := s.fyRepo.FindLastClosed(companyID)
	formalClose := prevErr == nil && prevFY.ClosingMode == entity.ClosingModeFormal

	type obLine struct {
		coaID  uuid.UUID
		debit  float64
		credit float64
		desc   string
	}
	var lines []obLine
	var totalDebit, totalCredit float64

	for _, row := range rows {
		coa, err := s.coaRepo.FindByCode(companyID, row.AccountCode)
		if err != nil {
			continue
		}

		normalBalance := s.getGroupNormalBalance(companyID, row.GroupName)

		// Formal close: skip P&L accounts (already zeroed by closing JE)
		if formalClose && s.isPLGroup(row.GroupName) {
			continue
		}

		netBalance := row.TotalDebit - row.TotalCredit
		if netBalance == 0 {
			continue
		}

		desc := fmt.Sprintf("Opening Balance - %s", row.AccountName)
		_ = normalBalance // normalBalance used for context; net balance drives debit/credit
		if netBalance > 0 {
			lines = append(lines, obLine{coaID: coa.Id, debit: netBalance, credit: 0, desc: desc})
			totalDebit += netBalance
		} else {
			lines = append(lines, obLine{coaID: coa.Id, debit: 0, credit: -netBalance, desc: desc})
			totalCredit += -netBalance
		}
	}

	if len(lines) == 0 {
		return fy, fmt.Errorf("all account balances are zero — nothing to carry forward")
	}

	// Find fiscal period for start date
	var periodID *uuid.UUID
	if p, err := s.fpRepo.FindByDate(companyID, fy.StartDate); err == nil {
		periodID = &p.Id
	}

	jeNumber := fmt.Sprintf("JE-OB-%s", fy.StartDate.Format("2006"))
	je := entity.JournalEntry{
		CompanyId:      companyID,
		FiscalPeriodId: periodID,
		JournalNumber:  jeNumber,
		Type:           entity.JournalTypeGeneral,
		Date:           fy.StartDate,
		Description:    fmt.Sprintf("Opening Balance — %s (carried from %s)", fy.Name, asOf.Format("2006-01-02")),
		Status:         entity.JournalStatusPosted,
		TotalDebit:     totalDebit,
		TotalCredit:    totalCredit,
		CreatedBy:      performedBy,
	}

	jlLines := make([]entity.JournalLine, 0, len(lines))
	for _, l := range lines {
		jlLines = append(jlLines, entity.JournalLine{
			CoaId: l.coaID, Description: l.desc, Debit: l.debit, Credit: l.credit,
		})
	}

	createdJE, err := s.jeRepo.Create(je, jlLines)
	if err != nil {
		return fy, fmt.Errorf("failed to save opening balance journal entry: %w", err)
	}

	fy.OpeningJEId = &createdJE.Id
	return fy, s.fyRepo.Update(fy)
}

// DeleteOpeningBalance removes the opening JE for regeneration.
func (s *FiscalService) DeleteOpeningBalance(companyID, fyID uuid.UUID) (entity.FiscalYear, error) {
	fy, err := s.fyRepo.FindById(companyID, fyID)
	if err != nil {
		return fy, err
	}
	if fy.Status == entity.FiscalYearClosed {
		return fy, fmt.Errorf("cannot delete opening balance of a closed fiscal year")
	}
	if fy.OpeningJEId == nil {
		return fy, fmt.Errorf("no opening balance found for this fiscal year")
	}

	if err := s.jeRepo.Delete(companyID, *fy.OpeningJEId); err != nil {
		return fy, fmt.Errorf("failed to delete opening balance journal entry: %w", err)
	}

	fy.OpeningJEId = nil
	return fy, s.fyRepo.Update(fy)
}

// ─── Private helpers ──────────────────────────────────────────────────────────

func (s *FiscalService) getGroupNormalBalance(companyID uuid.UUID, groupName string) string {
	g, err := s.coaGroupRepo.FindByName(companyID, groupName)
	if err != nil {
		return ""
	}
	return g.NormalBalance
}

func (s *FiscalService) isExpenseGroup(groupName string) bool {
	lower := strings.ToLower(strings.TrimSpace(groupName))
	return strings.Contains(lower, "beban") || strings.Contains(lower, "expense")
}

func (s *FiscalService) isPLGroup(groupName string) bool {
	lower := strings.ToLower(strings.TrimSpace(groupName))
	return strings.Contains(lower, "beban") || strings.Contains(lower, "expense") ||
		strings.Contains(lower, "pendapatan") || strings.Contains(lower, "revenue") ||
		strings.Contains(lower, "income")
}

// ─── Period generators & mappers (unchanged from original) ───────────────────

func generatePeriods(fyId uuid.UUID, fyStart, actualStart, fyEnd time.Time, pType entity.PeriodType) []entity.FiscalPeriod {
	var periods []entity.FiscalPeriod
	periodNumber := 1

	if actualStart.After(fyStart) {
		stubEnd := firstDayOfNextMonth(fyStart).Add(-time.Second)
		if stubEnd.After(fyEnd) {
			stubEnd = fyEnd
		}
		periods = append(periods, entity.FiscalPeriod{
			FiscalYearId: fyId, Name: fmt.Sprintf("Stub – %s", fyStart.Format("Jan 2006")),
			PeriodNumber: periodNumber, StartDate: actualStart, EndDate: stubEnd,
			IsStub: true, Status: entity.FiscalPeriodOpen,
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
		default:
			periodEnd = firstDayOfNextMonth(cursor).Add(-time.Second)
		}
		if periodEnd.After(fyEnd) {
			periodEnd = fyEnd
		}
		periods = append(periods, entity.FiscalPeriod{
			FiscalYearId: fyId, Name: periodName(cursor, pType),
			PeriodNumber: periodNumber, StartDate: cursor, EndDate: periodEnd,
			IsStub: false, Status: entity.FiscalPeriodOpen,
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

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func mapFiscalYear(fy entity.FiscalYear) response.FiscalYearResponse {
	return response.FiscalYearResponse{
		Id:          fy.Id,
		Name:        fy.Name,
		StartDate:   fy.StartDate,
		EndDate:     fy.EndDate,
		PeriodType:  string(fy.PeriodType),
		Status:      string(fy.Status),
		ClosingMode: string(fy.ClosingMode),
		ClosingJEId: fy.ClosingJEId,
		OpeningJEId: fy.OpeningJEId,
		ClosedAt:    fy.ClosedAt,
		ClosedBy:    fy.ClosedBy,
		CreatedBy:   fy.CreatedBy,
		CreatedAt:   fy.CreatedAt,
		UpdatedAt:   fy.UpdatedAt,
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
