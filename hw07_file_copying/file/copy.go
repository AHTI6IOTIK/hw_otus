package file

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrUnsupportedFile = errors.New("unsupported file")
)

func Copy(
	src ISourceFile,
	dst IDestinationFile,
) error {
	fileSize := src.Size()
	if fileSize <= 0 {
		return ErrUnsupportedFile
	}

	limit := fileSize
	if l := dst.Limit(); l > 0 {
		limit = l
	}

	if limit > fileSize {
		limit = fileSize
	}

	var err error
	defer func(err error) {
		if err != nil {
			fmt.Printf("\r")
		} else {
			fmt.Println()
		}
	}(err)

	p := make([]byte, limit)
	offset := 0
	for offset < limit {
		n, err := src.Read(p)
		offset += n

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	prg := progress(limit)
	for i := 0; i < limit; i++ {
		_, err := dst.Write(p[i : i+1])
		if err != nil {
			break
		}

		fmt.Print(prg(i + 1))
	}

	return nil
}

func progress(total int) func(step int) string {
	const width = int(50)

	return func(step int) string {
		percent := step * 100 / total
		current := width * percent / 100
		var s strings.Builder

		s.WriteString(strings.Repeat("#", current))
		s.WriteString(strings.Repeat(" ", width-current))

		return fmt.Sprintf("\r[%s] %d%v", s.String(), percent, "%")
	}
}
