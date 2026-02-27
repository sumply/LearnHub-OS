package encoder

import (
	"encoding/json"
	"io"
)

func JSON(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
