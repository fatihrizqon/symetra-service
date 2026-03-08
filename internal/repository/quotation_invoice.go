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

// ─── Company Configuration ────────────────────────────────────────────────────

type ICompanyConfigurationRepository interface {
	FindByCompanyId(companyId uuid.UUID) (entity.CompanyConfiguration, error)
	Upsert(cfg entity.CompanyConfiguration) (entity.CompanyConfiguration, error)
}

type CompanyConfigurationRepository struct {
	Db *gorm.DB
}

func NewCompanyConfigurationRepository(db *gorm.DB) ICompanyConfigurationRepository {
	return &CompanyConfigurationRepository{Db: db}
}

func (r *CompanyConfigurationRepository) FindByCompanyId(companyId uuid.UUID) (entity.CompanyConfiguration, error) {
	var cfg entity.CompanyConfiguration
	err := r.Db.
		Preload("ArAccount").Preload("ApAccount").
		Preload("SalesRevenueAccount").Preload("ServiceRevenueAccount").
		Preload("TaxPayableAccount").Preload("TaxReceivableAccount").
		Preload("BankAccount").Preload("CashAccount").
		Where("company_id = ?", companyId).First(&cfg).Error
	if err != nil {
		// Return default config if not found yet
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.CompanyConfiguration{
				CompanyId:       companyId,
				EnableTax:       false,
				TaxRate:         0.11,
				InvoicePrefix:   "INV",
				QuotationPrefix: "QUO",
				InvoiceDueDays:  30,
			}, nil
		}
		return cfg, err
	}
	return cfg, nil
}

func (r *CompanyConfigurationRepository) Upsert(cfg entity.CompanyConfiguration) (entity.CompanyConfiguration, error) {
	var existing entity.CompanyConfiguration
	err := r.Db.Where("company_id = ?", cfg.CompanyId).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new
		if cfg.Id == uuid.Nil {
			cfg.Id = uuid.New()
		}
		if err := r.Db.Create(&cfg).Error; err != nil {
			return cfg, err
		}
	} else {
		// Update existing
		cfg.Id = existing.Id
		if err := r.Db.Save(&cfg).Error; err != nil {
			return cfg, err
		}
	}
	return r.FindByCompanyId(cfg.CompanyId)
}

// ─── Quotation ────────────────────────────────────────────────────────────────

var quotationSortColumns = map[string]string{
	"quotation_number": "quotations.quotation_number",
	"quotation_date":   "quotations.quotation_date",
	"grand_total":      "quotations.grand_total",
	"status":           "quotations.status",
	"created_at":       "quotations.created_at",
}

type IQuotationRepository interface {
	Create(q entity.Quotation, items []entity.QuotationItem) (entity.Quotation, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int, error)
	FindById(companyId, id uuid.UUID) (entity.Quotation, error)
	Update(q entity.Quotation, items []entity.QuotationItem) (entity.Quotation, error)
	Delete(companyId, id uuid.UUID) error
	UpdateStatus(companyId, id uuid.UUID, status entity.QuotationStatus) error
	GenerateQuotationNumber(companyId uuid.UUID, prefix string) (string, error)
}

type QuotationRepository struct {
	Db *gorm.DB
}

func NewQuotationRepository(db *gorm.DB) IQuotationRepository {
	return &QuotationRepository{Db: db}
}

func (r *QuotationRepository) GenerateQuotationNumber(companyId uuid.UUID, prefix string) (string, error) {
	now := time.Now()
	monthPrefix := fmt.Sprintf("%s-%d%02d", prefix, now.Year(), now.Month())
	var count int64
	if err := r.Db.Model(&entity.Quotation{}).
		Where("company_id = ? AND quotation_number LIKE ?", companyId, monthPrefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%04d", monthPrefix, count+1), nil
}

func (r *QuotationRepository) Create(q entity.Quotation, items []entity.QuotationItem) (entity.Quotation, error) {
	tx := r.Db.Begin()
	q.Id = uuid.New()
	if err := tx.Create(&q).Error; err != nil {
		tx.Rollback()
		return q, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].QuotationId = q.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return q, err
	}
	tx.Commit()
	return r.FindById(q.CompanyId, q.Id)
}

func (r *QuotationRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int, error) {
	var entities []entity.Quotation
	var totalCount int64
	query := r.Db.Model(&entity.Quotation{}).Preload("Customer").Where("quotations.company_id = ?", companyId)
	query = util.ApplySearch(query, qp)
	query = entity.Quotation{}.ApplyFilters(query, qp.Filters)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}
	query = util.ApplySort(query, qp, quotationSortColumns, "quotations.created_at")
	query = util.ApplyPagination(query, qp)
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *QuotationRepository) FindById(companyId, id uuid.UUID) (entity.Quotation, error) {
	var q entity.Quotation
	err := r.Db.Preload("Customer").Preload("Items").
		Where("id = ? AND company_id = ?", id, companyId).First(&q).Error
	if err != nil {
		return q, errors.New("quotation not found")
	}
	return q, nil
}

