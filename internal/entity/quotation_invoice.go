package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Company Configuration ────────────────────────────────────────────────────

func (CompanyConfiguration) TableName() string { return "company_configurations" }

type CompanyConfiguration struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;" json:"company_id"`

	EnableTax bool    `gorm:"type:boolean;not null;default:false;" json:"enable_tax"`
	TaxRate   float64 `gorm:"type:numeric(5,4);not null;default:0.11;" json:"tax_rate"`

	// COA FKs — all nullable; validated in service before journal creation
	ArAccountId             *uuid.UUID `gorm:"type:uuid;" json:"ar_account_id"`
	ArAccount               *COA       `gorm:"foreignKey:ArAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ar_account,omitempty"`
	ApAccountId             *uuid.UUID `gorm:"type:uuid;" json:"ap_account_id"`
	ApAccount               *COA       `gorm:"foreignKey:ApAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ap_account,omitempty"`
	SalesRevenueAccountId   *uuid.UUID `gorm:"type:uuid;" json:"sales_revenue_account_id"`
	SalesRevenueAccount     *COA       `gorm:"foreignKey:SalesRevenueAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sales_revenue_account,omitempty"`
	ServiceRevenueAccountId *uuid.UUID `gorm:"type:uuid;" json:"service_revenue_account_id"`
	ServiceRevenueAccount   *COA       `gorm:"foreignKey:ServiceRevenueAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"service_revenue_account,omitempty"`
	TaxPayableAccountId     *uuid.UUID `gorm:"type:uuid;" json:"tax_payable_account_id"`
	TaxPayableAccount       *COA       `gorm:"foreignKey:TaxPayableAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax_payable_account,omitempty"`
	TaxReceivableAccountId  *uuid.UUID `gorm:"type:uuid;" json:"tax_receivable_account_id"`
	TaxReceivableAccount    *COA       `gorm:"foreignKey:TaxReceivableAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax_receivable_account,omitempty"`
	BankAccountId           *uuid.UUID `gorm:"type:uuid;" json:"bank_account_id"`
	BankAccount             *COA       `gorm:"foreignKey:BankAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"bank_account,omitempty"`
	CashAccountId           *uuid.UUID `gorm:"type:uuid;" json:"cash_account_id"`
	CashAccount             *COA       `gorm:"foreignKey:CashAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cash_account,omitempty"`

	InvoicePrefix   string `gorm:"type:character varying;not null;default:'INV';" json:"invoice_prefix"`
	QuotationPrefix string `gorm:"type:character varying;not null;default:'QUO';" json:"quotation_prefix"`
	InvoiceDueDays  int    `gorm:"type:integer;not null;default:30;" json:"invoice_due_days"`

	// Purchases
	DefaultExpenseAccountId *uuid.UUID `gorm:"type:uuid;" json:"default_expense_account_id"`
	DefaultExpenseAccount   *COA       `gorm:"foreignKey:DefaultExpenseAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"default_expense_account,omitempty"`
	PurchaseOrderPrefix     string     `gorm:"type:character varying;not null;default:'PO';" json:"purchase_order_prefix"`
	BillPrefix              string     `gorm:"type:character varying;not null;default:'BILL';" json:"bill_prefix"`
	BillDueDays             int        `gorm:"type:integer;not null;default:30;" json:"bill_due_days"`

	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

// ─── Quotation ────────────────────────────────────────────────────────────────

func (Quotation) TableName() string { return "quotations" }

type QuotationStatus string

const (
	QuotationStatusDraft     QuotationStatus = "draft"
	QuotationStatusSent      QuotationStatus = "sent"
	QuotationStatusAccepted  QuotationStatus = "accepted"
	QuotationStatusDeclined  QuotationStatus = "declined"
	QuotationStatusExpired   QuotationStatus = "expired"
	QuotationStatusConverted QuotationStatus = "converted"
)

type Quotation struct {
	Id              uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId       uuid.UUID       `gorm:"type:uuid;not null;index;" json:"company_id"`
	QuotationNumber string          `gorm:"type:character varying;not null;" json:"quotation_number"`
	CustomerId      uuid.UUID       `gorm:"type:uuid;not null;index;" json:"customer_id"`
	Customer        Customer        `gorm:"foreignKey:CustomerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"customer,omitempty"`
	QuotationDate   time.Time       `gorm:"type:date;not null;" json:"quotation_date"`
	ExpiryDate      *time.Time      `gorm:"type:date;" json:"expiry_date,omitempty"`
	Subtotal        float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal   float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp             float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate         float64         `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount       float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal      float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	Status          QuotationStatus `gorm:"type:character varying;not null;default:'draft';" json:"status"`
	Notes           string          `gorm:"type:text;" json:"notes"`
	// Set when converted to invoice
	ConvertedInvoiceId *uuid.UUID      `gorm:"type:uuid;" json:"converted_invoice_id,omitempty"`
	CreatedBy          uuid.UUID       `gorm:"type:uuid;not null;" json:"created_by"`
	Items              []QuotationItem `gorm:"foreignKey:QuotationId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt          time.Time       `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Quotation) SearchableFields() []string {
	return []string{"quotation_number", "notes"}
}

