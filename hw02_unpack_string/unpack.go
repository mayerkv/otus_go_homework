package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	chars := []rune(s)

	for i, char := range chars {
		var next rune
		if i < len(chars)-1 {
			next = chars[i+1]
		}

		if unicode.IsDigit(char) && i == 0 {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(next) && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(char) {
			continue
		}

		if unicode.IsDigit(next) {
			cnt, _ := strconv.Atoi(string(next))
			for j := 0; j < cnt; j++ {
				b.WriteRune(char)
			}
			continue
		}

		if !unicode.IsDigit(next) {
			b.WriteRune(char)
			continue
		}
	}

	return b.String(), nil
}
