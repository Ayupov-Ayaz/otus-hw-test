package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

var (
	ErrInvalidLenCommand          = errors.New("invalid validation command 'len'")
	ErrInvalidLenCommandValueType = errors.New("invalid validation command 'len' value type")
	ErrLenInvalid                 = errors.New("invalid len")
)

func parseLenTagValue(tag string) (int, error) {
	// len:
	if len(tag) <= 4 {
		return 0, ErrInvalidLenCommand
	}

	value := tag[4:]

	exp, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return exp, nil
}

func checkLen(field, str string, exp int) error {
	got := utf8.RuneCountInString(str)

	if got != exp {
		return NewValidateError(field, fmt.Errorf("exp '%d', got '%d': %w", exp, got, ErrLenInvalid))
	}

	return nil
}

func checkLenSliceString(field string, v interface{}, exp int) error {
	strList, ok := v.([]string)
	if !ok {
		return fmt.Errorf("should be slice string:%w", ErrInvalidLenCommandValueType)
	}

	for _, str := range strList {
		if err := checkLen(field, str, exp); err != nil {
			return err
		}
	}

	return nil
}

func validateLen(v reflect.Value, field, tag string) error {
	exp, err := parseLenTagValue(tag)
	if err != nil {
		return err
	}

	value := v.Interface()

	if v.Kind() == reflect.String {
		return checkLen(field, castToString(value), exp)
	} else if v.Kind() == reflect.Slice {
		return checkLenSliceString(field, value, exp)
	}

	return ErrInvalidLenCommandValueType
}
