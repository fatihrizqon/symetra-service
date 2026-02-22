package request

import "github.com/google/uuid"

type VendorCreateRequest struct {
	Code    string     `validate:"required,min=1,max=20" json:"code"`
	Name    string     `validate:"required,min=1" json:"name"`
	Email   string     `validate:"omitempty,email" json:"email"`
	Phone   string     `validate:"omitempty" json:"phone"`
	Address string     `validate:"omitempty" json:"address"`
	CoaId   *uuid.UUID `validate:"omitempty" json:"coa_id"`
}

type VendorUpdateRequest struct {
	Id      uuid.UUID  `json:"-"`
	Code    string     `validate:"required,min=1,max=20" json:"code"`
	Name    string     `validate:"required,min=1" json:"name"`
	Email   string     `validate:"omitempty,email" json:"email"`
	Phone   string     `validate:"omitempty" json:"phone"`
	Address string     `validate:"omitempty" json:"address"`
	CoaId   *uuid.UUID `validate:"omitempty" json:"coa_id"`
}
