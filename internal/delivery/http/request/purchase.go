package request

import "github.com/google/uuid"

// ─── Purchase Order ───────────────────────────────────────────────────────────

type POItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,gt=0"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	Discount      float64 `json:"discount"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type POCreateRequest struct {
	VendorId   uuid.UUID      `json:"vendor_id" validate:"required"`
	PODate     string         `json:"po_date" validate:"required"`
	ExpiryDate string         `json:"expiry_date"`
	Notes      string         `json:"notes"`
	Items      []POItemRequest `json:"items" validate:"required,min=1,dive"`
}

type POUpdateRequest struct {
	Id         uuid.UUID      `json:"id" validate:"required"`
	VendorId   uuid.UUID      `json:"vendor_id" validate:"required"`
	PODate     string         `json:"po_date" validate:"required"`
	ExpiryDate string         `json:"expiry_date"`
	Notes      string         `json:"notes"`
	Items      []POItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ─── Bill ─────────────────────────────────────────────────────────────────────

type BillItemRequest struct {
	Description   string     `json:"description" validate:"required"`
	Qty           float64    `json:"qty" validate:"required,gt=0"`
	Price         float64    `json:"price" validate:"required,gte=0"`
	Discount      float64    `json:"discount"`
	TaxApplicable bool       `json:"tax_applicable"`
	AccountId     *uuid.UUID `json:"account_id"`
}

type BillCreateRequest struct {
	VendorId string            `json:"vendor_id" validate:"required"`
	BillDate string            `json:"bill_date" validate:"required"`
	DueDate  string            `json:"due_date" validate:"required"`
	Notes    string            `json:"notes"`
	Items    []BillItemRequest `json:"items" validate:"required,min=1,dive"`
}

type BillUpdateRequest struct {
	Id       uuid.UUID         `json:"id" validate:"required"`
	VendorId string            `json:"vendor_id" validate:"required"`
	BillDate string            `json:"bill_date" validate:"required"`
	DueDate  string            `json:"due_date" validate:"required"`
	Notes    string            `json:"notes"`
	Items    []BillItemRequest `json:"items" validate:"required,min=1,dive"`
}

type BillPaymentRequest struct {
	Amount           float64 `json:"amount" validate:"required,gt=0"`
	PaymentAccountId string  `json:"payment_account_id" validate:"required"`
	PaymentDate      string  `json:"payment_date" validate:"required"`
	Notes            string  `json:"notes"`
}
