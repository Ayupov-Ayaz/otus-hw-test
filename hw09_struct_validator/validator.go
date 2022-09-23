package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

var ErrShouldBeStruct = errors.New("validation type should be struct")

func validate(vErrors *ValidationErrors, tags []string, val reflect.Value, field string) error {
	var err error

	for _, tag := range tags {
		switch parseCommand(tag) {
		case inCommand:
			err = validateIn(val, field, tag)
		case lenCommand:
			err = validateLen(val, field, tag)
		case maxCommand:
			err = validateMax(val, field, tag)
		case minCommand:
			err = validateMin(val, field, tag)
		case regexpCommand:
			err = validateRegexp(val, field, tag)
		}

		if err != nil && !vErrors.add(err) {
			return err
		}
	}

	return nil
}

func Validate(v interface{}) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return ErrShouldBeStruct
	}

	vErrors := &ValidationErrors{}
	for i := 0; i < valueOf.NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)
		_tag := field.Tag.Get("validate")
		if _tag == "" {
			continue
		}

		val := valueOf.Field(i)
		tags := strings.Split(_tag, "|")

		if err := validate(vErrors, tags, val, field.Name); err != nil {
			return err
		}
	}

	return vErrors.getErrors()
}

type ValidationError struct {
	Field string
	Err   error
}

func NewValidateError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	buff := strings.Builder{}
	buff.WriteString("field '")
	buff.WriteString(v.Field)
	buff.WriteString("' invalid: ")
	buff.WriteString(v.Err.Error())

	return buff.String()
}

func (v ValidationErrors) Error() string {
	buff := strings.Builder{}

	last := len(v) - 1

	for i, err := range v {
		buff.WriteString(err.Error())

		if i != last {
			buff.WriteString("\n")
		}
	}

	return buff.String()
}

func (v *ValidationErrors) add(err error) bool {
	var vErr ValidationError
	if errors.As(err, &vErr) {
		*v = append(*v, vErr)
		return true
	}

	return false
}

func (v ValidationErrors) getErrors() error {
	if len(v) == 0 {
		return nil
	}

	return v
}
