package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrParseMinMaxCommandFailed   = errors.New("parse min max command failed")
	ErrExtractMinMaxValueFailed   = errors.New("extract min max value failed")
	ErrInvalidMaxCommandValueType = errors.New("invalid validation command 'max' type, exp Int")
	ErrInvalidMinCommandValueType = errors.New("invalid validation command 'min' type, exp Int")
	ErrMinInvalid                 = errors.New("invalid min")
	ErrInvalidMax                 = errors.New("invalid max")
)

func compareInt(value, cmdValue int, cmd Command) error {
	if cmd == minCommand &&
		value < cmdValue {
		return fmt.Errorf("value=%d should be less or equal %d: %w", value, cmdValue, ErrMinInvalid)
	} else if cmd == maxCommand &&
		value > cmdValue {
		return fmt.Errorf("value=%d should be more or equal %d: %w", value, cmdValue, ErrInvalidMax)
	}

	return nil
}

func extractMinMaxTagValue(tag string) (int, error) {
	if len(tag) <= 4 {
		return 0, ErrExtractMinMaxValueFailed
	}

	v, err := strconv.Atoi(tag[4:])
	if err != nil {
		return 0, err
	}

	return v, err
}

func extractAndCompareMin(value int, field string, tag string) error {
	cmdValue, err := extractMinMaxTagValue(tag)
	if err != nil {
		return err
	}

	if err := compareInt(value, cmdValue, minCommand); err != nil {
		return NewValidateError(field, err)
	}

	return nil
}

func extractAndCompareMax(value int, field string, tag string) error {
	cmdValue, err := extractMinMaxTagValue(tag)
	if err != nil {
		return err
	}

	if err := compareInt(value, cmdValue, maxCommand); err != nil {
		return NewValidateError(field, err)
	}

	return nil
}

func validateMax(v reflect.Value, field, tag string) error {
	kind := v.Kind()

	if kind == reflect.Slice {
		values, ok := v.Interface().([]int)
		if !ok {
			return fmt.Errorf("exp []int: %w", ErrInvalidMaxCommandValueType)
		}

		for _, val := range values {
			if err := extractAndCompareMax(val, field, tag); err != nil {
				return err
			}
		}

		return nil
	}

	if kind == reflect.Int {
		value, ok := v.Interface().(int)
		if !ok {
			return fmt.Errorf("exp int: %w", ErrInvalidMaxCommandValueType)
		}

		if err := extractAndCompareMax(value, field, tag); err != nil {
			return err
		}

		return nil
	}

	return ErrInvalidMaxCommandValueType
}

func validateMin(v reflect.Value, field, tag string) error {
	kind := v.Kind()

	if kind == reflect.Slice {
		values, ok := v.Interface().([]int)
		if !ok {
			return fmt.Errorf("exp []int: %w", ErrInvalidMinCommandValueType)
		}

		for _, val := range values {
			if err := extractAndCompareMin(val, field, tag); err != nil {
				return err
			}
		}

		return nil
	}

	if kind == reflect.Int {
		value, ok := v.Interface().(int)
		if !ok {
			return fmt.Errorf("exp int: %w", ErrInvalidMinCommandValueType)
		}

		if err := extractAndCompareMin(value, field, tag); err != nil {
			return err
		}

		return nil
	}

	return ErrInvalidMinCommandValueType
}