func (r *QuotationRepository) Update(q entity.Quotation, items []entity.QuotationItem) (entity.Quotation, error) {
	tx := r.Db.Begin()
	if err := tx.Save(&q).Error; err != nil {
		tx.Rollback()
		return q, err
	}
	if err := tx.Where("quotation_id = ?", q.Id).Delete(&entity.QuotationItem{}).Error; err != nil {
		tx.Rollback()
		return q, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].QuotationId = q.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return q, err
	}
	tx.Commit()
	return r.FindById(q.CompanyId, q.Id)
}

func (r *QuotationRepository) Delete(companyId, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.Quotation{}).Error
}

func (r *QuotationRepository) UpdateStatus(companyId, id uuid.UUID, status entity.QuotationStatus) error {
	return r.Db.Model(&entity.Quotation{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Update("status", status).Error
}

// ─── Invoice ──────────────────────────────────────────────────────────────────

var invoiceSortColumns = map[string]string{
	"invoice_number": "invoices.invoice_number",
	"invoice_date":   "invoices.invoice_date",
	"due_date":       "invoices.due_date",
	"grand_total":    "invoices.grand_total",
	"invoice_status": "invoices.invoice_status",
	"created_at":     "invoices.created_at",
}

type IInvoiceRepository interface {
	Create(inv entity.Invoice, items []entity.InvoiceItem) (entity.Invoice, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int, error)
	FindById(companyId, id uuid.UUID) (entity.Invoice, error)
	Update(inv entity.Invoice, items []entity.InvoiceItem) (entity.Invoice, error)
	Delete(companyId, id uuid.UUID) error
	Save(inv entity.Invoice) error
	GenerateInvoiceNumber(companyId uuid.UUID, prefix string) (string, error)
}

type InvoiceRepository struct {
	Db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) IInvoiceRepository {
	return &InvoiceRepository{Db: db}
}

func (r *InvoiceRepository) GenerateInvoiceNumber(companyId uuid.UUID, prefix string) (string, error) {
	now := time.Now()
	monthPrefix := fmt.Sprintf("%s-%d%02d", prefix, now.Year(), now.Month())
	var count int64
	if err := r.Db.Model(&entity.Invoice{}).
		Where("company_id = ? AND invoice_number LIKE ?", companyId, monthPrefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%04d", monthPrefix, count+1), nil
}

func (r *InvoiceRepository) Create(inv entity.Invoice, items []entity.InvoiceItem) (entity.Invoice, error) {
	tx := r.Db.Begin()
	inv.Id = uuid.New()
	if err := tx.Create(&inv).Error; err != nil {
		tx.Rollback()
		return inv, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].InvoiceId = inv.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return inv, err
	}
	tx.Commit()
	return r.FindById(inv.CompanyId, inv.Id)
}

func (r *InvoiceRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int, error) {
	var entities []entity.Invoice
	var totalCount int64
	query := r.Db.Model(&entity.Invoice{}).Preload("Customer").Where("invoices.company_id = ?", companyId)
	query = util.ApplySearch(query, qp)
	query = entity.Invoice{}.ApplyFilters(query, qp.Filters)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}
	query = util.ApplySort(query, qp, invoiceSortColumns, "invoices.created_at")
	query = util.ApplyPagination(query, qp)
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *InvoiceRepository) FindById(companyId, id uuid.UUID) (entity.Invoice, error) {
	var inv entity.Invoice
	err := r.Db.Preload("Customer").Preload("Items").Preload("Quotation").
		Where("id = ? AND company_id = ?", id, companyId).First(&inv).Error
	if err != nil {
		return inv, errors.New("invoice not found")
	}
	return inv, nil
}

func (r *InvoiceRepository) Update(inv entity.Invoice, items []entity.InvoiceItem) (entity.Invoice, error) {
	tx := r.Db.Begin()
	if err := tx.Save(&inv).Error; err != nil {
		tx.Rollback()
		return inv, err
	}
	if err := tx.Where("invoice_id = ?", inv.Id).Delete(&entity.InvoiceItem{}).Error; err != nil {
		tx.Rollback()
		return inv, err
	}
	for i := range items {
		items[i].Id = uuid.New()
		items[i].InvoiceId = inv.Id
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return inv, err
	}
	tx.Commit()
	return r.FindById(inv.CompanyId, inv.Id)
}

func (r *InvoiceRepository) Delete(companyId, id uuid.UUID) error {
	return r.Db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.Invoice{}).Error
}

func (r *InvoiceRepository) Save(inv entity.Invoice) error {
	return r.Db.Save(&inv).Error
}
