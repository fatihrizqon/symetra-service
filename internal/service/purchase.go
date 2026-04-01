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

// ─── Purchase Order Service ───────────────────────────────────────────────────

type IPurchaseOrderService interface {
	Create(companyId uuid.UUID, req request.POCreateRequest, callerId uuid.UUID) (response.POResponse, error)
	FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.POResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.POResponse, error)
	Update(companyId uuid.UUID, req request.POUpdateRequest) (response.POResponse, error)
	Delete(companyId, id uuid.UUID) error
	Send(companyId, id uuid.UUID) error
	Approve(companyId, id uuid.UUID) error
	Decline(companyId, id uuid.UUID) error
}

type PurchaseOrderService struct {
	repo     repository.IPurchaseOrderRepository
	cfgRepo  repository.ICompanyConfigurationRepository
	validate *validator.Validate
}

func NewPurchaseOrderService(
	repo repository.IPurchaseOrderRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	validate *validator.Validate,
) IPurchaseOrderService {
	return &PurchaseOrderService{repo: repo, cfgRepo: cfgRepo, validate: validate}
}

func toPOItemResponse(it entity.PurchaseOrderItem) response.POItemResponse {
	return response.POItemResponse{
		Id: it.Id, PurchaseOrderId: it.PurchaseOrderId,
		Description: it.Description, Qty: it.Qty, Price: it.Price,
		Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
	}
}

func toPOResponse(po entity.PurchaseOrder) response.POResponse {
	items := make([]response.POItemResponse, 0, len(po.Items))
	for _, it := range po.Items {
		items = append(items, toPOItemResponse(it))
	}
	vendorName := ""
	if po.Vendor.Id != uuid.Nil {
		vendorName = po.Vendor.Name
	}
	return response.POResponse{
		Id: po.Id, CompanyId: po.CompanyId, PONumber: po.PONumber,
		VendorId: po.VendorId, VendorName: vendorName,
		PODate: po.PODate, ExpiryDate: po.ExpiryDate,
		Subtotal: po.Subtotal, DiscountTotal: po.DiscountTotal,
		Dpp: po.Dpp, TaxRate: po.TaxRate, TaxAmount: po.TaxAmount, GrandTotal: po.GrandTotal,
		Status: string(po.Status), Notes: po.Notes,
		ConvertedBillId: po.ConvertedBillId,
		CreatedBy:       po.CreatedBy, Items: items,
		CreatedAt: po.CreatedAt, UpdatedAt: po.UpdatedAt,
	}
}

func (s *PurchaseOrderService) Create(companyId uuid.UUID, req request.POCreateRequest, callerId uuid.UUID) (response.POResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.POResponse{}, err
	}
	poDate, err := parseDate(req.PODate)
	if err != nil {
		return response.POResponse{}, err
	}
	var expiryDate *time.Time
	if req.ExpiryDate != "" {
		t, err := parseDate(req.ExpiryDate)
		if err != nil {
			return response.POResponse{}, err
		}
		expiryDate = &t
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	poNumber, err := s.repo.GeneratePONumber(companyId, cfg.PurchaseOrderPrefix)
	if err != nil {
		return response.POResponse{}, err
	}

	items := make([]entity.PurchaseOrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.PurchaseOrderItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	po := entity.PurchaseOrder{
		CompanyId: companyId, PONumber: poNumber,
		VendorId: req.VendorId, PODate: poDate, ExpiryDate: expiryDate,
		Subtotal: totals.Subtotal, DiscountTotal: totals.DiscountTotal,
		Dpp: totals.Dpp, TaxRate: cfg.TaxRate, TaxAmount: totals.TaxAmount,
		GrandTotal: totals.GrandTotal, Status: entity.POStatusDraft,
		Notes: req.Notes, CreatedBy: callerId,
	}
	created, err := s.repo.Create(po, items)
	if err != nil {
		return response.POResponse{}, err
	}
	return toPOResponse(created), nil
}

func (s *PurchaseOrderService) FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.POResponse, int, error) {
	entities, totalCount, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.POResponse, 0, len(entities))
	for _, po := range entities {
		resps = append(resps, toPOResponse(po))
	}
	return resps, totalCount, nil
}

