package entity

import (
	"time"

	"github.com/google/uuid"
)

func (JournalLine) TableName() string {
	return "journal_lines"
}

type JournalLine struct {
	Id             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	JournalEntryId uuid.UUID    `gorm:"type:uuid;not null;index;" json:"journal_entry_id"`
	CoaId          uuid.UUID    `gorm:"type:uuid;not null;index;" json:"coa_id"`
	COA            COA          `gorm:"foreignKey:CoaId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"coa,omitempty"`
	JournalEntry   JournalEntry `gorm:"foreignKey:JournalEntryId;" json:"-"`
	Description    string       `gorm:"type:text;" json:"description"`
	Debit          float64      `gorm:"type:numeric(20,4);not null;default:0;" json:"debit"`
	Credit         float64      `gorm:"type:numeric(20,4);not null;default:0;" json:"credit"`
	CreatedAt      time.Time    `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime;" json:"updated_at"`
}
