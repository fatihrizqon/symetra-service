package entity

import (
	"time"

	"github.com/google/uuid"
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
	return []string{"name", "code"}
}

type COAFilters struct {
	Status *string
}
