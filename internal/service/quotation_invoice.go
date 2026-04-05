package service

import (
	"errors"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ─── Shared helpers ───────────────────────────────────────────────────────────

// parseDate parses a YYYY-MM-DD string into time.Time.
// Single definition used by all services in this file.
func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return t, nil
}

type itemInput struct {
	Description   string
	Qty           float64
	Price         float64
	Discount      float64
	TaxApplicable bool
}

type calcResult struct {
	Subtotal      float64
	DiscountTotal float64
	Dpp           float64
	TaxAmount     float64
	GrandTotal    float64
}

// calculateTotals computes subtotal, dpp, tax, and grand total from line items.
// Tax is applied only to items with TaxApplicable=true when enableTax is true.
func calculateTotals(items []itemInput, taxRate float64, enableTax bool) calcResult {
	var subtotal, discountTotal float64
	for _, it := range items {
		subtotal += it.Qty * it.Price
		discountTotal += it.Discount
	}
	dpp := subtotal - discountTotal
	var taxAmount float64
	if enableTax {
		for _, it := range items {
			if it.TaxApplicable {
				lineNet := it.Qty*it.Price - it.Discount
				taxAmount += lineNet * taxRate
			}
		}
	}
	return calcResult{
		Subtotal:      subtotal,
		DiscountTotal: discountTotal,
		Dpp:           dpp,
		TaxAmount:     taxAmount,
		GrandTotal:    dpp + taxAmount,
	}
}

// ─── Company Configuration Service ───────────────────────────────────────────

type ICompanyConfigurationService interface {
	Get(companyId uuid.UUID) (response.CompanyConfigurationResponse, error)
	Upsert(companyId uuid.UUID, req request.CompanyConfigurationRequest) (response.CompanyConfigurationResponse, error)
}

type CompanyConfigurationService struct {
	repo     repository.ICompanyConfigurationRepository
	validate *validator.Validate
}

func NewCompanyConfigurationService(repo repository.ICompanyConfigurationRepository, validate *validator.Validate) ICompanyConfigurationService {
	return &CompanyConfigurationService{repo: repo, validate: validate}
}

func toConfigResponse(cfg entity.CompanyConfiguration) response.CompanyConfigurationResponse {
	r := response.CompanyConfigurationResponse{
		Id:                      cfg.Id,
		CompanyId:               cfg.CompanyId,
		EnableTax:               cfg.EnableTax,
		TaxRate:                 cfg.TaxRate,
		ArAccountId:             cfg.ArAccountId,
		ApAccountId:             cfg.ApAccountId,
		SalesRevenueAccountId:   cfg.SalesRevenueAccountId,
		ServiceRevenueAccountId: cfg.ServiceRevenueAccountId,
		TaxPayableAccountId:     cfg.TaxPayableAccountId,
		TaxReceivableAccountId:  cfg.TaxReceivableAccountId,
		BankAccountId:           cfg.BankAccountId,
		CashAccountId:           cfg.CashAccountId,
		InvoicePrefix:           cfg.InvoicePrefix,
		QuotationPrefix:         cfg.QuotationPrefix,
		InvoiceDueDays:          cfg.InvoiceDueDays,
		DefaultExpenseAccountId: cfg.DefaultExpenseAccountId,
		PurchaseOrderPrefix:     cfg.PurchaseOrderPrefix,
		BillPrefix:              cfg.BillPrefix,
		BillDueDays:             cfg.BillDueDays,
		CreatedAt:               cfg.CreatedAt,
		UpdatedAt:               cfg.UpdatedAt,
	}
	if cfg.ArAccount != nil {
		r.ArAccountName = cfg.ArAccount.Name
	}
	if cfg.ApAccount != nil {
		r.ApAccountName = cfg.ApAccount.Name
	}
	if cfg.SalesRevenueAccount != nil {
		r.SalesRevenueAccountName = cfg.SalesRevenueAccount.Name
	}
	if cfg.ServiceRevenueAccount != nil {
		r.ServiceRevenueAccountName = cfg.ServiceRevenueAccount.Name
	}
	if cfg.TaxPayableAccount != nil {
		r.TaxPayableAccountName = cfg.TaxPayableAccount.Name
	}
	if cfg.TaxReceivableAccount != nil {
		r.TaxReceivableAccountName = cfg.TaxReceivableAccount.Name
	}
	if cfg.BankAccount != nil {
		r.BankAccountName = cfg.BankAccount.Name
	}
	if cfg.CashAccount != nil {
		r.CashAccountName = cfg.CashAccount.Name
	}
	if cfg.DefaultExpenseAccount != nil {
		r.DefaultExpenseAccountName = cfg.DefaultExpenseAccount.Name
	}
	return r
}

