package main

import (
	"testing"

	"github.com/AHTI6IOTIK/hw_otus/hw07_file_copying/file"
	"github.com/AHTI6IOTIK/hw_otus/hw07_file_copying/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	type prepare func(src *mocks.ISourceFile, dst *mocks.IDestinationFile)
	testCases := []struct {
		name    string
		prepare prepare
		want    error
	}{
		{
			name: "file 0 size",
			prepare: func(src *mocks.ISourceFile, dst *mocks.IDestinationFile) {
				src.EXPECT().
					Size().
					Return(0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := mocks.NewISourceFile(t)
			dst := mocks.NewIDestinationFile(t)

			tc.prepare(src, dst)
			err := file.Copy(src, dst)
			assert.Error(t, err)
		})
	}
}
