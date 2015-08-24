// auto read from stdin, gzip, bzip, raw or url
// auto write data to stdout, stderr, gzip or raw

package xopen

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// Xcreate write data to stdout, stderr, file or gzip file
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