func (s *CompanyConfigurationService) Get(companyId uuid.UUID) (response.CompanyConfigurationResponse, error) {
	cfg, err := s.repo.FindByCompanyId(companyId)
	if err != nil {
		return response.CompanyConfigurationResponse{}, err
	}
	return toConfigResponse(cfg), nil
}

func (s *CompanyConfigurationService) Upsert(companyId uuid.UUID, req request.CompanyConfigurationRequest) (response.CompanyConfigurationResponse, error) {
	prefix := req.InvoicePrefix
	if prefix == "" {
		prefix = "INV"
	}
	qprefix := req.QuotationPrefix
	if qprefix == "" {
		qprefix = "QUO"
	}
	dueDays := req.InvoiceDueDays
	if dueDays <= 0 {
		dueDays = 30
	}
	taxRate := req.TaxRate
	if taxRate <= 0 {
		taxRate = 0.11
	}
	poPrefix := req.PurchaseOrderPrefix
	if poPrefix == "" {
		poPrefix = "PO"
	}
	billPrefix := req.BillPrefix
	if billPrefix == "" {
		billPrefix = "BILL"
	}
	billDueDays := req.BillDueDays
	if billDueDays <= 0 {
		billDueDays = 30
	}
	cfg := entity.CompanyConfiguration{
		CompanyId:               companyId,
		EnableTax:               req.EnableTax,
		TaxRate:                 taxRate,
		ArAccountId:             req.ArAccountId,
		ApAccountId:             req.ApAccountId,
		SalesRevenueAccountId:   req.SalesRevenueAccountId,
		ServiceRevenueAccountId: req.ServiceRevenueAccountId,
		TaxPayableAccountId:     req.TaxPayableAccountId,
		TaxReceivableAccountId:  req.TaxReceivableAccountId,
		BankAccountId:           req.BankAccountId,
		CashAccountId:           req.CashAccountId,
		InvoicePrefix:           prefix,
		QuotationPrefix:         qprefix,
		InvoiceDueDays:          dueDays,
		DefaultExpenseAccountId: req.DefaultExpenseAccountId,
		PurchaseOrderPrefix:     poPrefix,
		BillPrefix:              billPrefix,
		BillDueDays:             billDueDays,
	}
	saved, err := s.repo.Upsert(cfg)
	if err != nil {
		return response.CompanyConfigurationResponse{}, err
	}
	return toConfigResponse(saved), nil
}

// ─── Quotation Service ────────────────────────────────────────────────────────

type IQuotationService interface {
	Create(companyId uuid.UUID, req request.QuotationCreateRequest, createdBy uuid.UUID) (response.QuotationResponse, error)
	FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.QuotationResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.QuotationResponse, error)
	Update(companyId uuid.UUID, req request.QuotationUpdateRequest) (response.QuotationResponse, error)
	Delete(companyId, id uuid.UUID) error
	Send(companyId, id uuid.UUID) error
	Accept(companyId, id uuid.UUID) error
	Decline(companyId, id uuid.UUID) error
}

type QuotationService struct {
	repo     repository.IQuotationRepository
	cfgRepo  repository.ICompanyConfigurationRepository
	validate *validator.Validate
}

func NewQuotationService(repo repository.IQuotationRepository, cfgRepo repository.ICompanyConfigurationRepository, validate *validator.Validate) IQuotationService {
	return &QuotationService{repo: repo, cfgRepo: cfgRepo, validate: validate}
}

