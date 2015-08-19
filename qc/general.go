package gc

import "path"
import "gongs/lib"

type FileInfo struct {
	Md5  string
	Size string
	Name string
}

func Stat(filename string) *FileInfo {
	name := path.Base(filename)
	sizech := make(chan string, 1)
	md5ch := make(chan string, 1)
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
