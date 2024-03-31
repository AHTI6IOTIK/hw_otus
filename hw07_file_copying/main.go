package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/AHTI6IOTIK/hw_otus/hw07_file_copying/file"
)

type Conf struct {
	from, to      string
	limit, offset int64
}

var settings Conf

func init() {
	flag.StringVar(&settings.from, "from", "", "file to read from")
	flag.StringVar(&settings.to, "to", "", "file to write to")
	flag.Int64Var(&settings.limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&settings.offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	srcFile := file.NewSourceFile(
		settings.from,
		os.O_RDONLY,
		settings.offset,
	)
	srcFile.Stat()
	srcFile.CheckOffset()
	if srcFile.Err != nil {
		fmt.Println(srcFile.Err)
		return
	}

	dstFile := file.NewDestinationFile(
		settings.to,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		settings.limit,
	)
	if dstFile.Err != nil {
		err := dstFile.Err

		dstFile.Remove()
		if dstFile.Err != nil {
			err = errors.Join(err, dstFile.Err)
		}

		fmt.Println(err)

		return
	}

	err := file.Copy(srcFile, dstFile)
	if err != nil {
		fmt.Printf("copy error: %v\n", err)
	}
}
