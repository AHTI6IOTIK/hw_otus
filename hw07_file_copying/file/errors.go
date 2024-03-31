package file

import "errors"

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrStatFile              = errors.New("error getting file statistics")
	ErrNotLoadFileStat       = errors.New("not load file statistics")
)
