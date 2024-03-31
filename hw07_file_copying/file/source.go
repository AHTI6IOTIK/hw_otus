//go:generate mockery --outpkg mocks  --output ../mocks --all --disable-version-string --with-expecter
package file

import (
	"io"
)

type ISourceFile interface {
	CheckOffset()
	Read(p []byte) (n int, err error)
	Size() int
	Close()
}

type SourceFile struct {
	ProcessFile
	io.ReadCloser
	step   int64
	offset int64
}

func NewSourceFile(
	path string,
	flg int,
	offset int64,
) *SourceFile {
	file := NewProcessFile(path, flg)

	return &SourceFile{
		ProcessFile: *file,
		offset:      offset,
	}
}

func (s *SourceFile) CheckOffset() {
	if s.Err != nil {
		return
	}

	if s.fileInfo == nil {
		s.Err = ErrNotLoadFileStat
		return
	}

	if s.offset > s.fileInfo.Size() {
		s.Err = ErrOffsetExceedsFileSize
		return
	}
}

func (s *SourceFile) Read(p []byte) (n int, err error) {
	if s.Err != nil {
		return 0, s.Err
	}

	if s.offset > 0 {
		if s.step == 0 {
			s.step = s.offset
		}

		n, err = s.file.ReadAt(p, s.step)
		s.step += int64(n)
	} else {
		n, err = s.file.Read(p)
	}

	return n, err
}

func (s *SourceFile) Size() int {
	if s.Err != nil {
		return 0
	}

	fSize := s.fileInfo.Size()
	if s.offset > 0 {
		return int(fSize - s.offset)
	}

	return int(fSize)
}

func (s *SourceFile) Close() {
	s.ProcessFile.Close()
}
