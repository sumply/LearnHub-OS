package utils

import (
	"github.com/doug-martin/goqu/v9/exec"
)

type Rower[dtoType any] interface {
	DTO() dtoType
	Clear()
}

func ScanDTO[dtoType any](s exec.Scanner, rower Rower[dtoType]) ([]dtoType, error) {
	items := make([]dtoType, 0)

	for s.Next() {
		rower.Clear()
		if err := s.ScanStruct(rower); err != nil {
			return nil, err
		}

		items = append(items, rower.DTO())
	}

	return items, nil
}
