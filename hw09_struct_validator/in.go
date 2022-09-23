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
	ErrInInvalid               = errors.New("should be in")
)

func validateInForInt(field string, v int, in []string) error {
	for _, strExp := range in {
		exp, err := strconv.Atoi(strExp)
		if err != nil {
			return fmt.Errorf("in tag failed for value '%s': %w", strExp, err)
		}

		if exp == v {
			return nil
		}
	}

	return NewValidateError(field, fmt.Errorf("exp:'%v': %w", in, ErrInInvalid))
}

func validateInForSlice(value interface{}, field string, values []string) error {
	if intValues, ok := value.([]int); ok {
		for _, number := range intValues {
			if err := validateInForInt(field, number, values); err != nil {
				return err
			}
		}
	} else if strValues, ok := value.([]string); ok {
		for _, str := range strValues {
			if err := validateInForString(field, str, values); err != nil {
				return err
			}
		}
	}

	return ErrInForKindNotImplemented
}

func validateInForString(field, v string, in []string) error {
	for _, exp := range in {
		if exp == v {
			return nil
		}
	}

	return NewValidateError(field, fmt.Errorf("exp:'%v': %w", in, ErrInInvalid))
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

	switch v.Kind() {
	case reflect.Int:
		return validateInForInt(field, value.(int), values)
	case reflect.String:
		return validateInForString(field, castToString(value), values)
	case reflect.Slice:
		return validateInForSlice(value, field, values)
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128,
		reflect.Complex64, reflect.Float32, reflect.Float64, reflect.Func,
		reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Map, reflect.Ptr,
		reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		return fmt.Errorf("typeKind '%s': %w", v.Kind().String(),
			ErrInForKindNotImplemented) // это не я дурак, это линтер заставил написать так
	}

	return fmt.Errorf("typeKind '%s': %w", v.Kind().String(),
		ErrInForKindNotImplemented)
}
