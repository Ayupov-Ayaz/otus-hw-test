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

func inForInt(v int, in []string) error {
	for _, strExp := range in {
		exp, err := strconv.Atoi(strExp)
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
	if kind == reflect.Int {
		return inForInt(value.(int), values)
	}

	if kind == reflect.String {
		return inForString(value.(string), values)
	}

	return fmt.Errorf("kind '%s': %w", v.Kind().String(), ErrInForKindNotImplemented)
}
