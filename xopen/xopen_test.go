package lib

import (
	"os"
	"testing"
)

var test_string = "this is xopen test string"
var test_filename = "test_xopen"

func Test_Xopen_txtfile(t *testing.T) {
	filename := test_filename + ".txt"
	w, err := Xcreate(filename, "w")
	if err != nil {
		t.Fail()
	}

	n, err := w.Write([]byte(test_string))
	if n != len(test_string) || err != nil {
		t.Fail()
	}
	w.Close()

	r, err := Xopen(filename)
	if err != nil {
		t.Fail()
	}
	data := make([]byte, 4096)
	n, err = r.Read(data)
	if err != nil || n != len(test_string) {
		t.Fail()
	}
	for i := 0; i < n; i++ {
		if data[i] != test_string[i] {
			t.Fail()
		}
	}
	r.Close()
	if err := os.Remove(filename); err != nil {
		t.Fail()
	}
}

func Test_Xopen_gzfile(t *testing.T) {
	filename := test_filename + ".gz"
	w, err := Xcreate(filename, "w")
	if err != nil {
		t.Fail()
	}

	n, err := w.Write([]byte(test_string))
	if n != len(test_string) || err != nil {
		t.Fail()
	}
	w.Close()

	r, err := Xopen(filename)
	if err != nil {
		t.Fail()
	}
	data := make([]byte, 4096)
	n, err = r.Read(data)
	if err != nil || n != len(test_string) {
		t.Fail()
	}
	for i := 0; i < n; i++ {
		if data[i] != test_string[i] {
			t.Fail()
		}
	}
	r.Close()
	if err := os.Remove(filename); err != nil {
		t.Fail()
	}
}

func Test_Xopen_file_not_exists(t *testing.T) {
	filename := test_filename
	_, err := Xopen(filename)
	if err == nil {
		t.Fail()
	}
}