func (s *PurchaseOrderService) FindById(companyId, id uuid.UUID) (response.POResponse, error) {
	po, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.POResponse{}, err
	}
	return toPOResponse(po), nil
}

func (s *PurchaseOrderService) Update(companyId uuid.UUID, req request.POUpdateRequest) (response.POResponse, error) {
	existing, err := s.repo.FindById(companyId, req.Id)
	if err != nil {
		return response.POResponse{}, err
	}
	if existing.Status != entity.POStatusDraft {
		return response.POResponse{}, errors.New("only draft purchase orders can be edited")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.POResponse{}, err
	}
	poDate, err := parseDate(req.PODate)
	if err != nil {
		return response.POResponse{}, err
	}
	var expiryDate *time.Time
	if req.ExpiryDate != "" {
		t, err := parseDate(req.ExpiryDate)
		if err != nil {
			return response.POResponse{}, err
		}
		expiryDate = &t
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	items := make([]entity.PurchaseOrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.PurchaseOrderItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount: it.Qty*it.Price - it.Discount,
		}
	}

	existing.VendorId = req.VendorId
	existing.PODate = poDate
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
		return response.POResponse{}, err
	}
	return toPOResponse(updated), nil
}

func (s *PurchaseOrderService) Delete(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *PurchaseOrderService) Send(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be sent")
	}
	return s.repo.UpdateStatus(companyId, id, entity.POStatusSent)
}

func (s *PurchaseOrderService) Approve(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusSent {
		return errors.New("only sent purchase orders can be approved")
	}
	return s.repo.UpdateStatus(companyId, id, entity.POStatusApproved)
}

func (s *PurchaseOrderService) Decline(companyId, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.Status == entity.POStatusConverted || existing.Status == entity.POStatusDeclined {
		return errors.New("purchase order already declined or converted")
	}
	return s.repo.UpdateStatus(companyId, id, entity.POStatusDeclined)
}

// ─── Bill Service ─────────────────────────────────────────────────────────────

type IBillService interface {
	Create(companyId uuid.UUID, req request.BillCreateRequest, callerId uuid.UUID) (response.BillResponse, error)
	CreateFromPO(companyId, poId uuid.UUID, callerId uuid.UUID) (response.BillResponse, error)
	FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.BillResponse, error)
	Update(companyId uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error)
	Delete(companyId, id uuid.UUID) error
	Confirm(companyId, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error)
	AddPayment(companyId, id uuid.UUID, req request.BillPaymentRequest, callerId uuid.UUID) (response.BillResponse, error)
	Cancel(companyId, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error)
}

type BillService struct {
	billRepo    repository.IBillRepository
	poRepo      repository.IPurchaseOrderRepository
	cfgRepo     repository.ICompanyConfigurationRepository
	journalRepo repository.IJournalEntryRepository
	fiscalRepo  repository.IFiscalPeriodRepository
	validate    *validator.Validate
}

func NewBillService(
	billRepo repository.IBillRepository,
	poRepo repository.IPurchaseOrderRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	journalRepo repository.IJournalEntryRepository,
	fiscalRepo repository.IFiscalPeriodRepository,
	validate *validator.Validate,
) IBillService {
	return &BillService{
		billRepo: billRepo, poRepo: poRepo,
		cfgRepo: cfgRepo, journalRepo: journalRepo,
		fiscalRepo: fiscalRepo, validate: validate,
	}
}

func toBillItemResponse(it entity.BillItem) response.BillItemResponse {
	r := response.BillItemResponse{
		Id: it.Id, BillId: it.BillId,
		Description: it.Description, Qty: it.Qty, Price: it.Price,
		Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
		AccountId: it.AccountId,
	}
	if it.Account != nil {
		r.AccountName = it.Account.Name
	}
	return r
}

