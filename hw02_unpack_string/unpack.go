package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	runes := []rune(s)
	last := runes[0]
	if unicode.IsDigit(last) {
		return "", ErrInvalidString
	}
	var builder strings.Builder

	for _, r := range runes[1:] {
		switch {
		case unicode.IsDigit(r):
			if unicode.IsDigit(last) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			rep := strings.Repeat(string(last), count)
			builder.WriteString(rep)
		case !unicode.IsDigit(last):
			builder.WriteRune(last)
		}
		last = r
	}

	if !unicode.IsDigit(last) {
		builder.WriteRune(last)
	}
	return builder.String(), nil
}
