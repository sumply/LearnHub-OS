package validator

import "github.com/go-playground/validator/v10"

var V *validator.Validate

func init() {
	V = validator.New()
	V.RegisterAlias("password", "min=8,max=16")
}
