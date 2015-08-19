package lib

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrMd5Sum = errors.New("Md5Sum Read bytes not equal Write bytes")
)

func Md5File(filename string) (string, error) {
	digest := md5.New()
	fmt.Println(os.TempDir())

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, 4096)
	for {
		rn, err := f.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		wn, err := digest.Write(buf[:rn])
		if err != nil {
			return "", err
		}
		if rn != wn {
			// return "", fmt.Errorf("rn: %d, wn: %d", rn, wn)
			return "", ErrMd5Sum
		}
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
