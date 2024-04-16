package arguments

import "fmt"

type InvalidArgumentsError struct {
	msg string
}

func (e InvalidArgumentsError) Error() string {
	return fmt.Sprintf("invalid argumetns: %s", e.msg)
}