func (Quotation) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("quotations.status IN ?", values)
	}
	if values, ok := filters["customer_id"]; ok {
		db = db.Where("quotations.customer_id IN ?", values)
	}
	return db
}

// ─── Quotation Item ───────────────────────────────────────────────────────────

func (QuotationItem) TableName() string { return "quotation_items" }

type QuotationItem struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	QuotationId   uuid.UUID `gorm:"type:uuid;not null;index;" json:"quotation_id"`
	Description   string    `gorm:"type:character varying;not null;" json:"description"`
	Qty           float64   `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price         float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount      float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable bool      `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount        float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	CreatedAt     time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

// ─── Invoice ──────────────────────────────────────────────────────────────────

func (Invoice) TableName() string { return "invoices" }

type InvoiceStatus string
type InvoiceTaxStatus string

const (
	// invoice_status workflow
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusConfirmed InvoiceStatus = "confirmed"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"

	// tax_status workflow (e-faktur)
	TaxStatusDraft      InvoiceTaxStatus = "draft"
	TaxStatusTaxable    InvoiceTaxStatus = "taxable"
	TaxStatusNotTaxable InvoiceTaxStatus = "not_taxable"
	TaxStatusReported   InvoiceTaxStatus = "reported"
	TaxStatusCancelled  InvoiceTaxStatus = "cancelled"
)

type Invoice struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId     uuid.UUID `gorm:"type:uuid;not null;index;" json:"company_id"`
	InvoiceNumber string    `gorm:"type:character varying;not null;" json:"invoice_number"`
	// Source quotation — nullable (invoice can be created standalone)
	QuotationId   *uuid.UUID       `gorm:"type:uuid;" json:"quotation_id,omitempty"`
	Quotation     *Quotation       `gorm:"foreignKey:QuotationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"quotation,omitempty"`
	CustomerId    uuid.UUID        `gorm:"type:uuid;not null;index;" json:"customer_id"`
	Customer      Customer         `gorm:"foreignKey:CustomerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"customer,omitempty"`
	InvoiceDate   time.Time        `gorm:"type:date;not null;" json:"invoice_date"`
	DueDate       time.Time        `gorm:"type:date;not null;" json:"due_date"`
	Subtotal      float64          `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal float64          `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp           float64          `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate       float64          `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount     float64          `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal    float64          `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	InvoiceStatus InvoiceStatus    `gorm:"type:character varying;not null;default:'draft';" json:"invoice_status"`
	TaxStatus     InvoiceTaxStatus `gorm:"type:character varying;not null;default:'draft';" json:"tax_status"`
	Notes         string           `gorm:"type:text;" json:"notes"`
	// Journal link — set on Confirm
	JournalEntryId *uuid.UUID    `gorm:"type:uuid;" json:"journal_entry_id,omitempty"`
	JournalEntry   *JournalEntry `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"journal_entry,omitempty"`
	// Payment journal — set on MarkPaid
	PaymentJournalId *uuid.UUID    `gorm:"type:uuid;" json:"payment_journal_id,omitempty"`
	PaymentJournal   *JournalEntry `gorm:"foreignKey:PaymentJournalId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"payment_journal,omitempty"`
	CreatedBy        uuid.UUID     `gorm:"type:uuid;not null;" json:"created_by"`
	Items            []InvoiceItem `gorm:"foreignKey:InvoiceId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt        time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Invoice) SearchableFields() []string {
	return []string{"invoice_number", "notes"}
}

func (Invoice) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["invoice_status"]; ok {
		db = db.Where("invoices.invoice_status IN ?", values)
	}
	if values, ok := filters["tax_status"]; ok {
		db = db.Where("invoices.tax_status IN ?", values)
	}
	if values, ok := filters["customer_id"]; ok {
		db = db.Where("invoices.customer_id IN ?", values)
	}
	return db
}

// ─── Invoice Item ─────────────────────────────────────────────────────────────

func (InvoiceItem) TableName() string { return "invoice_items" }

type InvoiceItem struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	InvoiceId     uuid.UUID `gorm:"type:uuid;not null;index;" json:"invoice_id"`
	Description   string    `gorm:"type:character varying;not null;" json:"description"`
	Qty           float64   `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price         float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount      float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable bool      `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount        float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	CreatedAt     time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}
