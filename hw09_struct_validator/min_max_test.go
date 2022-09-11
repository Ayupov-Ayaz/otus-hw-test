package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation_min_max(t *testing.T) {

	tests := []struct {
		name   string
		expErr error
		obj    interface{}
	}{
		{
			name:   "invalid type",
			expErr: ErrInvalidMinMaxCommandValueType,
			obj: struct {
				Id string `validate:"min"`
			}{},
		},
		{
			name:   "parse tag failed",
			expErr: ErrInvalidMinMaxCommand,
			obj: struct {
				Id int `validate:"min"`
			}{},
		},
		{
			name:   "min failed",
			expErr: ErrInvalidMin,
			obj: struct {
				Id int `validate:"min:10"`
			}{},
		},
		{
			name:   "max failed",
			expErr: ErrInvalidMax,
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
			name:   "min success, max failed",
			expErr: ErrInvalidMax,
			obj: struct {
				Id int `validate:"min:1|max11"`
			}{
				Id: 12,
			},
		},
		{
			name:   "max success, min failed",
			expErr: ErrInvalidMin,
			obj: struct {
				Id int `validate:"max:10|min:2"`
			}{
				Id: 1,
			},
		},
		{
			name:   "min, max in slice, min failed",
			expErr: ErrInvalidMin,
			obj: struct {
				Ids []int `validate:"min:10|max:20"`
			}{
				Ids: []int{9, 20},
			},
		},
		{
			name:   "min, max in slice, max failed",
			expErr: ErrInvalidMax,
			obj: struct {
				Ids []int `validate:"min:10|max:20"`
			}{
				Ids: []int{11, 21},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.obj)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
