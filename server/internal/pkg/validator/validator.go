package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func V(ctx context.Context, s any) error {
	return v.StructCtx(ctx, s)
}

func init() {
	v = validator.New()
	v.RegisterAlias("password", "min=8,max=16")
	v.RegisterAlias("user-name", "min=2,max=32,alphanumunicode")
	v.RegisterAlias("subject-name", "min=2,max=100")
	v.RegisterAlias("group-name", "min=2,max=30")
}