func toQuotationItemResponse(it entity.QuotationItem) response.QuotationItemResponse {
	return response.QuotationItemResponse{
		Id: it.Id, QuotationId: it.QuotationId,
		Description: it.Description, Qty: it.Qty, Price: it.Price,
		Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
	}
}

func toQuotationResponse(q entity.Quotation) response.QuotationResponse {
	items := make([]response.QuotationItemResponse, 0, len(q.Items))
	for _, it := range q.Items {
		items = append(items, toQuotationItemResponse(it))
	}
	customerName := ""
	if q.Customer.Id != uuid.Nil {
		customerName = q.Customer.Name
	}
	return response.QuotationResponse{
		Id: q.Id, CompanyId: q.CompanyId, QuotationNumber: q.QuotationNumber,
		CustomerId: q.CustomerId, CustomerName: customerName,
		QuotationDate: q.QuotationDate, ExpiryDate: q.ExpiryDate,
		Subtotal: q.Subtotal, DiscountTotal: q.DiscountTotal,
		Dpp: q.Dpp, TaxRate: q.TaxRate, TaxAmount: q.TaxAmount, GrandTotal: q.GrandTotal,
		Status: string(q.Status), Notes: q.Notes,
		ConvertedInvoiceId: q.ConvertedInvoiceId,
		CreatedBy:          q.CreatedBy, Items: items,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
}

func (s *QuotationService) Create(companyId uuid.UUID, req request.QuotationCreateRequest, createdBy uuid.UUID) (response.QuotationResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.QuotationResponse{}, err
	}
	qDate, err := parseDate(req.QuotationDate)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	var expiryDate *time.Time
	if req.ExpiryDate != "" {
		t, err := parseDate(req.ExpiryDate)
		if err != nil {
			return response.QuotationResponse{}, err
		}
		expiryDate = &t
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	quotationNumber, err := s.repo.GenerateQuotationNumber(companyId, cfg.QuotationPrefix)
	if err != nil {
		return response.QuotationResponse{}, err
	}

	items := make([]entity.QuotationItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.QuotationItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	q := entity.Quotation{
		CompanyId: companyId, QuotationNumber: quotationNumber,
		CustomerId: req.CustomerId, QuotationDate: qDate, ExpiryDate: expiryDate,
		Subtotal: totals.Subtotal, DiscountTotal: totals.DiscountTotal,
		Dpp: totals.Dpp, TaxRate: cfg.TaxRate, TaxAmount: totals.TaxAmount,
		GrandTotal: totals.GrandTotal, Status: entity.QuotationStatusDraft,
		Notes: req.Notes, CreatedBy: createdBy,
	}
	created, err := s.repo.Create(q, items)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	return toQuotationResponse(created), nil
}

func (s *QuotationService) FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.QuotationResponse, int, error) {
	entities, totalCount, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.QuotationResponse, 0, len(entities))
	for _, q := range entities {
		resps = append(resps, toQuotationResponse(q))
	}
	return resps, totalCount, nil
}

func (s *QuotationService) FindById(companyId, id uuid.UUID) (response.QuotationResponse, error) {
	q, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	return toQuotationResponse(q), nil
}