func toBillPaymentResponse(p entity.BillPayment) response.BillPaymentResponse {
	r := response.BillPaymentResponse{
		Id: p.Id, BillId: p.BillId, Amount: p.Amount,
		PaymentDate: p.PaymentDate, PaymentAccountId: p.PaymentAccountId,
		JournalEntryId: p.JournalEntryId, Notes: p.Notes,
		CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt,
	}
	if p.PaymentAccount != nil {
		r.PaymentAccountName = p.PaymentAccount.Name
	}
	return r
}

func toBillResponse(bill entity.Bill) response.BillResponse {
	items := make([]response.BillItemResponse, 0, len(bill.Items))
	for _, it := range bill.Items {
		items = append(items, toBillItemResponse(it))
	}
	payments := make([]response.BillPaymentResponse, 0, len(bill.Payments))
	for _, p := range bill.Payments {
		payments = append(payments, toBillPaymentResponse(p))
	}
	vendorName := ""
	if bill.Vendor.Id != uuid.Nil {
		vendorName = bill.Vendor.Name
	}
	return response.BillResponse{
		Id: bill.Id, CompanyId: bill.CompanyId, BillNumber: bill.BillNumber,
		PurchaseOrderId: bill.PurchaseOrderId,
		VendorId: bill.VendorId, VendorName: vendorName,
		BillDate: bill.BillDate, DueDate: bill.DueDate,
		Subtotal: bill.Subtotal, DiscountTotal: bill.DiscountTotal,
		Dpp: bill.Dpp, TaxRate: bill.TaxRate, TaxAmount: bill.TaxAmount,
		GrandTotal:    bill.GrandTotal,
		AmountPaid:    bill.AmountPaid, AmountDue: bill.AmountDue,
		BillStatus:    string(bill.BillStatus), PaymentStatus: string(bill.PaymentStatus),
		Notes:         bill.Notes, JournalEntryId: bill.JournalEntryId,
		CreatedBy:     bill.CreatedBy, Items: items, Payments: payments,
		CreatedAt: bill.CreatedAt, UpdatedAt: bill.UpdatedAt,
	}
}

func (s *BillService) buildBill(
	companyId uuid.UUID,
	vendorId uuid.UUID,
	billDate, dueDate time.Time,
	notes string,
	reqItems []request.BillItemRequest,
	callerId uuid.UUID,
	poId *uuid.UUID,
) (entity.Bill, []entity.BillItem, entity.CompanyConfiguration, error) {
	cfg, err := s.cfgRepo.FindByCompanyId(companyId)
	if err != nil {
		return entity.Bill{}, nil, cfg, err
	}
	inputs := make([]itemInput, len(reqItems))
	for i, it := range reqItems {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	billNumber, err := s.billRepo.GenerateBillNumber(companyId, cfg.BillPrefix)
	if err != nil {
		return entity.Bill{}, nil, cfg, err
	}

	items := make([]entity.BillItem, len(reqItems))
	for i, it := range reqItems {
		items[i] = entity.BillItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount:    it.Qty*it.Price - it.Discount,
			AccountId: it.AccountId,
		}
	}

	bill := entity.Bill{
		CompanyId: companyId, BillNumber: billNumber,
		PurchaseOrderId: poId, VendorId: vendorId,
		BillDate: billDate, DueDate: dueDate,
		Subtotal: totals.Subtotal, DiscountTotal: totals.DiscountTotal,
		Dpp: totals.Dpp, TaxRate: cfg.TaxRate, TaxAmount: totals.TaxAmount,
		GrandTotal:  totals.GrandTotal,
		BillStatus:  entity.BillStatusDraft,
		PaymentStatus: entity.BillPaymentUnpaid,
		Notes: notes, CreatedBy: callerId,
	}
	return bill, items, cfg, nil
}

