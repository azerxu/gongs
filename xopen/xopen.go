// auto read from stdin, gzip, bzip, raw or url
// auto write data to stdout, stderr, gzip or raw

package xopen

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// BzipReadCloser add a Close function for bzip2
type BzipReadCloser struct {
	r    io.Reader
	file *os.File
}

// Read Using bzip2 Reader's Read method
func (bz *BzipReadCloser) Read(buf []byte) (int, error) {
	return bz.r.Read(buf)
}

// Close Close BzipReadCloser
func (bz *BzipReadCloser) Close() error {
	return bz.file.Close()
}

// Xopen read from stdin, raw, gzip, bzip2 or url
func Xopen(filename string) (io.ReadCloser, error) {
	if filename == "-" { // check input is stdin or not
		return os.Stdin, nil
	}

	// check input from an url string or not
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		r, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		if r.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http error while loading %v. status: %v", filename, r.Status)
		}
		return r.Body, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(filename, ".bz2") || strings.HasSuffix(filename, ".bz") {
		return &BzipReadCloser{r: bzip2.NewReader(file), file: file}, nil
	}

	if gfile, err := gzip.NewReader(file); err == nil {
		return gfile, nil
	}

	if _, err = file.Seek(0, 0); err != nil {
		return nil, err
	}
	return file, nil
}

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
