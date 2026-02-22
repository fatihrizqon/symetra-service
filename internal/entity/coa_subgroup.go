package entity

import (
	"time"

	"github.com/google/uuid"
)

func (COASubGroup) TableName() string {
	return "coa_subgroups"
}

type COASubGroup struct {
	Id        uuid.UUID `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	GroupId   uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	Group     COAGroup  `gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"group,omitempty"`
	Code      string    `gorm:"type:character varying; not null; unique;" json:"code"`
	Name      string    `gorm:"type:character varying; not null;" json:"name"`
	Status    int       `gorm:"type:int; not null; default:1;" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COASubGroup) SearchableFields() []string {
	return []string{"name", "code"}
}

type COASubGroupFilters struct {
	Status *string
}
