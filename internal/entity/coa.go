package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (COA) TableName() string {
	return "chart_of_accounts"
}

type COA struct {
	Id         uuid.UUID   `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	SubgroupId uuid.UUID   `gorm:"type:uuid;not null;index" json:"subgroup_id"`
	SubGroup   COASubGroup `gorm:"foreignKey:SubgroupId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"subgroup,omitempty"`
	Code       string      `gorm:"type:character varying; not null; unique;" json:"code"`
	Name       string      `gorm:"type:character varying; not null;" json:"name"`
	Status     int         `gorm:"type:int; not null; default:1;" json:"status"`
	CreatedAt  time.Time   `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COA) SearchableFields() []string {
	return []string{"chart_of_accounts.name", "chart_of_accounts.code"}
}

// ApplyFilters applies COA-specific filter logic to the given GORM query.
// group_type filter performs a JOIN to coa_groups to filter by group type
// (e.g. ?group_type=revenue returns only COAs belonging to revenue groups).
func (COA) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("chart_of_accounts.status IN ?", values)
	}

	if values, ok := filters["group_type"]; ok {
		db = db.
			Joins("JOIN coa_subgroups ON coa_subgroups.id = chart_of_accounts.subgroup_id").
			Joins("JOIN coa_groups ON coa_groups.id = coa_subgroups.group_id").
			Where("coa_groups.type IN ?", values)
	}

	return db
}
