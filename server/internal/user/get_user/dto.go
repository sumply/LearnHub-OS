package get_user

import (
	"github.com/google/uuid"
)

type OutputUser struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Role      string     `json:"role"`
	Groups    uuid.UUIDs `json:"groups,omitempty"`
	Group     *uuid.UUID `json:"group,omitempty"`
}

type Output struct {
	Users []OutputUser `json:"users"`
}
