package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(st string) (string, error) {
	var result strings.Builder

	sta := strings.Split(st, "")
	count := len(sta)
	next := 0

	for i := 0; i < count; i += next {
		currentItem := sta[i]
		_, err := strconv.Atoi(currentItem)
		if err == nil {
			return "", ErrInvalidString
		}

		nextItem := ""
		size := 0

		if i+1 < count {
			nextItem = sta[i+1]
		}

		size, err = strconv.Atoi(nextItem)

		if err != nil {
			result.WriteString(currentItem)
			next = 1
		} else {
			result.WriteString(strings.Repeat(currentItem, size))
			next = 2
		}
	}

	return result.String(), nil
}
