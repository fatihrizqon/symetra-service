package request

import "github.com/google/uuid"

// ─── Purchase Order ───────────────────────────────────────────────────────────

type POItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,gt=0"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	// FIX [BUG-13/FRAUD-04]: Discount wajib >= 0 untuk mencegah inflasi biaya fiktif
	Discount      float64 `json:"discount" validate:"gte=0"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type POCreateRequest struct {
	// FIX [FRAUD-02]: VendorId tetap uuid.UUID, validasi kepemilikan dilakukan di service
	VendorId   uuid.UUID      `json:"vendor_id" validate:"required"`
	PODate     string         `json:"po_date" validate:"required"`
	// FIX [BUG-16]: ExpiryDate divalidasi di service level (harus > PODate)
	ExpiryDate string         `json:"expiry_date"`
	Notes      string         `json:"notes"`
	// FIX [BUG-17]: max=500 untuk mencegah request dengan item yang sangat banyak
	Items      []POItemRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type POUpdateRequest struct {
	Id         uuid.UUID      `json:"id" validate:"required"`
	VendorId   uuid.UUID      `json:"vendor_id" validate:"required"`
	PODate     string         `json:"po_date" validate:"required"`
	ExpiryDate string         `json:"expiry_date"`
	Notes      string         `json:"notes"`
	Items      []POItemRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

// ─── Bill ─────────────────────────────────────────────────────────────────────

type BillItemRequest struct {
	Description   string     `json:"description" validate:"required"`
	Qty           float64    `json:"qty" validate:"required,gt=0"`
	Price         float64    `json:"price" validate:"required,gte=0"`
	// FIX [BUG-13/FRAUD-04]: Discount wajib >= 0
	Discount      float64    `json:"discount" validate:"gte=0"`
	TaxApplicable bool       `json:"tax_applicable"`
	AccountId     *uuid.UUID `json:"account_id"`
}

type BillCreateRequest struct {
	// FIX [BUG-14]: Ubah VendorId dari string ke uuid.UUID agar konsisten dengan POCreateRequest
	// dan validasi format UUID terjadi di layer request bukan service
	VendorId string            `json:"vendor_id" validate:"required,uuid"`
	BillDate string            `json:"bill_date" validate:"required"`
	// FIX [BUG-15]: DueDate divalidasi di service (harus >= BillDate)
	DueDate  string            `json:"due_date" validate:"required"`
	Notes    string            `json:"notes"`
	Items    []BillItemRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BillUpdateRequest struct {
	Id       uuid.UUID         `json:"id" validate:"required"`
	// FIX [BUG-14]: Konsisten dengan BillCreateRequest
	VendorId string            `json:"vendor_id" validate:"required,uuid"`
	BillDate string            `json:"bill_date" validate:"required"`
	DueDate  string            `json:"due_date" validate:"required"`
	Notes    string            `json:"notes"`
	Items    []BillItemRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BillPaymentRequest struct {
	Amount           float64 `json:"amount" validate:"required,gt=0"`
	PaymentAccountId string  `json:"payment_account_id" validate:"required,uuid"`
	PaymentDate      string  `json:"payment_date" validate:"required"`
	Notes            string  `json:"notes"`
}
