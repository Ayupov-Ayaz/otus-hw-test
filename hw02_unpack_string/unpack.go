package hw02unpackstring

import (
	"errors"
	"strconv"
)

var ErrInvalidString = errors.New("invalid string")

func makeDuplicate(expSymbol string, count int) string {
	var resp string

	for i := 0; i < count; i++ {
		resp += expSymbol
	}

	return resp
}

func castToNumber(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	return n, true
}

func isZero(s string) bool {
	n, ok := castToNumber(s)
	if !ok {
		return false
	}

	return n == 0
}

func isShielding(s string) bool {
	return s == `\`
}

func Unpack(str string) (string, error) {
	var (
		resp       string
		currSymbol string
		nextSymbol string
		lastSymbol string
	)

	runes := []rune(str)
	lastIndex := len(runes) - 1
	shielding := false

	for i := 0; i < len(runes); i++ {
		currSymbol = string(runes[i])
		number, isNumber := castToNumber(currSymbol)

		if shielding && (isNumber || isShielding(currSymbol)) {
			resp += currSymbol
			lastSymbol = currSymbol
			shielding = false
			continue
		} else if shielding {
			return "", ErrInvalidString
		}

		if isShielding(currSymbol) {
			shielding = true
			continue
		}

		if i < lastIndex {
			nextSymbol = string(runes[i+1])
		}

		if !isNumber {
			lastSymbol = currSymbol

			if !isZero(nextSymbol) { //skip
				resp += currSymbol
			}

		} else {
			if _, isNumber = castToNumber(nextSymbol); isNumber || lastSymbol == "" {
				return "", ErrInvalidString
			}

			if number == 0 {
				continue
			} else {
				resp += makeDuplicate(lastSymbol, number-1)
			}
		}
	}

	return resp, nil
}
