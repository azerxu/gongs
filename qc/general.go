package gc

import (
	"gongs/lib"
	"path/filepath"
)

type FileInfo struct {
	Md5  string
	Size string
	Name string
}

func Stat(filename string) *FileInfo {
	name := filepath.Base(filename)
	sizech := make(chan string)
	md5ch := make(chan string)
	go func(ch chan string, filename string) {
		ch <- lib.FileSize(filename)
	}(sizech, filename)
	go func(ch chan string, filename string) {
		ch <- lib.Md5File(filename)
	}(md5ch, filename)

	var fsize, md5sum string
	for {
		select {
		case fsize <- sizech:
			sizech = nil
		case md5sum <- md5ch:
			md5ch = nil
		}
		if sizech == nil && md5ch == nil {
			break
		}
	}

	return &FileInfo{
		Name: name,
		Size: fsize,
		Md5:  md5sum,
	}
}
