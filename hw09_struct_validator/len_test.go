package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_len(t *testing.T) {
	tests := []struct {
		obj         interface{}
		expectedErr error
	}{
		{
			expectedErr: ErrInvalidLenCommandValueType,
			obj: struct {
				Name int `validate:"len"`
			}{},
		},
		{
			expectedErr: ErrInvalidLenCommand,
			obj: struct {
				Name string `validate:"len"`
			}{},
		},
		{
			expectedErr: ErrInvalidLen,
			obj: struct {
				Name string `validate:"len:1"`
			}{},
		},
		{
			expectedErr: ErrInvalidLen,
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
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.obj)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
