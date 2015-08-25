package scan

import (
	"fmt"
	"gongs/xopen"
	"os"
	"testing"
)

const (
	test_filename = "test_fq.fastq"
	test_name     = "@test"
	test_seq      = "ATCG"
	test_qual     = "@AAA"
)

func create_test_fastq_file(filename string) error {
	o, err := xopen.Xcreate(filename, "w")
	if err != nil {
		return err
	}
	defer o.Close()

	for i := 0; i < 10000; i++ {
		fmt.Fprintf(o, "%s\n%s\n+\n%s\n", test_name, test_seq, test_qual)
	}
	return nil
}

func TestScannerTxt(t *testing.T) {
	if err := create_test_fastq_file(test_filename); err != nil {
		t.Error(err)
		t.Fail()
	}
	defer os.Remove(test_filename)

	file, err := xopen.Xopen(test_filename)
	if err != nil {
		t.Error("TestScanner Xopen file error:", err)
		t.Fail()
	}

	s := New(file)
	var line []byte
	lid := 0
	for s.Scan() {
		line = s.Bytes()
		lid++
		if lid != s.Lid() {
			t.Error("lid error get lid:", lid, "expect:", s.Lid())
		}
		switch lid % 4 {
		case 1:
			if string(line) != test_name {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_name)
			}
		case 2:
			if string(line) != test_seq {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_seq)
			}
		case 3:
			if string(line) != "+" {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", "+")
			}
		case 0:
			if string(line) != test_qual {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_qual)
			}
		}
	}

	if err := s.Err(); err != nil {
		t.Error("TestScanner Error:", err)
	}
}

func TestScannerGz(t *testing.T) {
	if err := create_test_fastq_file(test_filename + ".gz"); err != nil {
		t.Error(err)
		t.Fail()
	}
	defer os.Remove(test_filename + ".gz")

	file, err := xopen.Xopen(test_filename + ".gz")
	if err != nil {
		t.Error("TestScanner Xopen file error:", err)
		t.Fail()
	}

	s := New(file)
	var line []byte
	lid := 0
	for s.Scan() {
		line = s.Bytes()
		lid++
		if lid != s.Lid() {
			t.Error("lid error get lid:", lid, "expect:", s.Lid())
		}
		switch lid % 4 {
		case 1:
			if string(line) != test_name {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_name)
			}
		case 2:
			if string(line) != test_seq {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_seq)
			}
		case 3:
			if string(line) != "+" {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", "+")
			}
		case 0:
			if string(line) != test_qual {
				t.Error("lineid:", lid, s.lid, "get line:", string(line), "expect:", test_qual)
			}
		}
	}

	if err := s.Err(); err != nil {
		t.Error("TestScanner Error:", err)
	}
}