func (s *QuotationService) Update(companyId uuid.UUID, req request.QuotationUpdateRequest) (response.QuotationResponse, error) {
	existing, err := s.repo.FindById(companyId, req.Id)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	if existing.Status != entity.QuotationStatusDraft {
		return response.QuotationResponse{}, errors.New("only draft quotations can be edited")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.QuotationResponse{}, err
	}
	qDate, err := parseDate(req.QuotationDate)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	var expiryDate *time.Time
	if req.ExpiryDate != "" {
		t, err := parseDate(req.ExpiryDate)
		if err != nil {
			return response.QuotationResponse{}, err
		}
		expiryDate = &t
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	items := make([]entity.QuotationItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.QuotationItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	existing.CustomerId = req.CustomerId
	existing.QuotationDate = qDate
	existing.ExpiryDate = expiryDate
	existing.Notes = req.Notes
	existing.Subtotal = totals.Subtotal
	existing.DiscountTotal = totals.DiscountTotal
	existing.Dpp = totals.Dpp
	existing.TaxRate = cfg.TaxRate
	existing.TaxAmount = totals.TaxAmount
	existing.GrandTotal = totals.GrandTotal

	updated, err := s.repo.Update(existing, items)
	if err != nil {
		return response.QuotationResponse{}, err
	}
	return toQuotationResponse(updated), nil
}

func (s *QuotationService) Delete(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *QuotationService) Send(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be sent")
	}
	return s.repo.UpdateStatus(companyId, id, entity.QuotationStatusSent)
}

func (s *QuotationService) Accept(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.QuotationStatusSent {
		return errors.New("only sent quotations can be accepted")
	}
	return s.repo.UpdateStatus(companyId, id, entity.QuotationStatusAccepted)
}

func (s *QuotationService) Decline(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status == entity.QuotationStatusConverted || existing.Status == entity.QuotationStatusDeclined {
		return errors.New("quotation already declined or converted")
	}
	return s.repo.UpdateStatus(companyId, id, entity.QuotationStatusDeclined)
}

// ─── Invoice Service ──────────────────────────────────────────────────────────

type IInvoiceService interface {
	Create(companyId uuid.UUID, req request.InvoiceCreateRequest, createdBy uuid.UUID) (response.InvoiceResponse, error)
	CreateFromQuotation(companyId, quotationId uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error)
	FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.InvoiceResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.InvoiceResponse, error)
	Update(companyId uuid.UUID, req request.InvoiceUpdateRequest) (response.InvoiceResponse, error)
	Delete(companyId, id uuid.UUID) error
	Confirm(companyId, id uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error)
	MarkPaid(companyId, id uuid.UUID, req request.InvoiceMarkPaidRequest, createdBy uuid.UUID) (response.InvoiceResponse, error)
	Cancel(companyId, id uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error)
}

type InvoiceService struct {
	invoiceRepo   repository.IInvoiceRepository
	quotationRepo repository.IQuotationRepository
	cfgRepo       repository.ICompanyConfigurationRepository
	journalRepo   repository.IJournalEntryRepository
	fiscalRepo    repository.IFiscalPeriodRepository
	validate      *validator.Validate
}

func NewInvoiceService(
	invoiceRepo repository.IInvoiceRepository,
	quotationRepo repository.IQuotationRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	journalRepo repository.IJournalEntryRepository,
	fiscalRepo repository.IFiscalPeriodRepository,
	validate *validator.Validate,
) IInvoiceService {
	return &InvoiceService{
		invoiceRepo: invoiceRepo, quotationRepo: quotationRepo,
		cfgRepo: cfgRepo, journalRepo: journalRepo,
		fiscalRepo: fiscalRepo, validate: validate,
	}
}

func toInvoiceItemResponse(it entity.InvoiceItem) response.InvoiceItemResponse {
	return response.InvoiceItemResponse{
		Id: it.Id, InvoiceId: it.InvoiceId,
		Description: it.Description, Qty: it.Qty, Price: it.Price,
		Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
	}
}

func toInvoiceResponse(inv entity.Invoice) response.InvoiceResponse {
	items := make([]response.InvoiceItemResponse, 0, len(inv.Items))
	for _, it := range inv.Items {
		items = append(items, toInvoiceItemResponse(it))
	}
	customerName := ""
	if inv.Customer.Id != uuid.Nil {
		customerName = inv.Customer.Name
	}
	return response.InvoiceResponse{
		Id: inv.Id, CompanyId: inv.CompanyId, InvoiceNumber: inv.InvoiceNumber,
		QuotationId: inv.QuotationId, CustomerId: inv.CustomerId, CustomerName: customerName,
		InvoiceDate: inv.InvoiceDate, DueDate: inv.DueDate,
		Subtotal: inv.Subtotal, DiscountTotal: inv.DiscountTotal,
		Dpp: inv.Dpp, TaxRate: inv.TaxRate, TaxAmount: inv.TaxAmount, GrandTotal: inv.GrandTotal,
		InvoiceStatus: string(inv.InvoiceStatus), TaxStatus: string(inv.TaxStatus),
		Notes: inv.Notes, JournalEntryId: inv.JournalEntryId,
		PaymentJournalId: inv.PaymentJournalId,
		CreatedBy:        inv.CreatedBy, Items: items,
		CreatedAt: inv.CreatedAt, UpdatedAt: inv.UpdatedAt,
	}
}

