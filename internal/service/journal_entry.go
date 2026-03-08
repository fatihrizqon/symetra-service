package service

import (
	"errors"
	"math"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IJournalEntryService interface {
	Create(companyID uuid.UUID, req request.JournalEntryCreateRequest, createdBy uuid.UUID, journalType entity.JournalType) (response.JournalEntryResponse, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.JournalEntryResponse, int, error)
	FindById(companyID, id uuid.UUID) (response.JournalEntryResponse, error)
	Update(companyID uuid.UUID, req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error)
	Delete(companyID, id uuid.UUID) error
	Post(companyID, id uuid.UUID) error
	Void(companyID, id uuid.UUID) error
}

type JournalEntryService struct {
	IJournalEntryRepository repository.IJournalEntryRepository
	IInvoiceRepository      repository.IInvoiceRepository
	validate                *validator.Validate
}

func NewJournalEntryService(repo repository.IJournalEntryRepository, invoiceRepo repository.IInvoiceRepository, validate *validator.Validate) IJournalEntryService {
	return &JournalEntryService{
		IJournalEntryRepository: repo,
		IInvoiceRepository:      invoiceRepo,
		validate:                validate,
	}
}

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
		Type:          string(e.Type),
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

func (s *JournalEntryService) Create(companyID uuid.UUID, req request.JournalEntryCreateRequest, createdBy uuid.UUID, journalType entity.JournalType) (response.JournalEntryResponse, error) {
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

	if journalType != entity.JournalTypeGeneral &&
		journalType != entity.JournalTypeRevenue &&
		journalType != entity.JournalTypeExpense {
		journalType = entity.JournalTypeGeneral
	}

	journalNumber, err := s.IJournalEntryRepository.GenerateJournalNumber(companyID, journalType)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	lines, totalDebit, totalCredit := buildLines(req.Lines)

	entry := entity.JournalEntry{
		CompanyId:     companyID,
		JournalNumber: journalNumber,
		Type:          journalType,
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

func (s *JournalEntryService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.JournalEntryResponse, int, error) {
	entities, totalCount, err := s.IJournalEntryRepository.FindAll(companyID, qp)
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

func (s *JournalEntryService) FindById(companyID, id uuid.UUID) (response.JournalEntryResponse, error) {
	entry, err := s.IJournalEntryRepository.FindById(companyID, id)
	if err != nil {
		return response.JournalEntryResponse{}, err
	}
	return toJournalResponse(entry), nil
}

func (s *JournalEntryService) Update(companyID uuid.UUID, req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error) {
	existing, err := s.IJournalEntryRepository.FindById(companyID, req.Id)
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

func (s *JournalEntryService) Delete(companyID, id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.JournalStatusDraft {
		return errors.New("only draft journal entries can be deleted")
	}
	return s.IJournalEntryRepository.Delete(companyID, id)
}

func (s *JournalEntryService) Post(companyID, id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(companyID, id)
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
	return s.IJournalEntryRepository.Post(companyID, id)
}

func (s *JournalEntryService) Void(companyID, id uuid.UUID) error {
	existing, err := s.IJournalEntryRepository.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.JournalStatusPosted {
		return errors.New("only posted journal entries can be voided")
	}
	// Guard: jurnal yang terikat ke invoice tidak boleh di-void manual.
	// Batalkan lewat workflow Cancel Invoice.
	if s.IInvoiceRepository.IsLinkedToJournal(id) {
		return errors.New("jurnal ini terikat ke invoice dan tidak dapat di-void secara manual — batalkan lewat menu Cancel Invoice")
	}
	return s.IJournalEntryRepository.Void(companyID, id)
}
