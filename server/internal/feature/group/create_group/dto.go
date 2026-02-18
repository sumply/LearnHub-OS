package create_group

import "github.com/google/uuid"

type Input struct {
	Name string `json:"name"`

	CuratorID  *uuid.UUID `json:"curator_id"`
	StudentIDs uuid.UUIDs `json:"student_ids"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