func (s *InvoiceService) buildInvoice(
	companyId, customerId uuid.UUID,
	invoiceDate, dueDate time.Time,
	notes string,
	reqItems []request.InvoiceItemRequest,
	createdBy uuid.UUID,
	quotationId *uuid.UUID,
) (entity.Invoice, []entity.InvoiceItem, entity.CompanyConfiguration, error) {
	cfg, err := s.cfgRepo.FindByCompanyId(companyId)
	if err != nil {
		return entity.Invoice{}, nil, cfg, err
	}
	inputs := make([]itemInput, len(reqItems))
	for i, it := range reqItems {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	invoiceNumber, err := s.invoiceRepo.GenerateInvoiceNumber(companyId, cfg.InvoicePrefix)
	if err != nil {
		return entity.Invoice{}, nil, cfg, err
	}

	items := make([]entity.InvoiceItem, len(reqItems))
	for i, it := range reqItems {
		items[i] = entity.InvoiceItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	taxStatus := entity.TaxStatusNotTaxable
	if cfg.EnableTax {
		taxStatus = entity.TaxStatusDraft
	}

	inv := entity.Invoice{
		CompanyId: companyId, InvoiceNumber: invoiceNumber,
		QuotationId: quotationId, CustomerId: customerId,
		InvoiceDate: invoiceDate, DueDate: dueDate,
		Subtotal: totals.Subtotal, DiscountTotal: totals.DiscountTotal,
		Dpp: totals.Dpp, TaxRate: cfg.TaxRate, TaxAmount: totals.TaxAmount,
		GrandTotal:    totals.GrandTotal,
		InvoiceStatus: entity.InvoiceStatusDraft, TaxStatus: taxStatus,
		Notes: notes, CreatedBy: createdBy,
	}
	return inv, items, cfg, nil
}

func (s *InvoiceService) Create(companyId uuid.UUID, req request.InvoiceCreateRequest, createdBy uuid.UUID) (response.InvoiceResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.InvoiceResponse{}, err
	}
	invDate, err := parseDate(req.InvoiceDate)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	inv, items, _, err := s.buildInvoice(companyId, req.CustomerId, invDate, dueDate, req.Notes, req.Items, createdBy, nil)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	created, err := s.invoiceRepo.Create(inv, items)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	return toInvoiceResponse(created), nil
}

func (s *InvoiceService) CreateFromQuotation(companyId, quotationId uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error) {
	q, err := s.quotationRepo.FindById(companyId, quotationId)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if q.Status == entity.QuotationStatusConverted {
		return response.InvoiceResponse{}, errors.New("quotation already converted to invoice")
	}
	if q.Status == entity.QuotationStatusDeclined || q.Status == entity.QuotationStatusExpired {
		return response.InvoiceResponse{}, errors.New("cannot convert a declined or expired quotation")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	invoiceNumber, err := s.invoiceRepo.GenerateInvoiceNumber(companyId, cfg.InvoicePrefix)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	dueDate := q.QuotationDate.AddDate(0, 0, cfg.InvoiceDueDays)

	taxStatus := entity.TaxStatusNotTaxable
	if cfg.EnableTax {
		taxStatus = entity.TaxStatusDraft
	}

	invItems := make([]entity.InvoiceItem, len(q.Items))
	for i, it := range q.Items {
		invItems[i] = entity.InvoiceItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
		}
	}

	inv := entity.Invoice{
		CompanyId: companyId, InvoiceNumber: invoiceNumber,
		QuotationId: &quotationId, CustomerId: q.CustomerId,
		InvoiceDate: q.QuotationDate, DueDate: dueDate,
		Subtotal: q.Subtotal, DiscountTotal: q.DiscountTotal,
		Dpp: q.Dpp, TaxRate: q.TaxRate, TaxAmount: q.TaxAmount,
		GrandTotal:    q.GrandTotal,
		InvoiceStatus: entity.InvoiceStatusDraft, TaxStatus: taxStatus,
		Notes: q.Notes, CreatedBy: createdBy,
	}

	created, err := s.invoiceRepo.Create(inv, invItems)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	// Mark quotation as converted
	convertedId := created.Id
	q.Status = entity.QuotationStatusConverted
	q.ConvertedInvoiceId = &convertedId
	s.quotationRepo.Update(q, nil)

	return toInvoiceResponse(created), nil
}

func (s *InvoiceService) FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.InvoiceResponse, int, error) {
	entities, totalCount, err := s.invoiceRepo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.InvoiceResponse, 0, len(entities))
	for _, inv := range entities {
		resps = append(resps, toInvoiceResponse(inv))
	}
	return resps, totalCount, nil
}

