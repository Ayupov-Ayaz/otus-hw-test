package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_in(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          1,
			expectedErr: ErrShouldBeStruct,
		},
		{
			in:          "string",
			expectedErr: ErrShouldBeStruct,
		},
		{
			in:          Response{},
			expectedErr: nil,
		},
		{
			in: struct {
				Name string `validate:"in"`
			}{},
			expectedErr: ErrInvalidInCommand,
		},
		{
			in: struct {
				Name string `validate:"in:"`
			}{},
			expectedErr: ErrTagValueIsEmpty,
		},
		{
			in: Response{
				Code: 100,
			},
			expectedErr: ErrShouldBeIn,
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
				Id int `validate:"in:123,234,567,894"`
			}{
				Id: 894,
			},
		},
		{
			in: Response{
				Code: 200,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
