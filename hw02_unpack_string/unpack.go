package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(st string) (string, error) {
	var result strings.Builder

	splitStr := strings.Split(st, "")
	count := len(splitStr)
	next := 0

	for i := 0; i < count; i += next {
		currentItem := splitStr[i]
		_, err := strconv.Atoi(currentItem)
		if err == nil {
			return "", ErrInvalidString
		}

		nextItem := ""
		repeatCount := 0

		if i+1 < count {
			nextItem = splitStr[i+1]
		}

		repeatCount, err = strconv.Atoi(nextItem)

		if err != nil {
			result.WriteString(currentItem)
			next = 1
		} else {
			result.WriteString(strings.Repeat(currentItem, repeatCount))
			next = 2
		}
	}

	return result.String(), nil
}
