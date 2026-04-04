package request

import "github.com/google/uuid"

// ─── Company Configuration ────────────────────────────────────────────────────

type CompanyConfigurationRequest struct {
	EnableTax bool    `json:"enable_tax"`
	TaxRate   float64 `json:"tax_rate"`

	ArAccountId             *uuid.UUID `json:"ar_account_id"`
	ApAccountId             *uuid.UUID `json:"ap_account_id"`
	SalesRevenueAccountId   *uuid.UUID `json:"sales_revenue_account_id"`
	ServiceRevenueAccountId *uuid.UUID `json:"service_revenue_account_id"`
	TaxPayableAccountId     *uuid.UUID `json:"tax_payable_account_id"`
	TaxReceivableAccountId  *uuid.UUID `json:"tax_receivable_account_id"`
	BankAccountId           *uuid.UUID `json:"bank_account_id"`
	CashAccountId           *uuid.UUID `json:"cash_account_id"`

	InvoicePrefix   string `json:"invoice_prefix"`
	QuotationPrefix string `json:"quotation_prefix"`
	InvoiceDueDays  int    `json:"invoice_due_days"`

	DefaultExpenseAccountId *uuid.UUID `json:"default_expense_account_id"`
	PurchaseOrderPrefix     string     `json:"purchase_order_prefix"`
	BillPrefix              string     `json:"bill_prefix"`
	BillDueDays             int        `json:"bill_due_days"`
}

// ─── Quotation ────────────────────────────────────────────────────────────────

type QuotationItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,gt=0"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	Discount      float64 `json:"discount"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type QuotationCreateRequest struct {
	CustomerId    uuid.UUID              `json:"customer_id" validate:"required"`
	QuotationDate string                 `json:"quotation_date" validate:"required"`
	ExpiryDate    string                 `json:"expiry_date"`
	Notes         string                 `json:"notes"`
	Items         []QuotationItemRequest `json:"items" validate:"required,min=1,dive"`
}

type QuotationUpdateRequest struct {
	Id            uuid.UUID              `json:"id" validate:"required"`
	CustomerId    uuid.UUID              `json:"customer_id" validate:"required"`
	QuotationDate string                 `json:"quotation_date" validate:"required"`
	ExpiryDate    string                 `json:"expiry_date"`
	Notes         string                 `json:"notes"`
	Items         []QuotationItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ─── Invoice ──────────────────────────────────────────────────────────────────

type InvoiceItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,gt=0"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	Discount      float64 `json:"discount"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type InvoiceCreateRequest struct {
	CustomerId  uuid.UUID            `json:"customer_id" validate:"required"`
	InvoiceDate string               `json:"invoice_date" validate:"required"`
	DueDate     string               `json:"due_date" validate:"required"`
	Notes       string               `json:"notes"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

type InvoiceUpdateRequest struct {
	Id          uuid.UUID            `json:"id" validate:"required"`
	CustomerId  uuid.UUID            `json:"customer_id" validate:"required"`
	InvoiceDate string               `json:"invoice_date" validate:"required"`
	DueDate     string               `json:"due_date" validate:"required"`
	Notes       string               `json:"notes"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

// MarkPaidRequest — specify which account received the payment
type InvoiceMarkPaidRequest struct {
	PaymentAccountId uuid.UUID `json:"payment_account_id" validate:"required"` // bank or cash COA
	PaymentDate      string    `json:"payment_date" validate:"required"`
}
