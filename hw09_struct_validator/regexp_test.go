package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRegexp(t *testing.T) {
	tests := []struct {
		checkErr func(err error) bool
		obj      interface{}
	}{
		{
			checkErr: systemError(ErrInvalidRegexpCommand),
			obj: struct {
				Name string `validate:"regexp"`
			}{},
		},
		{
			checkErr: validateError(NewValidateError("Name", ErrRegexpInvalid)),
			obj: struct {
				Name string `validate:"regexp:\\d+"`
			}{
				Name: "dsg",
			},
		},
		{
			obj: struct {
				ID string `validate:"regexp:\\d+"`
			}{
				ID: "\\1",
			},
		},
		{
			checkErr: validateError(NewValidateError("Ids", ErrLenInvalid)),
			obj: struct {
				Ids []string `validate:"regexp:\\d+|len:3"`
			}{
				Ids: []string{"123", "234", "2345"},
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
