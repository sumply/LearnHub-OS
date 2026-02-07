package get_users

import "server/internal/dto"

type Output struct {
	Users []dto.UserLastAttempt `json:"users"`
}
