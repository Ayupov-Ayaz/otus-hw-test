package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func systemError(expErr error) func(gotErr error) bool {
	return func(gotErr error) bool {
		return errors.Is(gotErr, expErr)
	}
}

func validateError(expErr ...ValidationError) func(gotErr error) bool {
	return func(gotErr error) bool {
		var vErrs ValidationErrors

		if errors.As(gotErr, &vErrs) {
			if len(vErrs) != len(expErr) {
				return false
			}

			for i, err := range vErrs {
				if err.Field != expErr[i].Field {
					return false
				}

				if !errors.Is(err.Err, expErr[i].Err) {
					return false
				}
			}

			return true
		}

		return false
	}
}

func TestValidateStructs(t *testing.T) {
	tests := []struct {
		checkErr func(err error) bool
		obj      interface{}
	}{
		{
			obj: User{
				Name:  "sffs",
				meta:  []byte(`{}`),
				ID:    "123",
				Age:   17,
				Email: "email.gmail.com",
				Role:  "user",
				Phones: []string{
					"12345678910",
					"1234567890",
				},
			},
			checkErr: validateError(
				NewValidateError("ID", ErrLenInvalid),
				NewValidateError("Age", ErrMinInvalid),
				NewValidateError("Email", ErrRegexpInvalid),
				NewValidateError("Role", ErrInInvalid),
				NewValidateError("Phones", ErrLenInvalid),
			),
		},
		{
			obj: User{
				ID: func() string {
					var b strings.Builder
					for i := 0; i < 36; i++ {
						b.WriteString("*")
					}
					return b.String()
				}(),
				Age:   19,
				Email: "email@gmail.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
					"19435678906",
				},
				meta: nil,
			},
			checkErr: systemError(nil),
		},
		{
			checkErr: validateError(NewValidateError("Version", ErrLenInvalid)),
			obj:      App{Version: "1234"},
		},
		{
			obj: App{Version: "12345"},
		},
		{
			obj: Token{},
		},
		{
			checkErr: validateError(NewValidateError("Code", ErrInInvalid)),
			obj: Response{
				Code: 300,
			},
		},
		{
			checkErr: systemError(nil),
			obj:      Response{Code: 200},
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
