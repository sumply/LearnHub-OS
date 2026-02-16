package get_by_id

import (
	"server/internal/dto"
)

type Output struct {
	User dto.User `json:"user"`
}
