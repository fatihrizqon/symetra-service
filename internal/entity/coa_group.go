package entity

import (
	"time"

	"github.com/google/uuid"
)

func (COAGroup) TableName() string {
	return "coa_groups"
}

type COAGroup struct {
	Id        uuid.UUID     `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	Code      string        `gorm:"type:character varying; not null; unique;" json:"code"`
	Name      string        `gorm:"type:character varying; not null;" json:"name"`
	Status    int           `gorm:"type:int; not null; default:1;" json:"status"`
	SubGroups []COASubGroup `gorm:"foreignKey:GroupId"`
	CreatedAt time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COAGroup) SearchableFields() []string {
	return []string{"name", "code"}
}

type COAGroupFilters struct {
	Status *string
}
