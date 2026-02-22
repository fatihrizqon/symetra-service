package entity

import (
	"time"

	"github.com/google/uuid"
)

func (JournalEntry) TableName() string {
	return "journal_entries"
}

// JournalStatus represents the lifecycle state of a journal entry.
// Allowed transitions: draft → posted → void
type JournalStatus string

const (
	JournalStatusDraft  JournalStatus = "draft"
	JournalStatusPosted JournalStatus = "posted"
	JournalStatusVoid   JournalStatus = "void"
)

type JournalEntry struct {
	Id            uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	JournalNumber string        `gorm:"type:character varying;not null;unique;" json:"journal_number"`
	Date          time.Time     `gorm:"type:date;not null;" json:"date"`
	Description   string        `gorm:"type:text;not null;" json:"description"`
	Status        JournalStatus `gorm:"type:character varying;not null;default:'draft';" json:"status"`
	TotalDebit    float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"total_debit"`
	TotalCredit   float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"total_credit"`
	CreatedBy     uuid.UUID     `gorm:"type:uuid;" json:"created_by"`
	Lines         []JournalLine `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lines,omitempty"`
	CreatedAt     time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (JournalEntry) SearchableFields() []string {
	return []string{"journal_number", "description"}
}

type JournalEntryFilters struct {
	Status *string
}
