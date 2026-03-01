package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (COAGroup) TableName() string {
	return "coa_groups"
}

type COAGroup struct {
	Id            uuid.UUID     `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	Code          string        `gorm:"type:character varying; not null; unique;" json:"code"`
	Name          string        `gorm:"type:character varying; not null;" json:"name"`
	NormalBalance string        `gorm:"type:character varying; not null;" json:"normal_balance"`
	Status        int           `gorm:"type:int; not null; default:1;" json:"status"`
	SubGroups     []COASubGroup `gorm:"foreignKey:GroupId"`
	CreatedAt     time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COAGroup) SearchableFields() []string {
	return []string{"name", "code"}
}

// ApplyFilters applies COAGroup-specific filter logic to the given GORM query.
// Supports multi-value filters via repeated params (e.g. ?status=1&status=2).
func (COAGroup) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["code"]; ok {
		db = db.Where("code IN ?", values)
	}
	if values, ok := filters["normal_balance"]; ok {
		db = db.Where("normal_balance IN ?", values)
	}
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}
	return db
}
