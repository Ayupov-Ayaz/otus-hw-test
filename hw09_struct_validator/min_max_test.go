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
			errors[i] = NewValidateError("Id", err[i])
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
			checkErr: systemError(ErrInvalidMinMaxCommandValueType),
			obj: struct {
				Id string `validate:"min"`
			}{},
		},
		{
			name:     "parse tag failed",
			checkErr: systemError(ErrInvalidMinMaxCommand),
			obj: struct {
				Id int `validate:"min"`
			}{},
		},
		{
			name:     "min failed",
			checkErr: vErr(ErrInvalidMin),
			obj: struct {
				Id int `validate:"min:10"`
			}{},
		},
		{
			name:     "max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				Id int `validate:"max:10"`
			}{
				Id: 11,
			},
		},
		{
			name: "min success",
			obj: struct {
				Id int `validate:"min:10"`
			}{
				Id: 11,
			},
		},
		{
			name: "max success",
			obj: struct {
				Id int `validate:"max:10"`
			}{
				Id: 9,
			},
		},
		{
			name:     "min success, max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				Id int `validate:"min:1|max11"`
			}{
				Id: 12,
			},
		},
		{
			name:     "max success, min failed",
			checkErr: vErr(ErrInvalidMin),
			obj: struct {
				Id int `validate:"max:10|min:2"`
			}{
				Id: 1,
			},
		},
		{
			name:     "min, max in slice, min failed",
			checkErr: vErr(ErrInvalidMin),
			obj: struct {
				Id []int `validate:"min:10|max:20"`
			}{
				Id: []int{9, 20},
			},
		},
		{
			name:     "min, max in slice, max failed",
			checkErr: vErr(ErrInvalidMax),
			obj: struct {
				Id []int `validate:"min:10|max:20"`
			}{
				Id: []int{11, 21},
			},
		},
		{
			name:     "min and max value invalid",
			checkErr: vErr(ErrInvalidMin, ErrInvalidMax),
			obj: struct {
				Id []int `validate:"min:10|max:20"`
			}{
				Id: []int{0, 21},
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
