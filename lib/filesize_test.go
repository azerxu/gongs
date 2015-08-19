package lib

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func TestHumanity(t *testing.T) {
	sizes := []int64{1, 1024, 4100, 2048000, 1234567890, 45678901234567890}
	strs := []string{"1.00B", "1.00KB", "4.00KB", "1.95MB", "1.15GB", "40.57PB"}
	for i, size := range sizes {
		if s := Humanity(size); s != strs[i] {
			t.Log("size:", size, "Humanity:", s, "expect string:", strs[i])
			t.Fail()
		}
	}
}

func TestFileSize(t *testing.T) {
	size := 1234567
	tempfile, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer os.Remove(tempfile.Name())

	f := bufio.NewWriter(tempfile)
	for i := 0; i < size; i++ {
		f.Write([]byte("ATCG"))
	}
	f.Flush()
	tempfile.Close()

	fsize := FileSize(tempfile.Name())
	if s := Humanity(int64(size * 4)); s != fsize {
		t.Log("Humanity size:", s, "filesize:", fsize)
		t.Fail()
	}
}
