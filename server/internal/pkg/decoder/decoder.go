package decoder

import (
	"encoding/json"
	"io"
)

func JSON[T any](r io.Reader) (T, error) {
	tt := *new(T)
	if err := json.NewDecoder(r).Decode(tt); err != nil {
		return tt, err
	}
	return tt, nil
}
