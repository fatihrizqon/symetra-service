package service

import (
	"errors"
	"math"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IJournalEntryService interface {
	Create(req request.JournalEntryCreateRequest, createdBy uuid.UUID) (response.JournalEntryResponse, error)
	FindAll(qp *util.QueryParams) ([]response.JournalEntryResponse, int, error)
	FindById(id uuid.UUID) (response.JournalEntryResponse, error)
	Update(req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error)
	Delete(id uuid.UUID) error
	Post(id uuid.UUID) error
	Void(id uuid.UUID) error
}

type JournalEntryService struct {
	IJournalEntryRepository repository.IJournalEntryRepository
	validate                *validator.Validate
}

func NewJournalEntryService(repo repository.IJournalEntryRepository, validate *validator.Validate) IJournalEntryService {
	return &JournalEntryService{
		IJournalEntryRepository: repo,
		validate:                validate,
	}
}

// toResponse maps a JournalEntry entity to its response DTO.
func toJournalResponse(e entity.JournalEntry) response.JournalEntryResponse {
	lines := make([]response.JournalLineResponse, 0, len(e.Lines))
	for _, l := range e.Lines {
		lr := response.JournalLineResponse{
			Id:             l.Id,
			JournalEntryId: l.JournalEntryId,
			CoaId:          l.CoaId,
			Description:    l.Description,
			Debit:          l.Debit,
			Credit:         l.Credit,
		}
		if l.COA.Id != uuid.Nil {
			lr.CoaCode = l.COA.Code
			lr.CoaName = l.COA.Name
		}
		lines = append(lines, lr)
	}

	return response.JournalEntryResponse{
		Id:            e.Id,
		JournalNumber: e.JournalNumber,
		Date:          e.Date,
		Description:   e.Description,
		Status:        string(e.Status),
		TotalDebit:    e.TotalDebit,
		TotalCredit:   e.TotalCredit,
		CreatedBy:     e.CreatedBy,
		Lines:         lines,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// validateBalance enforces that total debit equals total credit.
func validateBalance(lines []request.JournalLineRequest) error {
	var totalDebit, totalCredit float64
	for _, l := range lines {
		if l.Debit > 0 && l.Credit > 0 {
			return errors.New("a journal line cannot have both debit and credit amounts")
		}
		totalDebit += l.Debit
		totalCredit += l.Credit
	}

	diff := math.Abs(totalDebit - totalCredit)
	if diff > 0.001 {
		return errors.New("journal entry is not balanced: total debit must equal total credit")
	}

	return nil
}

// parseDate converts a date string (YYYY-MM-DD) to time.Time.
func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return t, nil
}

// buildLines converts request line DTOs to entity lines and computes totals.
func buildLines(reqLines []request.JournalLineRequest) ([]entity.JournalLine, float64, float64) {
	lines := make([]entity.JournalLine, 0, len(reqLines))
	var totalDebit, totalCredit float64

	for _, rl := range reqLines {
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       rl.CoaId,
			Description: rl.Description,
			Debit:       rl.Debit,
			Credit:      rl.Credit,
		})
		totalDebit += rl.Debit
		totalCredit += rl.Credit
	}

	return lines, totalDebit, totalCredit
}

// Create creates a new draft journal entry after validating balance.
func (s *JournalEntryService) Create(req request.JournalEntryCreateRequest, createdBy uuid.UUID) (response.JournalEntryResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.JournalEntryResponse{}, err
	}

	if err := validateBalance(req.Lines); err != nil {
		return response.JournalEntryResponse{}, err
	}

	date, err := parseDate(req.Date)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	journalNumber, err := s.IJournalEntryRepository.GenerateJournalNumber()
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	lines, totalDebit, totalCredit := buildLines(req.Lines)

	entry := entity.JournalEntry{
		JournalNumber: journalNumber,
		Date:          date,
		Description:   req.Description,
		Status:        entity.JournalStatusDraft,
		TotalDebit:    totalDebit,
		TotalCredit:   totalCredit,
		CreatedBy:     createdBy,
	}

	created, err := s.IJournalEntryRepository.Create(entry, lines)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	return toJournalResponse(created), nil
}

// FindAll retrieves paginated journal entries.
func (s *JournalEntryService) FindAll(qp *util.QueryParams) ([]response.JournalEntryResponse, int, error) {
	entities, totalCount, err := s.IJournalEntryRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.JournalEntryResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.JournalEntryResponse, 0, len(entities))

	for _, e := range entities {
		resps = append(resps, toJournalResponse(e))
	}

	return resps, totalCount, nil
}

// FindById retrieves a single journal entry by ID.
func (s *JournalEntryService) FindById(id uuid.UUID) (response.JournalEntryResponse, error) {
	entry, err := s.IJournalEntryRepository.FindById(id)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}
	return toJournalResponse(entry), nil
}

// Update replaces the header and lines of a draft journal entry.
func (s *JournalEntryService) Update(req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error) {
	existing, err := s.IJournalEntryRepository.FindById(req.Id)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	if existing.Status != entity.JournalStatusDraft {
		return response.JournalEntryResponse{}, errors.New("only draft journal entries can be edited")
	}

	if err := s.validate.Struct(req); err != nil {
		return response.JournalEntryResponse{}, err
	}

	if err := validateBalance(req.Lines); err != nil {
		return response.JournalEntryResponse{}, err
	}

	date, err := parseDate(req.Date)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	lines, totalDebit, totalCredit := buildLines(req.Lines)

	existing.Date = date
	existing.Description = req.Description
	existing.TotalDebit = totalDebit
	existing.TotalCredit = totalCredit

	updated, err := s.IJournalEntryRepository.Update(existing, lines)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	return toJournalResponse(updated), nil
}

// Delete removes a draft journal entry. Posted/void entries cannot be deleted.
func (s *JournalEntryService) Delete(id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(id)
	if err != nil {
		return err
	}

	if existing.Status != entity.JournalStatusDraft {
		return errors.New("only draft journal entries can be deleted")
	}

	return s.IJournalEntryRepository.Delete(id)
}

// Post transitions a balanced draft entry to posted.
func (s *JournalEntryService) Post(id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(id)
	if err != nil {
		return err
	}

	if existing.Status != entity.JournalStatusDraft {
		return errors.New("only draft journal entries can be posted")
	}

	diff := math.Abs(existing.TotalDebit - existing.TotalCredit)
	if diff > 0.001 {
		return errors.New("cannot post an unbalanced journal entry")
	}

	return s.IJournalEntryRepository.Post(id)
}

// Void transitions a posted entry to void.
func (s *JournalEntryService) Void(id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(id)
	if err != nil {
		return err
	}

	if existing.Status != entity.JournalStatusPosted {
		return errors.New("only posted journal entries can be voided")
	}

	return s.IJournalEntryRepository.Void(id)
}
