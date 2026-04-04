package request

import (
	"github.com/google/uuid"
)

type JournalLineRequest struct {
	CoaId       uuid.UUID `validate:"required" json:"coa_id"`
	Description string    `validate:"required,min=1" json:"description"`
	Debit       float64   `validate:"min=0" json:"debit"`
	Credit      float64   `validate:"min=0" json:"credit"`
}

type JournalEntryCreateRequest struct {
	Date        string               `validate:"required" json:"date"`
	Description string               `validate:"required,min=1" json:"description"`
	Lines       []JournalLineRequest `validate:"required,min=2,dive" json:"lines"`
}

type JournalEntryUpdateRequest struct {
	Id          uuid.UUID            `json:"-"`
	Date        string               `validate:"required" json:"date"`
	Description string               `validate:"required,min=1" json:"description"`
	Lines       []JournalLineRequest `validate:"required,min=2,dive" json:"lines"`
}
