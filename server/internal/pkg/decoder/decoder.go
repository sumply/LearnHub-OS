package decoder

import (
	"encoding/json"
	"fmt"
	"io"
)

func JSON[T any](r io.Reader) (T, error) {
	var tt T
	fmt.Println(tt)
	if err := json.NewDecoder(r).Decode(&tt); err != nil {
		return tt, err
	}
	return tt, nil
}