func (s *BillService) Create(companyId uuid.UUID, req request.BillCreateRequest, callerId uuid.UUID) (response.BillResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BillResponse{}, err
	}
	vendorId, err := uuid.Parse(req.VendorId)
	if err != nil {
		return response.BillResponse{}, errors.New("invalid vendor_id")
	}
	billDate, err := parseDate(req.BillDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	bill, items, _, err := s.buildBill(companyId, vendorId, billDate, dueDate, req.Notes, req.Items, callerId, nil)
	if err != nil {
		return response.BillResponse{}, err
	}
	created, err := s.billRepo.Create(bill, items)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(created), nil
}

func (s *BillService) CreateFromPO(companyId, poId uuid.UUID, callerId uuid.UUID) (response.BillResponse, error) {
	po, err := s.poRepo.FindById(companyId, poId)
	if err != nil {
		return response.BillResponse{}, err
	}
	if po.Status == entity.POStatusConverted {
		return response.BillResponse{}, errors.New("purchase order already converted to bill")
	}
	if po.Status == entity.POStatusDeclined || po.Status == entity.POStatusExpired {
		return response.BillResponse{}, errors.New("cannot convert a declined or expired purchase order")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	billNumber, err := s.billRepo.GenerateBillNumber(companyId, cfg.BillPrefix)
	if err != nil {
		return response.BillResponse{}, err
	}

	dueDate := po.PODate.AddDate(0, 0, cfg.BillDueDays)

	billItems := make([]entity.BillItem, len(po.Items))
	for i, it := range po.Items {
		billItems[i] = entity.BillItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
		}
	}

	bill := entity.Bill{
		CompanyId: companyId, BillNumber: billNumber,
		PurchaseOrderId: &poId, VendorId: po.VendorId,
		BillDate: po.PODate, DueDate: dueDate,
		Subtotal: po.Subtotal, DiscountTotal: po.DiscountTotal,
		Dpp: po.Dpp, TaxRate: po.TaxRate, TaxAmount: po.TaxAmount,
		GrandTotal:    po.GrandTotal,
		BillStatus:    entity.BillStatusDraft,
		PaymentStatus: entity.BillPaymentUnpaid,
		Notes: po.Notes, CreatedBy: callerId,
	}

	created, err := s.billRepo.Create(bill, billItems)
	if err != nil {
		return response.BillResponse{}, err
	}

	// Mark PO as converted
	convertedId := created.Id
	po.Status = entity.POStatusConverted
	po.ConvertedBillId = &convertedId
	s.poRepo.Save(po)

	return toBillResponse(created), nil
}

func (s *BillService) FindAllTyped(companyId uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error) {
	entities, totalCount, err := s.billRepo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.BillResponse, 0, len(entities))
	for _, b := range entities {
		resps = append(resps, toBillResponse(b))
	}
	return resps, totalCount, nil
}

func (s *BillService) FindById(companyId, id uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(bill), nil
}

func (s *BillService) Update(companyId uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error) {
	existing, err := s.billRepo.FindById(companyId, req.Id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if existing.BillStatus != entity.BillStatusDraft {
		return response.BillResponse{}, errors.New("only draft bills can be edited")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.BillResponse{}, err
	}
	vendorId, err := uuid.Parse(req.VendorId)
	if err != nil {
		return response.BillResponse{}, errors.New("invalid vendor_id")
	}
	billDate, err := parseDate(req.BillDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.BillResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	items := make([]entity.BillItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.BillItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount:    it.Qty*it.Price - it.Discount,
			AccountId: it.AccountId,
		}
	}

	existing.VendorId = vendorId
	existing.BillDate = billDate
	existing.DueDate = dueDate
	existing.Notes = req.Notes
	existing.Subtotal = totals.Subtotal
	existing.DiscountTotal = totals.DiscountTotal
	existing.Dpp = totals.Dpp
	existing.TaxRate = cfg.TaxRate
	existing.TaxAmount = totals.TaxAmount
	existing.GrandTotal = totals.GrandTotal

	updated, err := s.billRepo.Update(existing, items)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(updated), nil
}

func (s *BillService) Delete(companyId, id uuid.UUID) error {
	existing, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return err
	}
	if existing.BillStatus != entity.BillStatusDraft {
		return errors.New("only draft bills can be deleted")
	}
	if err := s.billRepo.Delete(companyId, id); err != nil {
		return err
	}
	// Reset PO if this bill came from a PO
	if existing.PurchaseOrderId != nil {
		s.poRepo.UpdateStatus(companyId, *existing.PurchaseOrderId, entity.POStatusApproved)
		po, err := s.poRepo.FindById(companyId, *existing.PurchaseOrderId)
		if err == nil {
			po.ConvertedBillId = nil
			s.poRepo.Save(po)
		}
	}
	return nil
}

