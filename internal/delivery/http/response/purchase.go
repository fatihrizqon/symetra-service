package response

import (
	"time"

	"github.com/google/uuid"
)

// ─── Purchase Order ───────────────────────────────────────────────────────────

type POItemResponse struct {
	Id              uuid.UUID `json:"id"`
	PurchaseOrderId uuid.UUID `json:"purchase_order_id"`
	Description     string    `json:"description"`
	Qty             float64   `json:"qty"`
	Price           float64   `json:"price"`
	Discount        float64   `json:"discount"`
	TaxApplicable   bool      `json:"tax_applicable"`
	Amount          float64   `json:"amount"`
}

type POResponse struct {
	Id              uuid.UUID        `json:"id"`
	CompanyId       uuid.UUID        `json:"company_id"`
	PONumber        string           `json:"po_number"`
	VendorId        uuid.UUID        `json:"vendor_id"`
	VendorName      string           `json:"vendor_name"`
	PODate          time.Time        `json:"po_date"`
	ExpiryDate      *time.Time       `json:"expiry_date,omitempty"`
	Subtotal        float64          `json:"subtotal"`
	DiscountTotal   float64          `json:"discount_total"`
	Dpp             float64          `json:"dpp"`
	TaxRate         float64          `json:"tax_rate"`
	TaxAmount       float64          `json:"tax_amount"`
	GrandTotal      float64          `json:"grand_total"`
	Status          string           `json:"status"`
	Notes           string           `json:"notes"`
	ConvertedBillId *uuid.UUID       `json:"converted_bill_id,omitempty"`
	CreatedBy       uuid.UUID        `json:"created_by"`
	Items           []POItemResponse `json:"items,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// ─── Bill ─────────────────────────────────────────────────────────────────────

type BillItemResponse struct {
	Id            uuid.UUID  `json:"id"`
	BillId        uuid.UUID  `json:"bill_id"`
	Description   string     `json:"description"`
	Qty           float64    `json:"qty"`
	Price         float64    `json:"price"`
	Discount      float64    `json:"discount"`
	TaxApplicable bool       `json:"tax_applicable"`
	Amount        float64    `json:"amount"`
	AccountId     *uuid.UUID `json:"account_id,omitempty"`
	AccountName   string     `json:"account_name,omitempty"`
}

type BillPaymentResponse struct {
	Id               uuid.UUID  `json:"id"`
	BillId           uuid.UUID  `json:"bill_id"`
	Amount           float64    `json:"amount"`
	PaymentDate      time.Time  `json:"payment_date"`
	PaymentAccountId uuid.UUID  `json:"payment_account_id"`
	PaymentAccountName string   `json:"payment_account_name,omitempty"`
	JournalEntryId   *uuid.UUID `json:"journal_entry_id,omitempty"`
	Notes            string     `json:"notes"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
}

type BillResponse struct {
	Id              uuid.UUID             `json:"id"`
	CompanyId       uuid.UUID             `json:"company_id"`
	BillNumber      string                `json:"bill_number"`
	PurchaseOrderId *uuid.UUID            `json:"purchase_order_id,omitempty"`
	VendorId        uuid.UUID             `json:"vendor_id"`
	VendorName      string                `json:"vendor_name"`
	BillDate        time.Time             `json:"bill_date"`
	DueDate         time.Time             `json:"due_date"`
	Subtotal        float64               `json:"subtotal"`
	DiscountTotal   float64               `json:"discount_total"`
	Dpp             float64               `json:"dpp"`
	TaxRate         float64               `json:"tax_rate"`
	TaxAmount       float64               `json:"tax_amount"`
	GrandTotal      float64               `json:"grand_total"`
	AmountPaid      float64               `json:"amount_paid"`
	AmountDue       float64               `json:"amount_due"`
	BillStatus      string                `json:"bill_status"`
	PaymentStatus   string                `json:"payment_status"`
	Notes           string                `json:"notes"`
	JournalEntryId  *uuid.UUID            `json:"journal_entry_id,omitempty"`
	CreatedBy       uuid.UUID             `json:"created_by"`
	Items           []BillItemResponse    `json:"items,omitempty"`
	Payments        []BillPaymentResponse `json:"payments,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}
