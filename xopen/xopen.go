package lib

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
)

func Xopen(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return os.Stdin, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	gfile, gerr := gzip.NewReader(file)
	if gerr == nil {
		return gfile, nil
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func Xcreate(args ...string) (io.WriteCloser, error) {
	filename := "-"
	mode := "w"
	switch l := len(args); {
	case l > 1:
		mode = args[1]
		fallthrough
	case l > 0:
		filename = args[0]
	}

	if filename == "-" {
		return os.Stdout, nil
	} else if filename == "@" {
		return os.Stderr, nil
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if mode == "a" {
		flag = os.O_APPEND | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filename, ".gz") {
		gfile := gzip.NewWriter(file)
		return gfile, nil
	}
	return file, nil
}