func (s *InvoiceService) FindById(companyId, id uuid.UUID) (response.InvoiceResponse, error) {
	inv, err := s.invoiceRepo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	return toInvoiceResponse(inv), nil
}

func (s *InvoiceService) Update(companyId uuid.UUID, req request.InvoiceUpdateRequest) (response.InvoiceResponse, error) {
	existing, err := s.invoiceRepo.FindById(companyId, req.Id)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if existing.InvoiceStatus != entity.InvoiceStatusDraft {
		return response.InvoiceResponse{}, errors.New("only draft invoices can be edited")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.InvoiceResponse{}, err
	}
	invDate, err := parseDate(req.InvoiceDate)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	items := make([]entity.InvoiceItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.InvoiceItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	existing.CustomerId = req.CustomerId
	existing.InvoiceDate = invDate
	existing.DueDate = dueDate
	existing.Notes = req.Notes
	existing.Subtotal = totals.Subtotal
	existing.DiscountTotal = totals.DiscountTotal
	existing.Dpp = totals.Dpp
	existing.TaxRate = cfg.TaxRate
	existing.TaxAmount = totals.TaxAmount
	existing.GrandTotal = totals.GrandTotal

	updated, err := s.invoiceRepo.Update(existing, items)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	return toInvoiceResponse(updated), nil
}

func (s *InvoiceService) Delete(companyId, id uuid.UUID) error {
	existing, err := s.invoiceRepo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.InvoiceStatus != entity.InvoiceStatusDraft {
		return errors.New("only draft invoices can be deleted")
	}

	if err := s.invoiceRepo.Delete(companyId, id); err != nil {
		return err
	}

	// Kembalikan quotation ke accepted jika invoice ini hasil convert
	if existing.QuotationId != nil {
		s.quotationRepo.UpdateStatus(companyId, *existing.QuotationId, entity.QuotationStatusAccepted)
		s.quotationRepo.ClearConvertedInvoice(companyId, *existing.QuotationId)
	}

	return nil
}

// Confirm — Draft → Confirmed. Creates AR journal entry (auto-posted).
//
//	Dr. Accounts Receivable   grand_total
//	  Cr. Sales Revenue       dpp
//	  Cr. Tax Payable (PPN)   tax_amount  (only if enable_tax && tax_amount > 0)
func (s *InvoiceService) Confirm(companyId, id uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error) {
	inv, err := s.invoiceRepo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if inv.InvoiceStatus != entity.InvoiceStatusDraft {
		return response.InvoiceResponse{}, errors.New("only draft invoices can be confirmed")
	}

	cfg, err := s.cfgRepo.FindByCompanyId(companyId)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if cfg.ArAccountId == nil {
		return response.InvoiceResponse{}, errors.New("accounts receivable account not configured — set it in Company → Configuration")
	}
	if cfg.SalesRevenueAccountId == nil {
		return response.InvoiceResponse{}, errors.New("sales revenue account not configured — set it in Company → Configuration")
	}
	if cfg.EnableTax && inv.TaxAmount > 0 && cfg.TaxPayableAccountId == nil {
		return response.InvoiceResponse{}, errors.New("tax payable account not configured — set it in Company → Configuration")
	}

	period, err := s.fiscalRepo.FindByDate(companyId, inv.InvoiceDate)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("no open fiscal period for invoice date: " + err.Error())
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeRevenue)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	lines := []entity.JournalLine{
		{
			Id:          uuid.New(),
			CoaId:       *cfg.ArAccountId,
			Description: "Piutang Invoice " + inv.InvoiceNumber,
			Debit:       inv.GrandTotal,
			Credit:      0,
		},
		{
			Id:          uuid.New(),
			CoaId:       *cfg.SalesRevenueAccountId,
			Description: "Pendapatan Invoice " + inv.InvoiceNumber,
			Debit:       0,
			Credit:      inv.Dpp,
		},
	}
	if cfg.EnableTax && inv.TaxAmount > 0 && cfg.TaxPayableAccountId != nil {
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *cfg.TaxPayableAccountId,
			Description: "PPN Keluaran Invoice " + inv.InvoiceNumber,
			Debit:       0,
			Credit:      inv.TaxAmount,
		})
	}

	periodId := period.Id
	journalEntry := entity.JournalEntry{
		CompanyId:      companyId,
		FiscalPeriodId: &periodId,
		JournalNumber:  journalNumber,
		Type:           entity.JournalTypeRevenue,
		Date:           inv.InvoiceDate,
		Description:    "Invoice " + inv.InvoiceNumber,
		Status:         entity.JournalStatusPosted,
		TotalDebit:     inv.GrandTotal,
		TotalCredit:    inv.GrandTotal,
		CreatedBy:      createdBy,
	}

	created, err := s.journalRepo.Create(journalEntry, lines)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	inv.InvoiceStatus = entity.InvoiceStatusConfirmed
	if cfg.EnableTax {
		inv.TaxStatus = entity.TaxStatusTaxable
	}
	inv.JournalEntryId = &created.Id
	if err := s.invoiceRepo.Save(inv); err != nil {
		return response.InvoiceResponse{}, err
	}

	return toInvoiceResponse(inv), nil
}

