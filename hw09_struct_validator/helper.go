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
	lenCommand
)

func isIn(tag string) bool {
	return strings.HasPrefix(tag, "in")
}

func isLen(tag string) bool {
	return strings.HasPrefix(tag, "len")
}

func parseCommand(tag string) Command {
	if isIn(tag) {
		return inCommand
	}

	if isLen(tag) {
		return lenCommand
	}

	return 0
}
