package repository

import (
	"context"

	"github.com/google/uuid"
)

type Saver[T any] interface {
	Save(context.Context, T) error
}

type Geter[T any] interface {
	Get(context.Context, uuid.UUID) (T, error)
}

type Lister[T any] interface {
	List(context.Context, []uuid.UUID) ([]T, error)
}

type Updater[T any] interface {
	Update(context.Context, T) error
}

type Remover[T any] interface {
	Remove(context.Context, uuid.UUID) error
}
