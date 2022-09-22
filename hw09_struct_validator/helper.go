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
	minCommand
	maxCommand
	regexpCommand
)

func isIn(tag string) bool {
	return strings.HasPrefix(tag, "in")
}

func isLen(tag string) bool {
	return strings.HasPrefix(tag, "len")
}

func isMin(tag string) bool {
	return strings.HasPrefix(tag, "min")
}

func isMax(tag string) bool {
	return strings.HasPrefix(tag, "max")
}

func isRegexp(tag string) bool {
	return strings.HasPrefix(tag, "regexp")
}

func parseCommand(tag string) Command {
	if isIn(tag) {
		return inCommand
	}

	if isLen(tag) {
		return lenCommand
	}

	if isMin(tag) {
		return minCommand
	}

	if isMax(tag) {
		return maxCommand
	}

	if isRegexp(tag) {
		return regexpCommand
	}

	return 0
}