// Confirm — Draft → Confirmed. Creates expense Journal Entry (auto-posted).
//
//	Per BillItem:
//	  Dr. item.AccountId (fallback DefaultExpenseAccountId) = item.amount
//	If any TaxApplicable item:
//	  Dr. TaxReceivableAccountId                           = tax_amount
//	  Cr. ApAccountId                                      = grand_total
func (s *BillService) Confirm(companyId, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if bill.BillStatus != entity.BillStatusDraft {
		return response.BillResponse{}, errors.New("only draft bills can be confirmed")
	}

	cfg, err := s.cfgRepo.FindByCompanyId(companyId)
	if err != nil {
		return response.BillResponse{}, err
	}
	if cfg.ApAccountId == nil {
		return response.BillResponse{}, errors.New("accounts payable account not configured — set it in Company → Configuration")
	}

	// Validate expense accounts
	hasTaxableItem := false
	for _, it := range bill.Items {
		if it.AccountId == nil && cfg.DefaultExpenseAccountId == nil {
			return response.BillResponse{}, errors.New("item '" + it.Description + "' has no expense account and no default expense account is configured")
		}
		if it.TaxApplicable {
			hasTaxableItem = true
		}
	}
	if hasTaxableItem && bill.TaxAmount > 0 && cfg.TaxReceivableAccountId == nil {
		return response.BillResponse{}, errors.New("tax receivable account not configured — set it in Company → Configuration")
	}

	period, err := s.fiscalRepo.FindByDate(bill.BillDate)
	if err != nil {
		return response.BillResponse{}, errors.New("no open fiscal period for bill date: " + err.Error())
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeExpense)
	if err != nil {
		return response.BillResponse{}, err
	}

	// Build debit lines — one per BillItem using its account (or fallback)
	var lines []entity.JournalLine
	for _, it := range bill.Items {
		acctId := it.AccountId
		if acctId == nil {
			acctId = cfg.DefaultExpenseAccountId
		}
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *acctId,
			Description: "Biaya Bill " + bill.BillNumber + " - " + it.Description,
			Debit:       it.Amount,
			Credit:      0,
		})
	}

	// Tax receivable debit
	if hasTaxableItem && bill.TaxAmount > 0 && cfg.TaxReceivableAccountId != nil {
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *cfg.TaxReceivableAccountId,
			Description: "PPN Masukan Bill " + bill.BillNumber,
			Debit:       bill.TaxAmount,
			Credit:      0,
		})
	}

	// AP credit
	lines = append(lines, entity.JournalLine{
		Id:          uuid.New(),
		CoaId:       *cfg.ApAccountId,
		Description: "Hutang Usaha Bill " + bill.BillNumber,
		Debit:       0,
		Credit:      bill.GrandTotal,
	})

	periodId := period.Id
	journalEntry := entity.JournalEntry{
		CompanyId:      companyId,
		FiscalPeriodId: &periodId,
		JournalNumber:  journalNumber,
		Type:           entity.JournalTypeExpense,
		Date:           bill.BillDate,
		Description:    "Bill " + bill.BillNumber,
		Status:         entity.JournalStatusPosted,
		TotalDebit:     bill.GrandTotal,
		TotalCredit:    bill.GrandTotal,
		CreatedBy:      callerId,
	}

	created, err := s.journalRepo.Create(journalEntry, lines)
	if err != nil {
		return response.BillResponse{}, err
	}

	bill.BillStatus = entity.BillStatusConfirmed
	bill.PaymentStatus = entity.BillPaymentUnpaid
	bill.AmountDue = bill.GrandTotal
	bill.JournalEntryId = &created.Id
	if err := s.billRepo.Save(bill); err != nil {
		return response.BillResponse{}, err
	}

	return toBillResponse(bill), nil
}

