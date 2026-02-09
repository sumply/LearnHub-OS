package get_by_id

import "server/internal/dto"

type Output struct {
	Quiz dto.Quiz `json:"quiz"`
}
