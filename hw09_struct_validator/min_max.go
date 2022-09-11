package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrInvalidMinMaxCommand          = errors.New("invalid validation command 'min' or 'max'")
	ErrInvalidMinMaxCommandValueType = errors.New("invalid validation command 'min' or 'max' type, exp Int")
	ErrInvalidMin                    = errors.New("invalid min")
	ErrInvalidMax                    = errors.New("invalid max")
)

func compareInt(value, cmdValue int, cmd Command) error {
	if cmd == minCommand &&
		value < cmdValue {
		return fmt.Errorf("value=%d should be less or equal %d: %w", value, cmdValue, ErrInvalidMin)
	} else if cmd == maxCommand &&
		value > cmdValue {
		return fmt.Errorf("value=%d should be more or equal %d: %w", value, cmdValue, ErrInvalidMax)
	}

	return nil
}

func extractMinMaxCommand(tag string) (int, Command, error) {
	if len(tag) <= 4 {
		return 0, 0, ErrInvalidMinMaxCommand
	}

	v, err := strconv.Atoi(tag[4:])
	if err != nil {
		return 0, 0, err
	}

	var cmd Command

	if isMin(tag) {
		cmd = minCommand
	} else if isMax(tag) {
		cmd = maxCommand
	} else {
		err = ErrInvalidMinMaxCommand
	}

	return v, cmd, err
}

func extractAndCompareMinMax(value int, field string, tag string) error {
	cmdValue, cmd, err := extractMinMaxCommand(tag)
	if err != nil {
		return err
	}

	if err := compareInt(value, cmdValue, cmd); err != nil {
		return NewValidateError(field, err)
	}

	return nil
}

func validateMinMax(v reflect.Value, field, tag string) error {
	kind := v.Kind()

	if kind != reflect.Slice && kind != reflect.Int {
		return ErrInvalidMinMaxCommandValueType
	}

	if kind == reflect.Slice {
		values, ok := v.Interface().([]int)
		if !ok {
			return fmt.Errorf("exp []int: %w", ErrInvalidMinMaxCommandValueType)
		}

		for _, val := range values {
			if err := extractAndCompareMinMax(val, field, tag); err != nil {
				return err
			}
		}

	} else {
		value, ok := v.Interface().(int)
		if !ok {
			return fmt.Errorf("exp int: %w", ErrInvalidMinMaxCommandValueType)
		}

		if err := extractAndCompareMinMax(value, field, tag); err != nil {
			return err
		}
	}

	return nil
}