// AddPayment — partial or full payment.
//
//	Dr. ApAccountId       = amount
//	Cr. paymentAccountId  = amount
func (s *BillService) AddPayment(companyId, id uuid.UUID, req request.BillPaymentRequest, callerId uuid.UUID) (response.BillResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BillResponse{}, err
	}
	bill, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if bill.BillStatus != entity.BillStatusConfirmed {
		return response.BillResponse{}, errors.New("only confirmed bills can receive payments")
	}
	if bill.PaymentStatus == entity.BillPaymentPaid {
		return response.BillResponse{}, errors.New("bill is already fully paid")
	}
	if req.Amount > bill.AmountDue {
		return response.BillResponse{}, errors.New("payment amount exceeds amount due")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
	if cfg.ApAccountId == nil {
		return response.BillResponse{}, errors.New("accounts payable account not configured")
	}

	paymentAccountId, err := uuid.Parse(req.PaymentAccountId)
	if err != nil {
		return response.BillResponse{}, errors.New("invalid payment_account_id")
	}
	paymentDate, err := parseDate(req.PaymentDate)
	if err != nil {
		return response.BillResponse{}, err
	}

	period, err := s.fiscalRepo.FindByDate(paymentDate)
	if err != nil {
		return response.BillResponse{}, errors.New("no open fiscal period for payment date: " + err.Error())
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeExpense)
	if err != nil {
		return response.BillResponse{}, err
	}

	lines := []entity.JournalLine{
		{
			Id:          uuid.New(),
			CoaId:       *cfg.ApAccountId,
			Description: "Pembayaran Bill " + bill.BillNumber,
			Debit:       req.Amount,
			Credit:      0,
		},
		{
			Id:          uuid.New(),
			CoaId:       paymentAccountId,
			Description: "Kas/Bank Pembayaran Bill " + bill.BillNumber,
			Debit:       0,
			Credit:      req.Amount,
		},
	}

	periodId := period.Id
	journalEntry := entity.JournalEntry{
		CompanyId:      companyId,
		FiscalPeriodId: &periodId,
		JournalNumber:  journalNumber,
		Type:           entity.JournalTypeExpense,
		Date:           paymentDate,
		Description:    "Pembayaran Bill " + bill.BillNumber,
		Status:         entity.JournalStatusPosted,
		TotalDebit:     req.Amount,
		TotalCredit:    req.Amount,
		CreatedBy:      callerId,
	}

	createdJE, err := s.journalRepo.Create(journalEntry, lines)
	if err != nil {
		return response.BillResponse{}, err
	}

	jeId := createdJE.Id
	payment := entity.BillPayment{
		BillId:           bill.Id,
		Amount:           req.Amount,
		PaymentDate:      paymentDate,
		PaymentAccountId: paymentAccountId,
		JournalEntryId:   &jeId,
		Notes:            req.Notes,
		CreatedBy:        callerId,
	}
	if _, err := s.billRepo.AddPayment(payment); err != nil {
		return response.BillResponse{}, err
	}

	// Update bill amounts
	bill.AmountPaid += req.Amount
	bill.AmountDue = bill.GrandTotal - bill.AmountPaid
	if bill.AmountDue <= 0 {
		bill.PaymentStatus = entity.BillPaymentPaid
		bill.AmountDue = 0
	} else {
		bill.PaymentStatus = entity.BillPaymentPartial
	}
	if err := s.billRepo.Save(bill); err != nil {
		return response.BillResponse{}, err
	}

	// Re-fetch to get payments loaded
	fresh, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(fresh), nil
}

