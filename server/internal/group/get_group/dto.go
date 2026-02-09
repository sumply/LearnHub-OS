package get_group

import "github.com/google/uuid"

type OutputGroup struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Output struct {
	Groups []OutputGroup
}
