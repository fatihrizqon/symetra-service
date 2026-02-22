package response

import (
	"time"

	"github.com/google/uuid"
)

type COASubGroupResponse struct {
	Id        uuid.UUID `json:"id"`
	GroupId   uuid.UUID `json:"group_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