// Cancel — reversal if confirmed, reject if has payments.
func (s *BillService) Cancel(companyId, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if bill.BillStatus == entity.BillStatusCancelled {
		return response.BillResponse{}, errors.New("bill already cancelled")
	}
	if len(bill.Payments) > 0 {
		return response.BillResponse{}, errors.New("cannot cancel a bill that has payments — reverse payments first")
	}

	// If confirmed, void journal entry
	if bill.BillStatus == entity.BillStatusConfirmed && bill.JournalEntryId != nil {
		cfg, _ := s.cfgRepo.FindByCompanyId(companyId)
		if cfg.ApAccountId == nil {
			return response.BillResponse{}, errors.New("accounts payable account not configured for reversal")
		}

		period, err := s.fiscalRepo.FindByDate(time.Now())
		if err != nil {
			return response.BillResponse{}, errors.New("no open fiscal period for reversal: " + err.Error())
		}

		journalNumber, _ := s.journalRepo.GenerateJournalNumber(companyId, entity.JournalTypeExpense)

		// Reverse: credit expense accounts, debit AP
		var reversalLines []entity.JournalLine
		for _, it := range bill.Items {
			acctId := it.AccountId
			if acctId == nil {
				acctId = cfg.DefaultExpenseAccountId
			}
			if acctId == nil {
				continue
			}
			reversalLines = append(reversalLines, entity.JournalLine{
				Id:          uuid.New(),
				CoaId:       *acctId,
				Description: "Pembatalan Bill " + bill.BillNumber + " - " + it.Description,
				Debit:       0,
				Credit:      it.Amount,
			})
		}
		hasTaxableItem := false
		for _, it := range bill.Items {
			if it.TaxApplicable {
				hasTaxableItem = true
				break
			}
		}
		if hasTaxableItem && bill.TaxAmount > 0 && cfg.TaxReceivableAccountId != nil {
			reversalLines = append(reversalLines, entity.JournalLine{
				Id:          uuid.New(),
				CoaId:       *cfg.TaxReceivableAccountId,
				Description: "Reversal PPN Masukan Bill " + bill.BillNumber,
				Debit:       0,
				Credit:      bill.TaxAmount,
			})
		}
		reversalLines = append(reversalLines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *cfg.ApAccountId,
			Description: "Pembatalan Hutang Bill " + bill.BillNumber,
			Debit:       bill.GrandTotal,
			Credit:      0,
		})

		periodId := period.Id
		reversal := entity.JournalEntry{
			CompanyId:      companyId,
			FiscalPeriodId: &periodId,
			JournalNumber:  journalNumber,
			Type:           entity.JournalTypeExpense,
			Date:           time.Now(),
			Description:    "Reversal Bill " + bill.BillNumber,
			Status:         entity.JournalStatusPosted,
			TotalDebit:     bill.GrandTotal,
			TotalCredit:    bill.GrandTotal,
			CreatedBy:      callerId,
		}
		s.journalRepo.Create(reversal, reversalLines)
	}

	bill.BillStatus = entity.BillStatusCancelled
	s.billRepo.Save(bill)

	// Reset PO if applicable
	if bill.PurchaseOrderId != nil {
		s.poRepo.UpdateStatus(companyId, *bill.PurchaseOrderId, entity.POStatusApproved)
		po, err := s.poRepo.FindById(companyId, *bill.PurchaseOrderId)
		if err == nil {
			po.ConvertedBillId = nil
			s.poRepo.Save(po)
		}
	}

	return toBillResponse(bill), nil
}
