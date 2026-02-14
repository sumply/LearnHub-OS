package filter

import "github.com/google/uuid"

type Null[T any] struct {
	V     T
	Valid bool
}

type Quiz struct {
	Owner       Null[uuid.UUID]
	User        Null[uuid.UUID]
	IsCompleted Null[bool]
}
