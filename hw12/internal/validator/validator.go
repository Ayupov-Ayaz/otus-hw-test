package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

func (v *Validator) Validate(ctx context.Context, i interface{}) error {
	return v.v.StructCtx(ctx, i)
}