// MarkPaid — Confirmed → Paid. Creates payment journal entry (auto-posted).
//
//	Dr. Bank/Cash             grand_total
//	  Cr. Accounts Receivable grand_total
func (s *InvoiceService) MarkPaid(companyId, id uuid.UUID, req request.InvoiceMarkPaidRequest, createdBy uuid.UUID) (response.InvoiceResponse, error) {
	inv, err := s.invoiceRepo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if inv.InvoiceStatus != entity.InvoiceStatusConfirmed {
		return response.InvoiceResponse{}, errors.New("only confirmed invoices can be marked as paid")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	if cfg.ArAccountId == nil {
		return response.InvoiceResponse{}, errors.New("accounts receivable account not configured")
	}

	paymentDate, err := parseDate(req.PaymentDate)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	period, err := s.fiscalRepo.FindByDate(companyId, paymentDate)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("no open fiscal period for payment date: " + err.Error())
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeRevenue)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	lines := []entity.JournalLine{
		{
			Id:          uuid.New(),
			CoaId:       req.PaymentAccountId,
			Description: "Pembayaran Invoice " + inv.InvoiceNumber,
			Debit:       inv.GrandTotal,
			Credit:      0,
		},
		{
			Id:          uuid.New(),
			CoaId:       *cfg.ArAccountId,
			Description: "Pelunasan Piutang Invoice " + inv.InvoiceNumber,
			Debit:       0,
			Credit:      inv.GrandTotal,
		},
	}

	periodId := period.Id
	journalEntry := entity.JournalEntry{
		CompanyId:      companyId,
		FiscalPeriodId: &periodId,
		JournalNumber:  journalNumber,
		Type:           entity.JournalTypeRevenue,
		Date:           paymentDate,
		Description:    "Pembayaran Invoice " + inv.InvoiceNumber,
		Status:         entity.JournalStatusPosted,
		TotalDebit:     inv.GrandTotal,
		TotalCredit:    inv.GrandTotal,
		CreatedBy:      createdBy,
	}

	created, err := s.journalRepo.Create(journalEntry, lines)
	if err != nil {
		return response.InvoiceResponse{}, err
	}

	inv.InvoiceStatus = entity.InvoiceStatusPaid
	inv.PaymentJournalId = &created.Id
	if err := s.invoiceRepo.Save(inv); err != nil {
		return response.InvoiceResponse{}, err
	}

	return toInvoiceResponse(inv), nil
}

