package response

import (
	"time"

	"github.com/google/uuid"
)

// ─── Company Configuration ────────────────────────────────────────────────────

type CompanyConfigurationResponse struct {
	Id        uuid.UUID `json:"id"`
	CompanyId uuid.UUID `json:"company_id"`

	EnableTax bool    `json:"enable_tax"`
	TaxRate   float64 `json:"tax_rate"`

	ArAccountId               *uuid.UUID `json:"ar_account_id"`
	ArAccountName             string     `json:"ar_account_name,omitempty"`
	ApAccountId               *uuid.UUID `json:"ap_account_id"`
	ApAccountName             string     `json:"ap_account_name,omitempty"`
	SalesRevenueAccountId     *uuid.UUID `json:"sales_revenue_account_id"`
	SalesRevenueAccountName   string     `json:"sales_revenue_account_name,omitempty"`
	ServiceRevenueAccountId   *uuid.UUID `json:"service_revenue_account_id"`
	ServiceRevenueAccountName string     `json:"service_revenue_account_name,omitempty"`
	TaxPayableAccountId       *uuid.UUID `json:"tax_payable_account_id"`
	TaxPayableAccountName     string     `json:"tax_payable_account_name,omitempty"`
	TaxReceivableAccountId    *uuid.UUID `json:"tax_receivable_account_id"`
	TaxReceivableAccountName  string     `json:"tax_receivable_account_name,omitempty"`
	BankAccountId             *uuid.UUID `json:"bank_account_id"`
	BankAccountName           string     `json:"bank_account_name,omitempty"`
	CashAccountId             *uuid.UUID `json:"cash_account_id"`
	CashAccountName           string     `json:"cash_account_name,omitempty"`

	InvoicePrefix   string `json:"invoice_prefix"`
	QuotationPrefix string `json:"quotation_prefix"`
	InvoiceDueDays  int    `json:"invoice_due_days"`

	DefaultExpenseAccountId   *uuid.UUID `json:"default_expense_account_id"`
	DefaultExpenseAccountName string     `json:"default_expense_account_name,omitempty"`
	PurchaseOrderPrefix       string     `json:"purchase_order_prefix"`
	BillPrefix                string     `json:"bill_prefix"`
	BillDueDays               int        `json:"bill_due_days"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── Quotation ────────────────────────────────────────────────────────────────

type QuotationItemResponse struct {
	Id            uuid.UUID `json:"id"`
	QuotationId   uuid.UUID `json:"quotation_id"`
	Description   string    `json:"description"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	TaxApplicable bool      `json:"tax_applicable"`
	Amount        float64   `json:"amount"`
}

type QuotationResponse struct {
	Id                 uuid.UUID               `json:"id"`
	CompanyId          uuid.UUID               `json:"company_id"`
	QuotationNumber    string                  `json:"quotation_number"`
	CustomerId         uuid.UUID               `json:"customer_id"`
	CustomerName       string                  `json:"customer_name"`
	QuotationDate      time.Time               `json:"quotation_date"`
	ExpiryDate         *time.Time              `json:"expiry_date,omitempty"`
	Subtotal           float64                 `json:"subtotal"`
	DiscountTotal      float64                 `json:"discount_total"`
	Dpp                float64                 `json:"dpp"`
	TaxRate            float64                 `json:"tax_rate"`
	TaxAmount          float64                 `json:"tax_amount"`
	GrandTotal         float64                 `json:"grand_total"`
	Status             string                  `json:"status"`
	Notes              string                  `json:"notes"`
	ConvertedInvoiceId *uuid.UUID              `json:"converted_invoice_id,omitempty"`
	CreatedBy          uuid.UUID               `json:"created_by"`
	Items              []QuotationItemResponse `json:"items,omitempty"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

// ─── Invoice ──────────────────────────────────────────────────────────────────

type InvoiceItemResponse struct {
	Id            uuid.UUID `json:"id"`
	InvoiceId     uuid.UUID `json:"invoice_id"`
	Description   string    `json:"description"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	TaxApplicable bool      `json:"tax_applicable"`
	Amount        float64   `json:"amount"`
}

type InvoiceResponse struct {
	Id               uuid.UUID             `json:"id"`
	CompanyId        uuid.UUID             `json:"company_id"`
	InvoiceNumber    string                `json:"invoice_number"`
	QuotationId      *uuid.UUID            `json:"quotation_id,omitempty"`
	CustomerId       uuid.UUID             `json:"customer_id"`
	CustomerName     string                `json:"customer_name"`
	InvoiceDate      time.Time             `json:"invoice_date"`
	DueDate          time.Time             `json:"due_date"`
	Subtotal         float64               `json:"subtotal"`
	DiscountTotal    float64               `json:"discount_total"`
	Dpp              float64               `json:"dpp"`
	TaxRate          float64               `json:"tax_rate"`
	TaxAmount        float64               `json:"tax_amount"`
	GrandTotal       float64               `json:"grand_total"`
	InvoiceStatus    string                `json:"invoice_status"`
	TaxStatus        string                `json:"tax_status"`
	Notes            string                `json:"notes"`
	JournalEntryId   *uuid.UUID            `json:"journal_entry_id,omitempty"`
	PaymentJournalId *uuid.UUID            `json:"payment_journal_id,omitempty"`
	CreatedBy        uuid.UUID             `json:"created_by"`
	Items            []InvoiceItemResponse `json:"items,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}
