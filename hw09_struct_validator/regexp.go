package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

var (
	ErrInvalidRegexpCommand   = errors.New("invalid regexp command")
	ErrRegexpInvalid          = errors.New("regexp invalid")
	ErrRegexpValueTypeInvalid = errors.New("regexp value type is invalid")
)

func parseRegexpTagValue(tag string) (string, error) {
	//regexp:
	if len(tag) <= 7 {
		return "", ErrInvalidRegexpCommand
	}

	return tag[7:], nil
}

func validateRegexp(v reflect.Value, field, tag string) error {
	rule, err := parseRegexpTagValue(tag)
	if err != nil {
		return err
	}

	_regexp, err := regexp.Compile(rule)
	if err != nil {
		return err
	}

	value := v.Interface()
	kind := v.Kind()

	if kind == reflect.String {
		if !_regexp.MatchString(value.(string)) {
			return NewValidateError(field, fmt.Errorf("rule:'%s':%w", rule, ErrRegexpInvalid))
		}

		return nil
	}

	if kind == reflect.Slice {
		strSlice, ok := value.([]string)
		if ok {
			for _, str := range strSlice {
				if !_regexp.MatchString(str) {
					return NewValidateError(field, fmt.Errorf("rule:'%s:%w", rule, ErrRegexpInvalid))
				}
			}

			return nil
		}
	}

	return ErrRegexpValueTypeInvalid
}
