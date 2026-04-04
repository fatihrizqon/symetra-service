package request

import "github.com/google/uuid"

type COAGroupCreateRequest struct {
	Code          string `validate:"required,min=1" json:"code"`
	Name          string `validate:"required,min=1" json:"name"`
	NormalBalance string `validate:"required,min=1" json:"normal_balance"`
}

type COAGroupUpdateRequest struct {
	Id            uuid.UUID
	Code          string `validate:"required,min=1,max=20" json:"code"`
	Name          string `validate:"required,min=1,max=20" json:"name"`
	NormalBalance string `validate:"required,min=1" json:"normal_balance"`
}
