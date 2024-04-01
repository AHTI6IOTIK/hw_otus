package file

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/AHTI6IOTIK/hw_otus/hw07_file_copying/progressbar"
)

var ErrUnsupportedFile = errors.New("unsupported file")

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

	pb, err := progressbar.NewProgressBar(50, os.Stdout)
	if err != nil {
		return fmt.Errorf("initialization progress bar: %w", err)
	}

	pb.Size(limit)
	defer pb.Stop()

	_, err = io.CopyN(io.MultiWriter(pb, dst), src, int64(limit))
	if err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
