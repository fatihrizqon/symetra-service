package entity

import (
	"time"

	"github.com/google/uuid"
)

func (COASubGroup) TableName() string {
	return "coa_subgroups"
}

type COASubGroup struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	GroupId   uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	Group     COAGroup  `gorm:"foreignKey:GroupId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"group,omitempty"`
	COA       []COA     `gorm:"foreignKey:SubgroupId"`
	Code      string
	Name      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (COASubGroup) SearchableFields() []string {
	return []string{"name", "code"}
}

type COASubGroupFilters struct {
	Status *string
}
