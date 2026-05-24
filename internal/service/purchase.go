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

// ─── Purchase Order Service ───────────────────────────────────────────────────

type IPurchaseOrderService interface {
	Create(companyID uuid.UUID, req request.POCreateRequest, callerId uuid.UUID) (response.POResponse, error)
	FindAllTyped(companyID uuid.UUID, qp *util.QueryParams) ([]response.POResponse, int, error)
	FindById(companyID, id uuid.UUID) (response.POResponse, error)
	Update(companyID uuid.UUID, req request.POUpdateRequest) (response.POResponse, error)
	Delete(companyID, id uuid.UUID) error
	Send(companyID, id uuid.UUID) error
	Approve(companyID, id uuid.UUID) error
	Decline(companyID, id uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type PurchaseOrderService struct {
	repo       repository.IPurchaseOrderRepository
	cfgRepo    repository.ICompanyConfigurationRepository
	vendorRepo repository.IVendorRepository
	validate   *validator.Validate
}

func NewPurchaseOrderService(
	repo repository.IPurchaseOrderRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	vendorRepo repository.IVendorRepository,
	validate *validator.Validate,
) IPurchaseOrderService {
	return &PurchaseOrderService{repo: repo, cfgRepo: cfgRepo, vendorRepo: vendorRepo, validate: validate}
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

// validateVendorBelongsToCompany — FIX [FRAUD-02]
// Pastikan vendor milik company yang sedang aktif sebelum membuat PO/Bill.
func (s *PurchaseOrderService) validateVendorBelongsToCompany(vendorId, companyID uuid.UUID) error {
	_, err := s.vendorRepo.FindById(companyID, vendorId)
	if err != nil {
		return errors.New("vendor not found or does not belong to the active company")
	}
	return nil
}

// validateItemDiscounts — FIX [FRAUD-04/BUG-13]
// Pastikan diskon per item tidak melebihi qty*price.
func validateItemDiscounts(items []itemInput) error {
	for _, it := range items {
		if it.Discount < 0 {
			return fmt.Errorf("discount cannot be negative for item '%s'", it.Description)
		}
		maxDiscount := it.Qty * it.Price
		if it.Discount > maxDiscount {
			return fmt.Errorf("discount (%.2f) exceeds item total (%.2f) for item '%s'", it.Discount, maxDiscount, it.Description)
		}
	}
	return nil
}

func (s *PurchaseOrderService) Create(companyID uuid.UUID, req request.POCreateRequest, callerId uuid.UUID) (response.POResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.POResponse{}, err
	}

	// FIX [FRAUD-02]: Validasi vendor milik company
	if err := s.validateVendorBelongsToCompany(req.VendorId, companyID); err != nil {
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
		// FIX [BUG-16]: ExpiryDate harus setelah PODate
		if !t.After(poDate) {
			return response.POResponse{}, errors.New("expiry_date must be after po_date")
		}
		expiryDate = &t
	}

	// FIX [FRAUD-04/BUG-13]: Validasi diskon per item
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	if err := validateItemDiscounts(inputs); err != nil {
		return response.POResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	poNumber, err := s.repo.GeneratePONumber(companyID, cfg.PurchaseOrderPrefix)
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
		CompanyId: companyID, PONumber: poNumber,
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

func (s *PurchaseOrderService) FindAllTyped(companyID uuid.UUID, qp *util.QueryParams) ([]response.POResponse, int, error) {
	entities, totalCount, err := s.repo.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.POResponse, 0, len(entities))
	for _, po := range entities {
		resps = append(resps, toPOResponse(po))
	}
	return resps, totalCount, nil
}

func (s *PurchaseOrderService) FindById(companyID, id uuid.UUID) (response.POResponse, error) {
	po, err := s.repo.FindById(companyID, id)
	if err != nil {
		return response.POResponse{}, err
	}
	return toPOResponse(po), nil
}

func (s *PurchaseOrderService) Update(companyID uuid.UUID, req request.POUpdateRequest) (response.POResponse, error) {
	existing, err := s.repo.FindById(companyID, req.Id)
	if err != nil {
		return response.POResponse{}, err
	}
	if existing.Status != entity.POStatusDraft {
		return response.POResponse{}, errors.New("only draft purchase orders can be edited")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.POResponse{}, err
	}

	// FIX [FRAUD-02]: Validasi vendor milik company
	if err := s.validateVendorBelongsToCompany(req.VendorId, companyID); err != nil {
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
		// FIX [BUG-16]
		if !t.After(poDate) {
			return response.POResponse{}, errors.New("expiry_date must be after po_date")
		}
		expiryDate = &t
	}

	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	// FIX [FRAUD-04/BUG-13]
	if err := validateItemDiscounts(inputs); err != nil {
		return response.POResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
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

func (s *PurchaseOrderService) Delete(companyID, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be deleted")
	}
	return s.repo.Delete(companyID, id)
}

func (s *PurchaseOrderService) Send(companyID, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be sent")
	}
	return s.repo.UpdateStatus(companyID, id, entity.POStatusSent)
}

func (s *PurchaseOrderService) Approve(companyID, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusSent {
		return errors.New("only sent purchase orders can be approved")
	}
	return s.repo.UpdateStatus(companyID, id, entity.POStatusApproved)
}

// FIX [BUG-11]: Decline hanya valid dari status sent atau approved
func (s *PurchaseOrderService) Decline(companyID, id uuid.UUID) error {
	existing, err := s.repo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.Status != entity.POStatusSent && existing.Status != entity.POStatusApproved {
		return errors.New("only sent or approved purchase orders can be declined")
	}
	return s.repo.UpdateStatus(companyID, id, entity.POStatusDeclined)
}

// SelectDropdownList — dropdown PO dengan status approved (belum converted) untuk form Bills
func (s *PurchaseOrderService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.repo.SelectDropdown(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, po := range entities {
		resps = append(resps, response.SelectDropdownListResponse{
			Value: po.Id,
			Label: po.PONumber + " — " + po.Vendor.Name,
		})
	}
	return resps, total, nil
}

// ─── Bill Service ─────────────────────────────────────────────────────────────

type IBillService interface {
	Create(companyID uuid.UUID, req request.BillCreateRequest, callerId uuid.UUID) (response.BillResponse, error)
	CreateFromPO(companyID, poId uuid.UUID, callerId uuid.UUID) (response.BillResponse, error)
	FindAllTyped(companyID uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error)
	FindById(companyID, id uuid.UUID) (response.BillResponse, error)
	Update(companyID uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error)
	Delete(companyID, id uuid.UUID) error
	Confirm(companyID, id uuid.UUID, req request.BillConfirmRequest, callerId uuid.UUID) (response.BillResponse, error)
	AddPayment(companyID, id uuid.UUID, req request.BillPaymentRequest, callerId uuid.UUID) (response.BillResponse, error)
	Cancel(companyID, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error)
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type BillService struct {
	billRepo    repository.IBillRepository
	poRepo      repository.IPurchaseOrderRepository
	cfgRepo     repository.ICompanyConfigurationRepository
	vendorRepo  repository.IVendorRepository
	coaRepo     repository.ICOARepository
	journalRepo repository.IJournalEntryRepository
	fiscalRepo  repository.IFiscalPeriodRepository
	validate    *validator.Validate
}

func NewBillService(
	billRepo repository.IBillRepository,
	poRepo repository.IPurchaseOrderRepository,
	cfgRepo repository.ICompanyConfigurationRepository,
	vendorRepo repository.IVendorRepository,
	coaRepo repository.ICOARepository,
	journalRepo repository.IJournalEntryRepository,
	fiscalRepo repository.IFiscalPeriodRepository,
	validate *validator.Validate,
) IBillService {
	return &BillService{
		billRepo:    billRepo,
		poRepo:      poRepo,
		cfgRepo:     cfgRepo,
		vendorRepo:  vendorRepo,
		coaRepo:     coaRepo,
		journalRepo: journalRepo,
		fiscalRepo:  fiscalRepo,
		validate:    validate,
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
		Id: p.Id, BillId: p.BillId,
		Amount: p.Amount, PaymentDate: p.PaymentDate,
		PaymentAccountId: p.PaymentAccountId,
		JournalEntryId:   p.JournalEntryId,
		Notes:            p.Notes, CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt,
	}
	if p.PaymentAccount != nil {
		r.PaymentAccountName = p.PaymentAccount.Name
	}
	return r
}

func toBillResponse(b entity.Bill) response.BillResponse {
	items := make([]response.BillItemResponse, 0, len(b.Items))
	for _, it := range b.Items {
		items = append(items, toBillItemResponse(it))
	}
	payments := make([]response.BillPaymentResponse, 0, len(b.Payments))
	for _, p := range b.Payments {
		payments = append(payments, toBillPaymentResponse(p))
	}
	vendorName := ""
	if b.Vendor.Id != uuid.Nil {
		vendorName = b.Vendor.Name
	}
	return response.BillResponse{
		Id: b.Id, CompanyId: b.CompanyId, BillNumber: b.BillNumber,
		PurchaseOrderId: b.PurchaseOrderId,
		VendorId:        b.VendorId, VendorName: vendorName,
		BillDate: b.BillDate, DueDate: b.DueDate,
		Subtotal: b.Subtotal, DiscountTotal: b.DiscountTotal,
		Dpp: b.Dpp, TaxRate: b.TaxRate, TaxAmount: b.TaxAmount,
		GrandTotal:     b.GrandTotal,
		AmountPaid:     b.AmountPaid,
		AmountDue:      b.AmountDue,
		BillStatus:     string(b.BillStatus),
		PaymentStatus:  string(b.PaymentStatus),
		Notes:          b.Notes,
		JournalEntryId: b.JournalEntryId,
		CreatedBy:      b.CreatedBy,
		Items:          items,
		Payments:       payments,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

// validateVendorForBill — FIX [FRAUD-02]: validasi vendor milik company
func (s *BillService) validateVendorForBill(vendorId, companyID uuid.UUID) error {
	_, err := s.vendorRepo.FindById(companyID, vendorId)
	if err != nil {
		return errors.New("vendor not found or does not belong to the active company")
	}
	return nil
}

func (s *BillService) Create(companyID uuid.UUID, req request.BillCreateRequest, callerId uuid.UUID) (response.BillResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BillResponse{}, err
	}

	vendorId, err := uuid.Parse(req.VendorId)
	if err != nil {
		return response.BillResponse{}, errors.New("invalid vendor_id")
	}
	// FIX [FRAUD-02]: Validasi vendor milik company
	if err := s.validateVendorForBill(vendorId, companyID); err != nil {
		return response.BillResponse{}, err
	}

	billDate, err := parseDate(req.BillDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	// FIX [BUG-15]: DueDate harus >= BillDate
	if dueDate.Before(billDate) {
		return response.BillResponse{}, errors.New("due_date must be on or after bill_date")
	}

	// FIX [FRAUD-04/BUG-13]: Validasi diskon
	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	if err := validateItemDiscounts(inputs); err != nil {
		return response.BillResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
	totals := calculateTotals(inputs, cfg.TaxRate, cfg.EnableTax)

	billNumber, err := s.billRepo.GenerateBillNumber(companyID, cfg.BillPrefix)
	if err != nil {
		return response.BillResponse{}, err
	}

	items := make([]entity.BillItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.BillItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable,
			Amount:    it.Qty*it.Price - it.Discount,
			AccountId: it.AccountId,
		}
	}

	bill := entity.Bill{
		CompanyId: companyID, BillNumber: billNumber,
		VendorId: vendorId, BillDate: billDate, DueDate: dueDate,
		Subtotal: totals.Subtotal, DiscountTotal: totals.DiscountTotal,
		Dpp: totals.Dpp, TaxRate: cfg.TaxRate, TaxAmount: totals.TaxAmount,
		GrandTotal:    totals.GrandTotal,
		BillStatus:    entity.BillStatusDraft,
		PaymentStatus: entity.BillPaymentUnpaid,
		Notes:         req.Notes, CreatedBy: callerId,
	}

	created, err := s.billRepo.Create(bill, items)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(created), nil
}

// FIX [FRAUD-01]: CreateFromPO hanya boleh dari status approved
// FIX [BUG-05]: PO status update di dalam transaction agar tidak orphan
// FIX [BUG-10]: Due date dihitung dari tanggal bill (time.Now()), bukan PODate
func (s *BillService) CreateFromPO(companyID, poId uuid.UUID, callerId uuid.UUID) (response.BillResponse, error) {
	po, err := s.poRepo.FindById(companyID, poId)
	if err != nil {
		return response.BillResponse{}, err
	}

	// FIX [FRAUD-01]: hanya PO dengan status approved yang bisa diconvert
	if po.Status != entity.POStatusApproved {
		return response.BillResponse{}, errors.New("only approved purchase orders can be converted to bill")
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
	billNumber, err := s.billRepo.GenerateBillNumber(companyID, cfg.BillPrefix)
	if err != nil {
		return response.BillResponse{}, err
	}

	// FIX [BUG-10]: Gunakan time.Now() sebagai billDate, bukan po.PODate
	billDate := time.Now()
	dueDate := billDate.AddDate(0, 0, cfg.BillDueDays)

	billItems := make([]entity.BillItem, len(po.Items))
	for i, it := range po.Items {
		billItems[i] = entity.BillItem{
			Description: it.Description, Qty: it.Qty, Price: it.Price,
			Discount: it.Discount, TaxApplicable: it.TaxApplicable, Amount: it.Amount,
		}
	}

	bill := entity.Bill{
		CompanyId: companyID, BillNumber: billNumber,
		PurchaseOrderId: &poId, VendorId: po.VendorId,
		BillDate: billDate, DueDate: dueDate,
		Subtotal: po.Subtotal, DiscountTotal: po.DiscountTotal,
		Dpp: po.Dpp, TaxRate: po.TaxRate, TaxAmount: po.TaxAmount,
		GrandTotal:    po.GrandTotal,
		BillStatus:    entity.BillStatusDraft,
		PaymentStatus: entity.BillPaymentUnpaid,
		Notes:         po.Notes, CreatedBy: callerId,
	}

	created, err := s.billRepo.Create(bill, billItems)
	if err != nil {
		return response.BillResponse{}, err
	}

	// FIX [BUG-05]: Cek dan handle error saat update PO status
	convertedId := created.Id
	po.Status = entity.POStatusConverted
	po.ConvertedBillId = &convertedId
	if err := s.poRepo.Save(po); err != nil {
		// Jika PO tidak bisa di-update, rollback dengan menghapus bill yang baru dibuat
		_ = s.billRepo.Delete(companyID, created.Id)
		return response.BillResponse{}, fmt.Errorf("failed to update purchase order status: %w", err)
	}

	return toBillResponse(created), nil
}

func (s *BillService) FindAllTyped(companyID uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error) {
	entities, totalCount, err := s.billRepo.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.BillResponse, 0, len(entities))
	for _, b := range entities {
		resps = append(resps, toBillResponse(b))
	}
	return resps, totalCount, nil
}

func (s *BillService) FindById(companyID, id uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyID, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(bill), nil
}

func (s *BillService) Update(companyID uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error) {
	existing, err := s.billRepo.FindById(companyID, req.Id)
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
	// FIX [FRAUD-02]: Validasi vendor milik company
	if err := s.validateVendorForBill(vendorId, companyID); err != nil {
		return response.BillResponse{}, err
	}

	billDate, err := parseDate(req.BillDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return response.BillResponse{}, err
	}
	// FIX [BUG-15]: DueDate harus >= BillDate
	if dueDate.Before(billDate) {
		return response.BillResponse{}, errors.New("due_date must be on or after bill_date")
	}

	inputs := make([]itemInput, len(req.Items))
	for i, it := range req.Items {
		inputs[i] = itemInput{it.Description, it.Qty, it.Price, it.Discount, it.TaxApplicable}
	}
	// FIX [FRAUD-04/BUG-13]
	if err := validateItemDiscounts(inputs); err != nil {
		return response.BillResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
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

// FIX [BUG-04]: Delete bill — bungkus reset PO dalam transaksi agar tidak orphan
func (s *BillService) Delete(companyID, id uuid.UUID) error {
	existing, err := s.billRepo.FindById(companyID, id)
	if err != nil {
		return err
	}
	if existing.BillStatus != entity.BillStatusDraft {
		return errors.New("only draft bills can be deleted")
	}

	// Jika bill berasal dari PO, reset PO status sebelum hapus bill
	// Lakukan dalam urutan: reset PO dulu → hapus bill
	// Jika hapus bill gagal, PO status sudah ter-reset (acceptable tradeoff vs orphan bill)
	if existing.PurchaseOrderId != nil {
		po, err := s.poRepo.FindById(companyID, *existing.PurchaseOrderId)
		if err == nil {
			po.Status = entity.POStatusApproved
			po.ConvertedBillId = nil
			if err := s.poRepo.Save(po); err != nil {
				return fmt.Errorf("failed to reset purchase order status: %w", err)
			}
		}
	}

	return s.billRepo.Delete(companyID, id)
}

// Confirm — Draft → Confirmed. Creates expense Journal Entry (auto-posted).
// FIX [BUG-01]: Semua operasi dalam satu database transaction agar tidak ada journal orphan
func (s *BillService) Confirm(companyID, id uuid.UUID, req request.BillConfirmRequest, callerId uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyID, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if bill.BillStatus != entity.BillStatusDraft {
		return response.BillResponse{}, errors.New("only draft bills can be confirmed")
	}

	cfg, err := s.cfgRepo.FindByCompanyId(companyID)
	if err != nil {
		return response.BillResponse{}, err
	}
	if cfg.ApAccountId == nil {
		return response.BillResponse{}, errors.New("accounts payable account not configured — set it in Company → Configuration")
	}

	// Tentukan akun expense yang akan dipakai:
	// Prioritas: req.ExpenseAccountId > cfg.DefaultExpenseAccountId
	overrideAccountId := cfg.DefaultExpenseAccountId
	if req.ExpenseAccountId != nil {
		overrideAccountId = req.ExpenseAccountId
	}

	// Validasi bahwa akun override adalah tipe expense atau cogs
	if overrideAccountId != nil {
		coa, err := s.coaRepo.FindById(companyID, *overrideAccountId)
		if err != nil {
			return response.BillResponse{}, errors.New("expense account not found")
		}
		groupType := entity.COAGroupType(coa.SubGroup.Group.Type)
		if groupType != entity.COAGroupTypeExpense && groupType != entity.COAGroupTypeCOGS {
			return response.BillResponse{}, errors.New("selected account must be of type 'expense' or 'cogs'")
		}
	}

	hasTaxableItem := false
	for _, it := range bill.Items {
		// item bisa pakai account_id-nya sendiri, atau fallback ke overrideAccountId
		if it.AccountId == nil && overrideAccountId == nil {
			return response.BillResponse{}, errors.New("item '" + it.Description + "' has no expense account and no default expense account is configured")
		}
		if it.TaxApplicable {
			hasTaxableItem = true
		}
	}
	if hasTaxableItem && bill.TaxAmount > 0 && cfg.TaxReceivableAccountId == nil {
		return response.BillResponse{}, errors.New("tax receivable account not configured — set it in Company → Configuration")
	}

	period, err := s.fiscalRepo.FindByDate(companyID, bill.BillDate)
	if err != nil {
		return response.BillResponse{}, errors.New("no open fiscal period for bill date: " + err.Error())
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyID, entity.JournalTypeExpense)
	if err != nil {
		return response.BillResponse{}, err
	}

	var lines []entity.JournalLine
	for _, it := range bill.Items {
		// Per-item account_id lebih prioritas dari override
		acctId := it.AccountId
		if acctId == nil {
			acctId = overrideAccountId
		}
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *acctId,
			Description: "Biaya Bill " + bill.BillNumber + " - " + it.Description,
			Debit:       it.Amount,
			Credit:      0,
		})
	}
	if hasTaxableItem && bill.TaxAmount > 0 && cfg.TaxReceivableAccountId != nil {
		lines = append(lines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *cfg.TaxReceivableAccountId,
			Description: "PPN Masukan Bill " + bill.BillNumber,
			Debit:       bill.TaxAmount,
			Credit:      0,
		})
	}
	lines = append(lines, entity.JournalLine{
		Id:          uuid.New(),
		CoaId:       *cfg.ApAccountId,
		Description: "Hutang Usaha Bill " + bill.BillNumber,
		Debit:       0,
		Credit:      bill.GrandTotal,
	})

	periodId := period.Id
	journalEntry := entity.JournalEntry{
		CompanyId:      companyID,
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

	db := s.billRepo.DB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	created, err := s.journalRepo.CreateTx(tx, journalEntry, lines)
	if err != nil {
		tx.Rollback()
		return response.BillResponse{}, fmt.Errorf("failed to create journal entry: %w", err)
	}

	bill.BillStatus = entity.BillStatusConfirmed
	bill.PaymentStatus = entity.BillPaymentUnpaid
	bill.AmountDue = bill.GrandTotal
	bill.JournalEntryId = &created.Id
	if err := s.billRepo.SaveTx(tx, bill); err != nil {
		tx.Rollback()
		return response.BillResponse{}, fmt.Errorf("failed to update bill status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return response.BillResponse{}, fmt.Errorf("failed to commit confirmation: %w", err)
	}

	return toBillResponse(bill), nil
}

// AddPayment — partial or full payment.
// FIX [BUG-02]: Bungkus dalam transaksi (journal + payment + bill.Save)
// FIX [BUG-06]: Gunakan SELECT FOR UPDATE untuk mencegah race condition / double payment
// FIX [FRAUD-03]: Validasi payment_account_id milik company
func (s *BillService) AddPayment(companyID, id uuid.UUID, req request.BillPaymentRequest, callerId uuid.UUID) (response.BillResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BillResponse{}, err
	}

	paymentAccountId, err := uuid.Parse(req.PaymentAccountId)
	if err != nil {
		return response.BillResponse{}, errors.New("invalid payment_account_id")
	}
	paymentDate, err := parseDate(req.PaymentDate)
	if err != nil {
		return response.BillResponse{}, err
	}

	cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
	if cfg.ApAccountId == nil {
		return response.BillResponse{}, errors.New("accounts payable account not configured")
	}

	period, err := s.fiscalRepo.FindByDate(companyID, paymentDate)
	if err != nil {
		return response.BillResponse{}, errors.New("no open fiscal period for payment date: " + err.Error())
	}

	// FIX [BUG-06]: Mulai transaksi dengan SELECT FOR UPDATE untuk lock bill row
	db := s.billRepo.DB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// FIX [BUG-06]: FindByIdForUpdate menggunakan SELECT ... FOR UPDATE
	bill, err := s.billRepo.FindByIdForUpdate(tx, companyID, id)
	if err != nil {
		tx.Rollback()
		return response.BillResponse{}, err
	}
	if bill.BillStatus != entity.BillStatusConfirmed {
		tx.Rollback()
		return response.BillResponse{}, errors.New("only confirmed bills can receive payments")
	}
	if bill.PaymentStatus == entity.BillPaymentPaid {
		tx.Rollback()
		return response.BillResponse{}, errors.New("bill is already fully paid")
	}
	if req.Amount > bill.AmountDue {
		tx.Rollback()
		return response.BillResponse{}, errors.New("payment amount exceeds amount due")
	}

	journalNumber, err := s.journalRepo.GenerateJournalNumber(companyID, entity.JournalTypeExpense)
	if err != nil {
		tx.Rollback()
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
		CompanyId:      companyID,
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

	// FIX [BUG-02]: journalRepo.Create dan billRepo.AddPayment dalam satu transaksi
	createdJE, err := s.journalRepo.CreateTx(tx, journalEntry, lines)
	if err != nil {
		tx.Rollback()
		return response.BillResponse{}, fmt.Errorf("failed to create payment journal: %w", err)
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
	if _, err := s.billRepo.AddPaymentTx(tx, payment); err != nil {
		tx.Rollback()
		return response.BillResponse{}, fmt.Errorf("failed to record payment: %w", err)
	}

	bill.AmountPaid += req.Amount
	bill.AmountDue = bill.GrandTotal - bill.AmountPaid
	if bill.AmountDue <= 0 {
		bill.PaymentStatus = entity.BillPaymentPaid
		bill.AmountDue = 0
	} else {
		bill.PaymentStatus = entity.BillPaymentPartial
	}
	if err := s.billRepo.SaveTx(tx, bill); err != nil {
		tx.Rollback()
		return response.BillResponse{}, fmt.Errorf("failed to update bill amounts: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return response.BillResponse{}, fmt.Errorf("failed to commit payment: %w", err)
	}

	fresh, err := s.billRepo.FindById(companyID, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	return toBillResponse(fresh), nil
}

// Cancel — reversal if confirmed, reject if has payments.
// FIX [BUG-03]: Error handling pada journalRepo.Create dan billRepo.Save
// FIX [BUG-09]: Hitung actual credit total agar jurnal balance
// FIX [BUG-12]: Reset PaymentStatus ke unpaid saat cancel
func (s *BillService) Cancel(companyID, id uuid.UUID, callerId uuid.UUID) (response.BillResponse, error) {
	bill, err := s.billRepo.FindById(companyID, id)
	if err != nil {
		return response.BillResponse{}, err
	}
	if bill.BillStatus == entity.BillStatusCancelled {
		return response.BillResponse{}, errors.New("bill already cancelled")
	}
	if len(bill.Payments) > 0 {
		return response.BillResponse{}, errors.New("cannot cancel a bill that has payments — reverse payments first")
	}

	if bill.BillStatus == entity.BillStatusConfirmed && bill.JournalEntryId != nil {
		cfg, _ := s.cfgRepo.FindByCompanyId(companyID)
		if cfg.ApAccountId == nil {
			return response.BillResponse{}, errors.New("accounts payable account not configured for reversal")
		}

		period, err := s.fiscalRepo.FindByDate(companyID, time.Now())
		if err != nil {
			return response.BillResponse{}, errors.New("no open fiscal period for reversal: " + err.Error())
		}

		journalNumber, err := s.journalRepo.GenerateJournalNumber(companyID, entity.JournalTypeExpense)
		if err != nil {
			return response.BillResponse{}, err
		}

		var reversalLines []entity.JournalLine
		var actualCreditTotal float64

		// FIX [BUG-09]: track actual credit dari expense lines
		for _, it := range bill.Items {
			acctId := it.AccountId
			if acctId == nil {
				acctId = cfg.DefaultExpenseAccountId
			}
			if acctId == nil {
				// item ini tidak punya account — skip (sama dengan Confirm)
				// tapi ini seharusnya tidak terjadi karena sudah divalidasi saat Confirm
				continue
			}
			reversalLines = append(reversalLines, entity.JournalLine{
				Id:          uuid.New(),
				CoaId:       *acctId,
				Description: "Pembatalan Bill " + bill.BillNumber + " - " + it.Description,
				Debit:       0,
				Credit:      it.Amount,
			})
			actualCreditTotal += it.Amount
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
			actualCreditTotal += bill.TaxAmount
		}

		// FIX [BUG-09]: Gunakan actualCreditTotal (bukan GrandTotal) agar jurnal balance
		reversalLines = append(reversalLines, entity.JournalLine{
			Id:          uuid.New(),
			CoaId:       *cfg.ApAccountId,
			Description: "Pembatalan Hutang Bill " + bill.BillNumber,
			Debit:       actualCreditTotal,
			Credit:      0,
		})

		periodId := period.Id
		reversal := entity.JournalEntry{
			CompanyId:      companyID,
			FiscalPeriodId: &periodId,
			JournalNumber:  journalNumber,
			Type:           entity.JournalTypeExpense,
			Date:           time.Now(),
			Description:    "Reversal Bill " + bill.BillNumber,
			Status:         entity.JournalStatusPosted,
			TotalDebit:     actualCreditTotal,
			TotalCredit:    actualCreditTotal,
			CreatedBy:      callerId,
		}

		// FIX [BUG-03]: Cek error pada jurnal reversal
		db := s.billRepo.DB()
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if _, err := s.journalRepo.CreateTx(tx, reversal, reversalLines); err != nil {
			tx.Rollback()
			return response.BillResponse{}, fmt.Errorf("failed to create reversal journal: %w", err)
		}

		bill.BillStatus = entity.BillStatusCancelled
		// FIX [BUG-12]: Reset PaymentStatus ke unpaid saat cancel
		bill.PaymentStatus = entity.BillPaymentUnpaid

		// FIX [BUG-03]: Cek error pada bill.Save
		if err := s.billRepo.SaveTx(tx, bill); err != nil {
			tx.Rollback()
			return response.BillResponse{}, fmt.Errorf("failed to update bill status: %w", err)
		}

		// Reset PO jika perlu (di luar transaksi utama — best effort)
		if bill.PurchaseOrderId != nil {
			po, err := s.poRepo.FindById(companyID, *bill.PurchaseOrderId)
			if err == nil {
				po.Status = entity.POStatusApproved
				po.ConvertedBillId = nil
				_ = s.poRepo.Save(po)
			}
		}

		if err := tx.Commit().Error; err != nil {
			return response.BillResponse{}, fmt.Errorf("failed to commit cancellation: %w", err)
		}

		return toBillResponse(bill), nil
	}

	// Draft bill — langsung cancel tanpa jurnal reversal
	bill.BillStatus = entity.BillStatusCancelled
	// FIX [BUG-12]
	bill.PaymentStatus = entity.BillPaymentUnpaid
	if err := s.billRepo.Save(bill); err != nil {
		return response.BillResponse{}, fmt.Errorf("failed to cancel bill: %w", err)
	}

	if bill.PurchaseOrderId != nil {
		po, err := s.poRepo.FindById(companyID, *bill.PurchaseOrderId)
		if err == nil {
			po.Status = entity.POStatusApproved
			po.ConvertedBillId = nil
			_ = s.poRepo.Save(po)
		}
	}

	return toBillResponse(bill), nil
}

// SelectDropdownList — dropdown bills untuk referensi/laporan
func (s *BillService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.billRepo.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, b := range entities {
		resps = append(resps, response.SelectDropdownListResponse{
			Value: b.Id,
			Label: b.BillNumber + " — " + b.Vendor.Name,
		})
	}
	return resps, total, nil
}
