package file

import (
	"fmt"
	"io"
	"os"
)

type ProcessFile struct {
	io.Closer
	file     *os.File
	fileInfo os.FileInfo
	Err      error
}

func NewProcessFile(
	path string,
	flg int,
) *ProcessFile {
	file, err := os.OpenFile(path, flg, 0644)

	return &ProcessFile{
		file: file,
		Err:  err,
	}
}

func (p *ProcessFile) Stat() {
	if p.Err != nil {
		return
	}
	fileInfo, err := p.file.Stat()
	if err != nil {
		p.Err = fmt.Errorf("%v: %v", ErrStatFile, err)
		return
	}

	p.fileInfo = fileInfo
}

func (p *ProcessFile) Close() {
	err := p.file.Close()

	if p.Err != nil {
		p.Err = fmt.Errorf("%w: %v", p.Err, err)
	} else {
		p.Err = err
	}
}
