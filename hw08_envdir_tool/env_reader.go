package main

import (
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

const (
	lfByte  = 0x0A
	nulByte = 0x00
)

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("stat dir: %s | %w", dir, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s: is not directory", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filePath, value string
	var fileData []byte
	var isRemove bool
	result := make(Environment)

	for _, entry := range entries {
		if strings.Contains(entry.Name(), "=") {
			continue
		}

		filePath = fmt.Sprintf(
			"%s/%s",
			dir,
			entry.Name(),
		)

		fileData, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read envfile: %s | %w", entry.Name(), err)
		}

		lfIndex := findIndex(fileData, lfByte)
		if lfIndex > 0 {
			fileData = fileData[:lfIndex]
		}

		nulIndex := findIndex(fileData, nulByte)
		if nulIndex > 0 {
			fileData[nulIndex] = lfByte
		}

		value = string(fileData)
		if len(fileData) == 0 {
			isRemove = true
		}

		value = strings.TrimRight(strings.TrimRight(value, "\t"), " ")
		result[entry.Name()] = EnvValue{
			Value:      value,
			NeedRemove: isRemove,
		}

		isRemove = false
	}

	return result, nil
}

func findIndex(data []byte, b int) int {
	for i, item := range data {
		if item == byte(b) {
			return i
		}
	}
	return 0
}
