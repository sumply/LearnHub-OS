package get_group

import "server/internal/dto"

type Output struct {
	Groups []dto.Group `json:"groups"`
}
