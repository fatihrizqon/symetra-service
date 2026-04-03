package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Purchase Order Repository ────────────────────────────────────────────────

var poSortColumns = map[string]string{
	"po_number":   "purchase_orders.po_number",
	"po_date":     "purchase_orders.po_date",
	"grand_total": "purchase_orders.grand_total",
	"status":      "purchase_orders.status",
	"created_at":  "purchase_orders.created_at",
}

type IPurchaseOrderRepository interface {
	Create(po entity.PurchaseOrder, items []entity.PurchaseOrderItem) (entity.PurchaseOrder, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error)
	FindById(companyId, id uuid.UUID) (entity.PurchaseOrder, error)
	Update(po entity.PurchaseOrder, items []entity.PurchaseOrderItem) (entity.PurchaseOrder, error)
	Delete(companyId, id uuid.UUID) error
	UpdateStatus(companyId, id uuid.UUID, status entity.PurchaseOrderStatus) error
	Save(po entity.PurchaseOrder) error
	GeneratePONumber(companyId uuid.UUID, prefix string) (string, error)
	// SelectDropdown untuk keperluan dropdown PO (bills/from-purchase-order)
	SelectDropdown(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error)
}

type PurchaseOrderRepository struct {
	Db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) IPurchaseOrderRepository {
	return &PurchaseOrderRepository{Db: db}
}

func (r *PurchaseOrderRepository) GeneratePONumber(companyId uuid.UUID, prefix string) (string, error) {
	now := time.Now()
	monthPrefix := fmt.Sprintf("%s-%d%02d", prefix, now.Year(), now.Month())
	var count int64
	if err := r.Db.Model(&entity.PurchaseOrder{}).
		Where("company_id = ? AND po_number LIKE ?", companyId, monthPrefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%04d", monthPrefix, count+1), nil
}

func (r *PurchaseOrderRepository) Create(po entity.PurchaseOrder, items []entity.PurchaseOrderItem) (entity.PurchaseOrder, error) {
	tx := r.Db.Begin()
	po.Id = uuid.New()
	if err := tx.Create(&po).Error; err != nil {
		tx.Rollback()
		return po, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].PurchaseOrderId = po.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return po, err
	}
	tx.Commit()
	return r.FindById(po.CompanyId, po.Id)
}

func (r *PurchaseOrderRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error) {
	var entities []entity.PurchaseOrder
	var totalCount int64
	query := r.Db.Model(&entity.PurchaseOrder{}).Preload("Vendor").Where("purchase_orders.company_id = ?", companyId)
	query = util.ApplySearch(query, qp)
	query = entity.PurchaseOrder{}.ApplyFilters(query, qp.Filters)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}
	query = util.ApplySort(query, qp, poSortColumns, "purchase_orders.created_at")
	query = util.ApplyPagination(query, qp)
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

// SelectDropdown — hanya PO dengan status approved (belum converted), untuk keperluan
// dropdown "Convert PO ke Bill"
func (r *PurchaseOrderRepository) SelectDropdown(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error) {
	var entities []entity.PurchaseOrder
	var totalCount int64
	query := r.Db.Model(&entity.PurchaseOrder{}).
		Preload("Vendor").
		Where("purchase_orders.company_id = ? AND purchase_orders.status = ?", companyId, entity.POStatusApproved)
	query = util.ApplySearch(query, qp)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}
	query = util.ApplyPagination(query, qp)
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *PurchaseOrderRepository) FindById(companyId, id uuid.UUID) (entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.Db.Preload("Vendor").Preload("Items").
		Where("id = ? AND company_id = ?", id, companyId).First(&po).Error
	if err != nil {
		return po, errors.New("purchase order not found")
	}
	return po, nil
}

func (r *PurchaseOrderRepository) Update(po entity.PurchaseOrder, items []entity.PurchaseOrderItem) (entity.PurchaseOrder, error) {
	tx := r.Db.Begin()
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		return po, err
	}
	if err := tx.Where("purchase_order_id = ?", po.Id).Delete(&entity.PurchaseOrderItem{}).Error; err != nil {
		tx.Rollback()
		return po, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].PurchaseOrderId = po.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return po, err
	}
	tx.Commit()
	return r.FindById(po.CompanyId, po.Id)
}

func (r *PurchaseOrderRepository) Delete(companyId, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.PurchaseOrder{}).Error
}

