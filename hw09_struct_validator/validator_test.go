package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
