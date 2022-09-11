package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInForKindNotImplemented = errors.New("cmd 'in' for kind not implemented")
	ErrShouldBeIn              = errors.New("should be in")
)

func inForInt(field string, v int, in []string) error {
	for _, strExp := range in {
		exp, err := strconv.Atoi(strExp)
		if err != nil {
			return fmt.Errorf("in tag failed for value '%s': %w", strExp, err)
		}

		if exp == v {
			return nil
		}
	}

	return NewValidateError(field, fmt.Errorf("exp:'%v': %w", in, ErrShouldBeIn))
}

func inForString(field, v string, in []string) error {
	for _, exp := range in {
		if exp == v {
			return nil
		}
	}

	return NewValidateError(field, fmt.Errorf("exp:'%v': %w", in, ErrShouldBeIn))
}

func parseInTagValue(tag string) ([]string, error) {
	if len(tag) < 3 {
		return nil, ErrInvalidInCommand
	}

	v := strings.Split(tag[3:], ",")

	if len(v) == 1 && v[0] == "" {
		return nil, ErrTagValueIsEmpty
	}

	return v, nil
}

func validateIn(v reflect.Value, field, tag string) error {
	values, err := parseInTagValue(tag)
	if err != nil {
		return err
	}

	value := v.Interface()
	kind := v.Kind()

	if kind == reflect.Int {
		return inForInt(field, value.(int), values)
	}

	if kind == reflect.String {
		if err := inForString(field, value.(string), values); err != nil {
			return err
		}

		return nil
	}

	if kind == reflect.Slice {
		if intValues, ok := value.([]int); ok {
			for _, number := range intValues {
				if err := inForInt(field, number, values); err != nil {
					return err
				}
			}

			return nil
		} else if strValues, ok := value.([]string); ok {
			for _, str := range strValues {
				if err := inForString(field, str, values); err != nil {
					return err
				}
			}

			return nil
		}
	}

	return fmt.Errorf("kind '%s': %w", v.Kind().String(), ErrInForKindNotImplemented)

}
