package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Customer) TableName() string {
	return "customers"
}

type Customer struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	Code      string     `gorm:"type:character varying;not null;unique;" json:"code"`
	Name      string     `gorm:"type:character varying;not null;" json:"name"`
	Email     string     `gorm:"type:character varying;" json:"email"`
	Phone     string     `gorm:"type:character varying;" json:"phone"`
	Address   string     `gorm:"type:text;" json:"address"`
	CoaId     *uuid.UUID `gorm:"type:uuid;" json:"coa_id"`
	COA       *COA       `gorm:"foreignKey:CoaId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"coa,omitempty"`
	Status    int        `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Customer) SearchableFields() []string {
	return []string{"code", "name", "email", "phone"}
}

// ApplyFilters applies Customer-specific filter logic to the given GORM query.
// Supports multi-value filters via repeated params (e.g. ?status=1&status=2).
func (Customer) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}
	return db
}