func (r *PurchaseOrderRepository) UpdateStatus(companyId, id uuid.UUID, status entity.PurchaseOrderStatus) error {
	return r.Db.Model(&entity.PurchaseOrder{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Update("status", status).Error
}

func (r *PurchaseOrderRepository) Save(po entity.PurchaseOrder) error {
	return r.Db.Save(&po).Error
}

// ─── Bill Repository ──────────────────────────────────────────────────────────

var billSortColumns = map[string]string{
	"bill_number":    "bills.bill_number",
	"bill_date":      "bills.bill_date",
	"due_date":       "bills.due_date",
	"grand_total":    "bills.grand_total",
	"bill_status":    "bills.bill_status",
	"payment_status": "bills.payment_status",
	"created_at":     "bills.created_at",
}

type IBillRepository interface {
	Create(bill entity.Bill, items []entity.BillItem) (entity.Bill, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int, error)
	// FIX [BUG-06]: FindByIdForUpdate menggunakan SELECT FOR UPDATE untuk mencegah race condition
	FindByIdForUpdate(tx *gorm.DB, companyId, id uuid.UUID) (entity.Bill, error)
	FindById(companyId, id uuid.UUID) (entity.Bill, error)
	Update(bill entity.Bill, items []entity.BillItem) (entity.Bill, error)
	Delete(companyId, id uuid.UUID) error
	Save(bill entity.Bill) error
	SaveTx(tx *gorm.DB, bill entity.Bill) error
	AddPayment(payment entity.BillPayment) (entity.BillPayment, error)
	AddPaymentTx(tx *gorm.DB, payment entity.BillPayment) (entity.BillPayment, error)
	GenerateBillNumber(companyId uuid.UUID, prefix string) (string, error)
	// FIX [BUG-22]: IsLinkedToJournal kini juga cek BillPayment.JournalEntryId
	IsLinkedToJournal(journalId uuid.UUID) bool
	// DB() untuk mengakses raw *gorm.DB untuk transaksi di service layer
	DB() *gorm.DB
}

type BillRepository struct {
	Db *gorm.DB
}

func NewBillRepository(db *gorm.DB) IBillRepository {
	return &BillRepository{Db: db}
}

// DB — expose raw gorm.DB untuk kebutuhan Begin() di service layer (BUG-01, BUG-02, BUG-03)
func (r *BillRepository) DB() *gorm.DB {
	return r.Db
}

func (r *BillRepository) GenerateBillNumber(companyId uuid.UUID, prefix string) (string, error) {
	now := time.Now()
	monthPrefix := fmt.Sprintf("%s-%d%02d", prefix, now.Year(), now.Month())
	var count int64
	if err := r.Db.Model(&entity.Bill{}).
		Where("company_id = ? AND bill_number LIKE ?", companyId, monthPrefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%04d", monthPrefix, count+1), nil
}

func (r *BillRepository) Create(bill entity.Bill, items []entity.BillItem) (entity.Bill, error) {
	tx := r.Db.Begin()
	bill.Id = uuid.New()
	if err := tx.Create(&bill).Error; err != nil {
		tx.Rollback()
		return bill, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].BillId = bill.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return bill, err
	}
	tx.Commit()
	return r.FindById(bill.CompanyId, bill.Id)
}

func (r *BillRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int, error) {
	var entities []entity.Bill
	var totalCount int64
	query := r.Db.Model(&entity.Bill{}).Preload("Vendor").Where("bills.company_id = ?", companyId)
	query = util.ApplySearch(query, qp)
	query = entity.Bill{}.ApplyFilters(query, qp.Filters)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}
	query = util.ApplySort(query, qp, billSortColumns, "bills.created_at")
	query = util.ApplyPagination(query, qp)
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *BillRepository) FindById(companyId, id uuid.UUID) (entity.Bill, error) {
	var bill entity.Bill
	err := r.Db.
		Preload("Vendor").
		Preload("Items").
		Preload("Items.Account").
		Preload("Payments").
		Preload("Payments.PaymentAccount").
		Preload("PurchaseOrder").
		Where("id = ? AND company_id = ?", id, companyId).First(&bill).Error
	if err != nil {
		return bill, errors.New("bill not found")
	}
	return bill, nil
}

// FIX [BUG-06]: FindByIdForUpdate — SELECT ... FOR UPDATE untuk mencegah race condition
// pada AddPayment. Harus dipanggil dalam konteks transaksi yang sudah dimulai.
func (r *BillRepository) FindByIdForUpdate(tx *gorm.DB, companyId, id uuid.UUID) (entity.Bill, error) {
	var bill entity.Bill
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&bill).Error
	if err != nil {
		return bill, errors.New("bill not found")
	}
	return bill, nil
}

func (r *BillRepository) Update(bill entity.Bill, items []entity.BillItem) (entity.Bill, error) {
	tx := r.Db.Begin()
	if err := tx.Save(&bill).Error; err != nil {
		tx.Rollback()
		return bill, err
	}
	if err := tx.Where("bill_id = ?", bill.Id).Delete(&entity.BillItem{}).Error; err != nil {
		tx.Rollback()
		return bill, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].BillId = bill.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return bill, err
	}
	tx.Commit()
	return r.FindById(bill.CompanyId, bill.Id)
}

func (r *BillRepository) Delete(companyId, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.Bill{}).Error
}

func (r *BillRepository) Save(bill entity.Bill) error {
	return r.Db.Save(&bill).Error
}

// SaveTx — Save dalam konteks transaksi yang sedang berjalan
func (r *BillRepository) SaveTx(tx *gorm.DB, bill entity.Bill) error {
	return tx.Save(&bill).Error
}

func (r *BillRepository) AddPayment(payment entity.BillPayment) (entity.BillPayment, error) {
	payment.Id = uuid.New()
	if err := r.Db.Create(&payment).Error; err != nil {
		return payment, err
	}
	return payment, nil
}

// AddPaymentTx — AddPayment dalam konteks transaksi yang sedang berjalan
func (r *BillRepository) AddPaymentTx(tx *gorm.DB, payment entity.BillPayment) (entity.BillPayment, error) {
	payment.Id = uuid.New()
	if err := tx.Create(&payment).Error; err != nil {
		return payment, err
	}
	return payment, nil
}

// FIX [BUG-22]: IsLinkedToJournal kini juga memeriksa BillPayment.JournalEntryId
// agar journal pembayaran tidak bisa di-void manual melalui halaman Journal Entry
func (r *BillRepository) IsLinkedToJournal(journalId uuid.UUID) bool {
	var count int64

	// Cek di Bill.JournalEntryId (jurnal konfirmasi)
	r.Db.Model(&entity.Bill{}).
		Where("journal_entry_id = ?", journalId).
		Count(&count)
	if count > 0 {
		return true
	}

	// FIX [BUG-22]: Cek juga di BillPayment.JournalEntryId (jurnal pembayaran)
	r.Db.Model(&entity.BillPayment{}).
		Where("journal_entry_id = ?", journalId).
		Count(&count)
	return count > 0
}
