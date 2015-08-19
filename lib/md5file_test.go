package lib

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	testString = "This is Test string.\n"
	testMd5    = "b63bd36ca6dda91fb958a8e00ca83824"
)

func Test_Md5File(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer os.Remove(f.Name())
	name := f.Name()
	f.Write([]byte(testString))
	f.Close()

	sum, err := Md5File(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if sum != testMd5 {
		t.Log("md5sum:", sum)
		t.Log("expect:", testMd5)
		t.Fail()
	}
}
