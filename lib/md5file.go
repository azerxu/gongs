package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const Md5SumError = "Md5Sum Read bytes not equal Write bytes"

func Md5File(filename string) string {
	digest := md5.New()
	fmt.Println(os.TempDir())

	f, err := os.Open(filename)
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	buf := make([]byte, 4096)
	for {
		rn, err := f.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err.Error()
		}
		wn, err := digest.Write(buf[:rn])
		if err != nil {
			return err.Error()
		}
		if rn != wn {
			return Md5SumError
		}
	}
	return hex.EncodeToString(digest.Sum(nil))
}
