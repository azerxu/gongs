package lib

import (
	"fmt"
	"math"
	"os"
)

// Humanity trans number to human readable size
func Humanity(size int64) string {
	units := []string{"", "K", "M", "G", "T", "P", "E", "Z"}
	s := float64(size)
	for _, unit := range units {
		if math.Abs(s) < 1024.0 {
			return fmt.Sprintf("%.2f%sB", s, unit)
		}
		s /= 1024.0
	}
	return fmt.Sprintf("%.2f%sB", s, "Y")
}

// FileSize get file size
func FileSize(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return "Openfile error"
	}
	defer f.Close()

	if fi, err := f.Stat(); err == nil {
		return Humanity(fi.Size())
	}
	return "file stat error"
}
