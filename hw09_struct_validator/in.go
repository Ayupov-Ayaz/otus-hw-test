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

func parseInt64(kind reflect.Kind, value interface{}) int64 {
	switch kind {
	case reflect.Int:
		return int64(value.(int))
	case reflect.Int8:
		return int64(value.(int8))
	case reflect.Int16:
		return int64(value.(int16))
	case reflect.Int32:
		return int64(value.(int32))
	default:
		return value.(int64)
	}
}

func parseUint64(kind reflect.Kind, value interface{}) uint64 {
	switch kind {
	case reflect.Uint:
		return uint64(value.(uint))
	case reflect.Uint8:
		return uint64(value.(uint8))
	case reflect.Uint16:
		return uint64(value.(uint16))
	case reflect.Uint32:
		return uint64(value.(uint32))
	default:
		return value.(uint64)
	}
}

func inForInt(v int64, in []string) error {
	for _, strExp := range in {
		exp, err := strconv.ParseInt(strExp, 10, 64)
		if err != nil {
			return fmt.Errorf("in tag failed for value '%s': %w", strExp, err)
		}

		if exp == v {
			return nil
		}
	}

	return fmt.Errorf("exp:'%v': %w", in, ErrShouldBeIn)
}

func inForUint(v uint64, in []string) error {
	for _, strExp := range in {
		exp, err := strconv.ParseUint(strExp, 10, 64)
		if err != nil {
			return fmt.Errorf("in tag failed for value '%s': %w", strExp, err)
		}

		if exp == v {
			return nil
		}
	}

	return fmt.Errorf("exp:'%v': %w", in, ErrShouldBeIn)
}

func inForString(v string, in []string) error {
	for _, exp := range in {
		if exp == v {
			return nil
		}
	}

	return fmt.Errorf("exp:'%v': %w", in, ErrShouldBeIn)
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

func validateIn(v reflect.Value, tag string) error {
	values, err := parseInTagValue(tag)
	if err != nil {
		return err
	}

	value := v.Interface()

	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		return inForInt(parseInt64(kind, value), values)
	}

	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		return inForUint(parseUint64(kind, value), values)
	}

	if kind == reflect.String {
		return inForString(value.(string), values)
	}

	return fmt.Errorf("kind '%s': %w", v.Kind().String(), ErrInForKindNotImplemented)
}
