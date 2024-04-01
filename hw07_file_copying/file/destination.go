//go:generate mockery --outpkg mocks  --output ../mocks --all --disable-version-string --with-expecter

package file

import (
	"errors"
	"io"
	"os"
)

type IDestinationFile interface {
	Write(p []byte) (n int, err error)
	Limit() int
	Remove() error
	Close()
}

type DestinationFile struct {
	ProcessFile
	io.WriteCloser
	limit int
}

func NewDestinationFile(
	path string,
	flg int,
	limit int64,
) *DestinationFile {
	file := NewProcessFile(path, flg)

	return &DestinationFile{
		ProcessFile: *file,
		limit:       int(limit),
	}
}

func (d *DestinationFile) Write(p []byte) (n int, err error) {
	if d.Err != nil {
		return 0, d.Err
	}

	if d.limit != 0 && len(p) > d.limit {
		return d.file.Write(p[:d.limit])
	}

	return d.file.Write(p)
}

func (d *DestinationFile) Limit() int {
	if d.Err != nil {
		return 0
	}

	return d.limit
}

func (d *DestinationFile) Remove() error {
	if d.file == nil {
		return errors.New("destination file is not load")
	}

	return os.Remove(d.file.Name())
}

func (d *DestinationFile) Close() {
	d.ProcessFile.Close()
}
