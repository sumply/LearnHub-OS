package get_by_id

import "server/internal/query"

type Output struct {
	Quiz query.Quiz `json:"quiz"`
}
