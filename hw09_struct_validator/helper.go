package hw09structvalidator

import (
	"errors"
	"strings"
)

var (
	ErrInvalidInCommand = errors.New("invalid validation command 'in'")
	ErrTagValueIsEmpty  = errors.New("validation value tag is empty")
)

type Command uint8

const (
	_ Command = iota
	inCommand
)

func isIn(tag string) bool {
	return strings.HasPrefix(tag, "in")
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

func parseCommand(tag string) Command {
	if isIn(tag) {
		return inCommand
	}

	return 0
}
