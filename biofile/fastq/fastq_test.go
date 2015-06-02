package fastq

import (
	"fmt"
	"gongs/lib"
	"os"
	"testing"
)

var test_fq_filename = "test_fq.fastq"
var test_fq_name = "test"
var test_fq_seq = "ATCG"
var test_fq_qual = "AAAA"

func create_test_fastq_file(filename string) error {
	o, err := lib.Xcreate(filename, "w")
	if err != nil {
		return err
	}
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(o, "@%s\n", test_fq_name)
		fmt.Fprintf(o, "%s\n", test_fq_seq)
		fmt.Fprintln(o, "+")
		fmt.Fprintf(o, "%s\n", test_fq_qual)
	}
	o.Close()
	return nil
}

func Test_Fastq(t *testing.T) {
	fq := Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
	if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
		t.Fail()
	}
}

func Test_FastqFile_txt(t *testing.T) {
	if err := create_test_fastq_file(test_fq_filename); err != nil {
		t.Fail()
	}
	fqfile, err := Open(test_fq_filename)
	defer fqfile.Close()
	if err != nil {
		t.Fail()
	}
	if fqfile.Name != test_fq_filename {
		t.Fail()
	}

	for fq := range fqfile.Load() {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
	if err := os.Remove(test_fq_filename); err != nil {
		t.Fail()
	}
}

func Test_FastqFile_gz(t *testing.T) {
	filename := test_fq_filename + ".gz"
	if err := create_test_fastq_file(filename); err != nil {
		t.Fail()
	}
	fqfile, err := Open(filename)
	defer fqfile.Close()
	if err != nil {
		t.Fail()
	}
	if fqfile.Name != filename {
		t.Fail()
	}

	for fq := range fqfile.Load() {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
	if err := os.Remove(filename); err != nil {
		t.Fail()
	}
}
