package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// ApplyFilters applies COASubGroup-specific filter logic to the given GORM query.
// Supports multi-value filters via repeated params (e.g. ?status=1&status=2).
func (COASubGroup) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}
	return db
}
