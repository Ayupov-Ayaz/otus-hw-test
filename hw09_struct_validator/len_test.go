package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_len(t *testing.T) {
	tests := []struct {
		obj      interface{}
		checkErr func(gotErr error) bool
	}{
		{
			checkErr: systemError(ErrInvalidLenCommand),
			obj: struct {
				Name int `validate:"len"`
			}{},
		},
		{
			checkErr: systemError(ErrInvalidLenCommandValueType),
			obj: struct {
				Name int `validate:"len:1"`
			}{},
		},
		{
			checkErr: validateError(NewValidateError("Name", ErrLenInvalid)),
			obj: struct {
				Name string `validate:"len:1"`
			}{},
		},
		{
			checkErr: validateError(NewValidateError("Name", ErrLenInvalid)),
			obj: struct {
				Name string `validate:"len:5"`
			}{
				Name: "123456",
			},
		},
		{
			obj: struct {
				Name string `validate:"len:5"`
			}{
				Name: "12345",
			},
		},
		{
			obj: struct {
				Names []string `validate:"len:3"`
			}{
				Names: []string{
					"1234",
					"12345",
				},
			},
			checkErr: validateError(NewValidateError("Names", ErrLenInvalid)),
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
