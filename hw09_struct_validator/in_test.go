package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_in(t *testing.T) {
	tests := []struct {
		in       interface{}
		checkErr func(gotErr error) bool
	}{
		{
			in:       1,
			checkErr: systemError(ErrShouldBeStruct),
		},
		{
			in:       "string",
			checkErr: systemError(ErrShouldBeStruct),
		},
		{
			in: struct {
				Name string `validate:"in"`
			}{},
			checkErr: systemError(ErrInvalidInCommand),
		},
		{
			in: struct {
				Name string `validate:"in:"`
			}{},
			checkErr: systemError(ErrTagValueIsEmpty),
		},
		{
			in: Response{
				Code: 100,
			},
			checkErr: validateError(NewValidateError("Code", ErrInInvalid)),
		},
		{
			in: struct {
				Name string `validate:"in:tommy,Tommy"`
			}{
				Name: "tommy",
			},
		},
		{
			in: struct {
				Name string `validate:"in:tommy,Tommy"`
			}{
				Name: "Tommy",
			},
		},
		{
			in: struct {
				ID int `validate:"in:123,234,567,894"`
			}{
				ID: 894,
			},
		},
		{
			in: Response{
				Code: 200,
			},
		},
		{
			in: struct {
				Codes []int `validate:"in:100,500,400"`
			}{
				Codes: []int{100, 200, 300, 400, 500},
			},
			checkErr: validateError(
				NewValidateError("Codes", ErrInInvalid),
			),
		},
		{
			in: struct {
				Codes []string `validate:"in:a,b,c"`
			}{
				Codes: []string{"a", "b", "c", "d"},
			},
			checkErr: validateError(
				NewValidateError("Codes", ErrInInvalid),
			),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			if tt.checkErr != nil {
				require.True(t, tt.checkErr(err))
			} else {
				require.Nil(t, err)
			}
		})
	}
}
