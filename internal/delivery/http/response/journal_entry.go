package response

import (
	"time"

	"github.com/google/uuid"
)

type JournalLineResponse struct {
	Id             uuid.UUID `json:"id"`
	JournalEntryId uuid.UUID `json:"journal_entry_id"`
	CoaId          uuid.UUID `json:"coa_id"`
	CoaCode        string    `json:"coa_code"`
	CoaName        string    `json:"coa_name"`
	Description    string    `json:"description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
}

type JournalEntryResponse struct {
	Id            uuid.UUID             `json:"id"`
	JournalNumber string                `json:"journal_number"`
	Type          string                `json:"type"`
	Date          time.Time             `json:"date"`
	Description   string                `json:"description"`
	Status        string                `json:"status"`
	TotalDebit    float64               `json:"total_debit"`
	TotalCredit   float64               `json:"total_credit"`
	CreatedBy     uuid.UUID             `json:"created_by"`
	Lines         []JournalLineResponse `json:"lines,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}
