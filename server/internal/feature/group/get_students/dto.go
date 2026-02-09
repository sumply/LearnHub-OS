package get_students

import "server/internal/dto"

type Output struct {
	Users []dto.User `json:"students"`
}
