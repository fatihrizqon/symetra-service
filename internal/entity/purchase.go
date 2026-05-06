package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Purchase Order ───────────────────────────────────────────────────────────

func (PurchaseOrder) TableName() string { return "purchase_orders" }

type PurchaseOrderStatus string

const (
	POStatusDraft     PurchaseOrderStatus = "draft"
	POStatusSent      PurchaseOrderStatus = "sent"
	POStatusApproved  PurchaseOrderStatus = "approved"
	POStatusDeclined  PurchaseOrderStatus = "declined"
	POStatusExpired   PurchaseOrderStatus = "expired"
	POStatusConverted PurchaseOrderStatus = "converted"
)

type PurchaseOrder struct {
	Id              uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId       uuid.UUID           `gorm:"type:uuid;not null;index;" json:"company_id"`
	PONumber        string              `gorm:"type:character varying;not null;uniqueIndex:idx_company_po_number;" json:"po_number"`
	VendorId        uuid.UUID           `gorm:"type:uuid;not null;index;" json:"vendor_id"`
	Vendor          Vendor              `gorm:"foreignKey:VendorId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vendor,omitempty"`
	PODate          time.Time           `gorm:"type:date;not null;" json:"po_date"`
	ExpiryDate      *time.Time          `gorm:"type:date;" json:"expiry_date,omitempty"`
	Subtotal        float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal   float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp             float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate         float64             `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount       float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal      float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	Status          PurchaseOrderStatus `gorm:"type:character varying;not null;default:'draft';" json:"status"`
	Notes           string              `gorm:"type:text;" json:"notes"`
	ConvertedBillId *uuid.UUID          `gorm:"type:uuid;" json:"converted_bill_id,omitempty"`
	CreatedBy       uuid.UUID           `gorm:"type:uuid;not null;" json:"created_by"`
	Items           []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt       time.Time           `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (PurchaseOrder) SearchableFields() []string {
	return []string{"po_number", "notes"}
}

func (PurchaseOrder) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("purchase_orders.status IN ?", values)
	}
	if values, ok := filters["vendor_id"]; ok {
		db = db.Where("purchase_orders.vendor_id IN ?", values)
	}
	return db
}

// ─── Purchase Order Item ──────────────────────────────────────────────────────

func (PurchaseOrderItem) TableName() string { return "purchase_order_items" }

type PurchaseOrderItem struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	PurchaseOrderId uuid.UUID `gorm:"type:uuid;not null;index;" json:"purchase_order_id"`
	Description     string    `gorm:"type:character varying;not null;" json:"description"`
	Qty             float64   `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price           float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount        float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable   bool      `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount          float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	CreatedAt       time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

// ─── Bill ─────────────────────────────────────────────────────────────────────

func (Bill) TableName() string { return "bills" }

type BillStatus string
type BillPaymentStatus string

const (
	BillStatusDraft     BillStatus = "draft"
	BillStatusConfirmed BillStatus = "confirmed"
	BillStatusCancelled BillStatus = "cancelled"

	BillPaymentUnpaid  BillPaymentStatus = "unpaid"
	BillPaymentPartial BillPaymentStatus = "partial"
	BillPaymentPaid    BillPaymentStatus = "paid"
)

type Bill struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID `gorm:"type:uuid;not null;index;" json:"company_id"`
	// FIX [CFG-01]: Tambah uniqueIndex composite company_id+bill_number
	BillNumber      string            `gorm:"type:character varying;not null;uniqueIndex:idx_company_bill_number;" json:"bill_number"`
	PurchaseOrderId *uuid.UUID        `gorm:"type:uuid;" json:"purchase_order_id,omitempty"`
	PurchaseOrder   *PurchaseOrder    `gorm:"foreignKey:PurchaseOrderId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"purchase_order,omitempty"`
	VendorId        uuid.UUID         `gorm:"type:uuid;not null;index;" json:"vendor_id"`
	Vendor          Vendor            `gorm:"foreignKey:VendorId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vendor,omitempty"`
	BillDate        time.Time         `gorm:"type:date;not null;" json:"bill_date"`
	DueDate         time.Time         `gorm:"type:date;not null;" json:"due_date"`
	Subtotal        float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal   float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp             float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate         float64           `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount       float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal      float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	AmountPaid      float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_paid"`
	AmountDue       float64           `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_due"`
	BillStatus      BillStatus        `gorm:"type:character varying;not null;default:'draft';" json:"bill_status"`
	PaymentStatus   BillPaymentStatus `gorm:"type:character varying;not null;default:'unpaid';" json:"payment_status"`
	Notes           string            `gorm:"type:text;" json:"notes"`
	JournalEntryId  *uuid.UUID        `gorm:"type:uuid;" json:"journal_entry_id,omitempty"`
	JournalEntry    *JournalEntry     `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"journal_entry,omitempty"`
	CreatedBy       uuid.UUID         `gorm:"type:uuid;not null;" json:"created_by"`
	Items           []BillItem        `gorm:"foreignKey:BillId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	Payments        []BillPayment     `gorm:"foreignKey:BillId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payments,omitempty"`
	CreatedAt       time.Time         `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Bill) SearchableFields() []string {
	return []string{"bill_number", "notes"}
}

func (Bill) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["bill_status"]; ok {
		db = db.Where("bills.bill_status IN ?", values)
	}
	if values, ok := filters["payment_status"]; ok {
		db = db.Where("bills.payment_status IN ?", values)
	}
	if values, ok := filters["vendor_id"]; ok {
		db = db.Where("bills.vendor_id IN ?", values)
	}
	return db
}

// ─── Bill Item ────────────────────────────────────────────────────────────────

func (BillItem) TableName() string { return "bill_items" }

type BillItem struct {
	Id            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	BillId        uuid.UUID  `gorm:"type:uuid;not null;index;" json:"bill_id"`
	Description   string     `gorm:"type:character varying;not null;" json:"description"`
	Qty           float64    `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price         float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount      float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable bool       `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount        float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	AccountId     *uuid.UUID `gorm:"type:uuid;" json:"account_id,omitempty"`
	Account       *COA       `gorm:"foreignKey:AccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}

// ─── Bill Payment ─────────────────────────────────────────────────────────────

func (BillPayment) TableName() string { return "bill_payments" }

type BillPayment struct {
	Id               uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	BillId           uuid.UUID     `gorm:"type:uuid;not null;index;" json:"bill_id"`
	Amount           float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	PaymentDate      time.Time     `gorm:"type:date;not null;" json:"payment_date"`
	PaymentAccountId uuid.UUID     `gorm:"type:uuid;not null;" json:"payment_account_id"`
	PaymentAccount   *COA          `gorm:"foreignKey:PaymentAccountId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"payment_account,omitempty"`
	JournalEntryId   *uuid.UUID    `gorm:"type:uuid;" json:"journal_entry_id,omitempty"`
	JournalEntry     *JournalEntry `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"journal_entry,omitempty"`
	Notes            string        `gorm:"type:text;" json:"notes"`
	CreatedBy        uuid.UUID     `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt        time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}
