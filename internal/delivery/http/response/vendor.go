package response

import (
	"time"

	"github.com/google/uuid"
)

type VendorResponse struct {
	Id        uuid.UUID  `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Address   string     `json:"address"`
	CoaId     *uuid.UUID `json:"coa_id,omitempty"`
	CoaCode   string     `json:"coa_code,omitempty"`
	CoaName   string     `json:"coa_name,omitempty"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
