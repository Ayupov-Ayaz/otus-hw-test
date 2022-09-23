package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation_min_max(t *testing.T) {
	vErr := func(err ...error) func(gotErr error) bool {
		errors := make([]ValidationError, len(err))
		for i := 0; i < len(err); i++ {
			errors[i] = NewValidateError("ID", err[i])
		}

		return validateError(errors...)
	}

	tests := []struct {
		name     string
		checkErr func(err error) bool
		obj      interface{}
	}{
		{
			name:     "invalid type",
			checkErr: systemError(ErrInvalidMinCommandValueType),
			obj: struct {
				ID string `validate:"min"`
			}{},
		},
		{
			name:     "parse tag failed",
			checkErr: systemError(ErrExtractMinMaxValueFailed),
			obj: struct {
				ID int `validate:"min"`
			}{},
		},
		{
			name:     "min failed",
			checkErr: vErr(ErrMinInvalid),
			obj: struct {
				ID int `validate:"min:10"`
			}{},
		},
		{
			name:     "max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				ID int `validate:"max:10"`
			}{
				ID: 11,
			},
		},
		{
			name: "min success",
			obj: struct {
				ID int `validate:"min:10"`
			}{
				ID: 11,
			},
		},
		{
			name: "max success",
			obj: struct {
				ID int `validate:"max:10"`
			}{
				ID: 9,
			},
		},
		{
			name:     "min success, max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				ID int `validate:"min:1|max11"`
			}{
				ID: 12,
			},
		},
		{
			name:     "max success, min failed",
			checkErr: vErr(ErrMinInvalid),
			obj: struct {
				ID int `validate:"max:10|min:2"`
			}{
				ID: 1,
			},
		},
		{
			name:     "min, max in slice, min failed",
			checkErr: vErr(ErrMinInvalid),
			obj: struct {
				ID []int `validate:"min:10|max:20"`
			}{
				ID: []int{9, 20},
			},
		},
		{
			name:     "min, max in slice, max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				ID []int `validate:"min:10|max:20"`
			}{
				ID: []int{11, 21},
			},
		},
		{
			name:     "min and max value invalid",
			checkErr: vErr(ErrMinInvalid, ErrInvalidMax),
			obj: struct {
				ID []int `validate:"min:10|max:20"`
			}{
				ID: []int{0, 21},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.obj)
			if tt.checkErr != nil {
				require.True(t, tt.checkErr(err))
			} else {
				require.Nil(t, err)
			}
		})
	}
}
