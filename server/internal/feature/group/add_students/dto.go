package add_students

import (
	"net/http"
	"server/internal/pkg/decoder"

	"github.com/google/uuid"
)

type Input struct {
	GroupID    uuid.UUID  `json:"-"`
	StudentIDs uuid.UUIDs `json:"student_ids"`
}

func InputFromRequest(r *http.Request) (Input, error) {
	return decoder.JSON[Input](r.Body)
}
