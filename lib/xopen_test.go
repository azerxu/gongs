package lib

import (
	"os"
	"testing"
)

var test_string = "this is xopen test string"
var test_filename = "test_xopen"

func Test_Xopen_txtfile_w(t *testing.T) {
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

func Test_Xopen_txtfile_a(t *testing.T) {
	filename := test_filename + ".txt"
	w, err := Xcreate(filename, "a")
	if err != nil {
		t.Log(err)
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
		t.Log(err)
		t.Fail()
	}
	for i := 0; i < n; i++ {
		if data[i] != test_string[i] {
			t.Log(err)
			t.Fail()
		}
	}
	r.Close()
	if err := os.Remove(filename); err != nil {
		t.Log(err)
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

func Test_xopen_file_stdin(t *testing.T) {
	f, err := Xopen("-")
	if err != nil {
		t.Fail()
	}
	// checks if we are getting data from stdin.
	if f.(*os.File).Fd() != 0 {
		t.Fail()
	}
}

func Test_Xcreate_file_stdout_w(t *testing.T) {
	f, err := Xcreate("-", "w")
	if err != nil {
		t.Fail()
	}
	if f.(*os.File).Fd() != 1 {
		t.Fail()
	}
}

func Test_Xcreate_file_stdout_a(t *testing.T) {
	f, err := Xcreate("-", "a")
	if err != nil {
		t.Fail()
	}
	if f.(*os.File).Fd() != 1 {
		t.Fail()
	}
}

func Test_Xcreate_file_stderr(t *testing.T) {
	f, err := Xcreate("@", "w")
	if err != nil {
		t.Fail()
	}
	if f.(*os.File).Fd() != 2 {
		t.Fail()
	}
}

func Test_Xcreate_cant_create_file(t *testing.T) {
	_, err := Xcreate("/xopen")
	t.Log(err)
	if err == nil {
		t.Fail()
	}
}

func Test_Xopen_file_nil(t *testing.T) {
	_, err := Xopen("tt")
	t.Log(err)
	if err == nil {
		t.Fail()
	}

}
