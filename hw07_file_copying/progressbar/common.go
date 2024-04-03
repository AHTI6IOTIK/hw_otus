package progressbar

import (
	"fmt"
	"io"
	"strings"
)

type ProgressBar struct {
	width int
	out   io.Writer
	buf   *strings.Builder
	size  int
	step  int
}

func NewProgressBar(
	width int,
	out io.Writer,
) (*ProgressBar, error) {
	if width <= 0 {
		return nil, ErrInvalidPbWidth
	}

	if out == nil {
		return nil, ErrInvalidPbOut
	}

	buf := new(strings.Builder)
	buf.Grow(width)

	return &ProgressBar{
		width: width,
		out:   out,
		buf:   buf,
	}, nil
}

func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.step += n

	fmt.Fprint(pb.out, "\r"+pb.progress())

	return n, nil
}

func (pb *ProgressBar) Size(size int) {
	pb.size = size
}

func (pb *ProgressBar) Stop() {
	fmt.Println()
}

func (pb *ProgressBar) progress() string {
	percent := pb.step * 100 / pb.size
	current := pb.width * percent / 100

	pb.buf.WriteString(strings.Repeat("#", current))
	pb.buf.WriteString(strings.Repeat(" ", pb.width-current))
	defer pb.buf.Reset()

	return fmt.Sprintf("[%s] %d%%", pb.buf.String(), percent)
}