// Cancel — creates reversal journal if invoice was confirmed, then sets status to cancelled.
// Paid invoices cannot be cancelled.
func (s *InvoiceService) Cancel(companyId, id uuid.UUID, createdBy uuid.UUID) (response.InvoiceResponse, error) {
	inv, err := s.invoiceRepo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, err
	}
	if inv.InvoiceStatus == entity.InvoiceStatusPaid {
		return response.InvoiceResponse{}, errors.New("paid invoices cannot be cancelled")
	}
	if inv.InvoiceStatus == entity.InvoiceStatusCancelled {
		return response.InvoiceResponse{}, errors.New("invoice already cancelled")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)

	// If confirmed, create reversal journal
	if inv.InvoiceStatus == entity.InvoiceStatusConfirmed && inv.JournalEntryId != nil {
		if cfg.ArAccountId == nil || cfg.SalesRevenueAccountId == nil {
			return response.InvoiceResponse{}, errors.New("COA accounts not configured for reversal")
		}

		period, err := s.fiscalRepo.FindByDate(companyId, time.Now())
		if err != nil {
			return response.InvoiceResponse{}, errors.New("no open fiscal period for reversal: " + err.Error())
		}

		journalNumber, _ := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeRevenue)

		lines := []entity.JournalLine{
			{
				Id:          uuid.New(),
				CoaId:       *cfg.SalesRevenueAccountId,
				Description: "Pembatalan Invoice " + inv.InvoiceNumber,
				Debit:       inv.Dpp,
				Credit:      0,
			},
			{
				Id:          uuid.New(),
				CoaId:       *cfg.ArAccountId,
				Description: "Pembatalan Piutang Invoice " + inv.InvoiceNumber,
				Debit:       0,
				Credit:      inv.GrandTotal,
			},
		}
		if cfg.EnableTax && inv.TaxAmount > 0 && cfg.TaxPayableAccountId != nil {
			lines = append(lines, entity.JournalLine{
				Id:          uuid.New(),
				CoaId:       *cfg.TaxPayableAccountId,
				Description: "Reversal PPN Invoice " + inv.InvoiceNumber,
				Debit:       inv.TaxAmount,
				Credit:      0,
			})
		}

		periodId := period.Id
		reversal := entity.JournalEntry{
			CompanyId:      companyId,
			FiscalPeriodId: &periodId,
			JournalNumber:  journalNumber,
			Type:           entity.JournalTypeRevenue,
			Date:           time.Now(),
			Description:    "Reversal Invoice " + inv.InvoiceNumber,
			Status:         entity.JournalStatusPosted,
			TotalDebit:     inv.GrandTotal,
			TotalCredit:    inv.GrandTotal,
			CreatedBy:      createdBy,
		}
		s.journalRepo.Create(reversal, lines)
	}

	inv.InvoiceStatus = entity.InvoiceStatusCancelled
	inv.TaxStatus = entity.TaxStatusCancelled
	s.invoiceRepo.Save(inv)

	// Kembalikan quotation ke accepted supaya bisa di-convert ulang atau di-decline secara eksplisit
	if inv.QuotationId != nil {
		s.quotationRepo.UpdateStatus(companyId, *inv.QuotationId, entity.QuotationStatusAccepted)
		// Reset ConvertedInvoiceId di quotation
		s.quotationRepo.ClearConvertedInvoice(companyId, *inv.QuotationId)
	}

	return toInvoiceResponse(inv), nil
}
