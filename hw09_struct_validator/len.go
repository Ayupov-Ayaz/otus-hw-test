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
	ErrInvalidLen                 = errors.New("invalid len")
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

func validateLen(v reflect.Value, tag string) error {
	value := v.Interface()
	if v.Kind() != reflect.String {
		return ErrInvalidLenCommandValueType
	}

	exp, err := parseLenTagValue(tag)
	if err != nil {
		return err
	}

	got := utf8.RuneCountInString(value.(string))

	if got != exp {
		return fmt.Errorf("exp '%d', got '%d': %w", exp, got, ErrInvalidLen)
	}

	return nil
}
