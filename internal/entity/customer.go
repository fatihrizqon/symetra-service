package entity

import (
	"time"

	"github.com/google/uuid"
)

func (Customer) TableName() string {
	return "customers"
}

type Customer struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	Code      string    `gorm:"type:character varying;not null;unique;" json:"code"`
	Name      string    `gorm:"type:character varying;not null;" json:"name"`
	Email     string    `gorm:"type:character varying;" json:"email"`
	Phone     string    `gorm:"type:character varying;" json:"phone"`
	Address   string    `gorm:"type:text;" json:"address"`
	CoaId     *uuid.UUID `gorm:"type:uuid;" json:"coa_id"`
	COA       *COA      `gorm:"foreignKey:CoaId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"coa,omitempty"`
	Status    int       `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Customer) SearchableFields() []string {
	return []string{"code", "name", "email", "phone"}
}

type CustomerFilters struct {
	Status *string
}
