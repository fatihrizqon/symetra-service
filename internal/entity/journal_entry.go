package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (JournalEntry) TableName() string { return "journal_entries" }

type JournalStatus string

const (
	JournalStatusDraft  JournalStatus = "draft"
	JournalStatusPosted JournalStatus = "posted"
	JournalStatusVoid   JournalStatus = "void"
)

type JournalType string

const (
	JournalTypeGeneral JournalType = "general"
	JournalTypeRevenue JournalType = "revenue"
	JournalTypeExpense JournalType = "expense"
)

var JournalNumberPrefix = map[JournalType]string{
	JournalTypeGeneral: "JE",
	JournalTypeRevenue: "RV",
	JournalTypeExpense: "EX",
}

type JournalEntry struct {
	Id             uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId      uuid.UUID     `gorm:"type:uuid;not null;index;" json:"company_id"`        // ← NEW
	FiscalPeriodId *uuid.UUID    `gorm:"type:uuid;index;" json:"fiscal_period_id,omitempty"` // ← NEW
	JournalNumber  string        `gorm:"type:character varying;not null;" json:"journal_number"`
	Type           JournalType   `gorm:"type:character varying;not null;default:'general';" json:"type"`
	Date           time.Time     `gorm:"type:date;not null;" json:"date"`
	Description    string        `gorm:"type:text;not null;" json:"description"`
	Status         JournalStatus `gorm:"type:character varying;not null;default:'draft';" json:"status"`
	TotalDebit     float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"total_debit"`
	TotalCredit    float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"total_credit"`
	CreatedBy      uuid.UUID     `gorm:"type:uuid;" json:"created_by"`
	Lines          []JournalLine `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lines,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (JournalEntry) SearchableFields() []string {
	return []string{"journal_number", "description"}
}

func (JournalEntry) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}
	if values, ok := filters["type"]; ok {
		db = db.Where("type IN ?", values)
	}
	return db
}
