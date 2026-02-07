package utils

import (
	"fmt"

	"github.com/doug-martin/goqu/v9/exec"
)

type DTOMapper[dtoType any] interface {
	DTO() dtoType
}

func ScanDTO[rowType, dtoType any](s exec.Scanner) ([]dtoType, error) {
	items := make([]dtoType, 0)

	for s.Next() {
		var row rowType

		if err := s.ScanStruct(&row); err != nil {
			return nil, err
		}

		if rower, ok := any(&row).(DTOMapper[dtoType]); ok {
			items = append(items, rower.DTO())
		} else {
			return nil, fmt.Errorf("type %T does not implement DTOMapper[%T]", row, *new(dtoType))
		}
	}

	return items, nil
}
