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
var test_fq_qual = "@AAA"

func create_test_fastq_file(filename string) error {
	o, err := lib.Xcreate(filename, "w")
	if err != nil {
		return err
	}
	fmt.Fprintln(o, "# this is a test fastq file")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(o, "@%s\n", test_fq_name)
		fmt.Fprintf(o, "%s\n", test_fq_seq)
		fmt.Fprintln(o, "+")
		fmt.Fprintf(o, "%s\n", test_fq_qual)
	}
	fmt.Fprintln(o, "  ")
	fmt.Fprintf(o, "@%s\n", test_fq_name)
	fmt.Fprintf(o, "%s\n%s\n", test_fq_seq[:2], test_fq_seq[2:])
	fmt.Fprintln(o, "+")
	fmt.Fprintf(o, "%s\n%s\n", test_fq_qual[:1], test_fq_qual[1:])
	o.Close()
	return nil
}

func Test_Fastq(t *testing.T) {
	fq := Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
	if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
		t.Fail()
	}
}

func Test_Fastq_String(t *testing.T) {
	fq := Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
	if fq.String() != fmt.Sprintf("@%s\n%s\n+\n%s", test_fq_name, test_fq_seq, test_fq_qual) {
		t.Fail()
	}
}

func Test_FastqFile_stdin(t *testing.T) {
	fqfile, err := Open("-")
	if err != nil {
		t.Fail()
	}
	defer fqfile.Close()
	if fqfile.Name != "STDIN" {
		t.Fail()
	}
}

func Test_FastqFile_txt(t *testing.T) {
	if err := create_test_fastq_file(test_fq_filename); err != nil {
		t.Fail()
	}
	fqfile, err := Open(test_fq_filename)
	if err != nil {
		t.Fail()
	}
	defer fqfile.Close()

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

func Test_FastqFile_nil(t *testing.T) {
	_, err := Open("tt_inposable_FYLRMane")
	if err == nil {
		t.Fail()
	}
}

// func Test_FastqFile_errformat(t *testing.T) {
// 	o, err := lib.Xcreate("tt")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	fmt.Fprintln(o, "@tt")
// 	fmt.Fprintln(o, "aaaa")
// 	fmt.Fprintln(o, " ")
// 	fmt.Fprintln(o, "+")
// 	fmt.Fprintln(o, "aaa")

// 	o.Close()
// 	defer os.Remove("tt")

// 	fqfile, err := Open("tt")
// 	if err != nil {
// 		t.Fail()
// 	}

// 	for _ = range fqfile.Load() {
// 		t.Fail()
// 	}
// }

func Test_FastqFile_emptyfile(t *testing.T) {
	o, err := lib.Xcreate("tt")
	if err != nil {
		t.Fail()
	}
	o.Close()
	defer os.Remove("tt")

	fqfile, err := Open("tt")
	if err != nil {
		t.Fail()
	}
	for _ = range fqfile.Load() {
		t.Fail()
	}
}

func Test_Load_txt(t *testing.T) {
	create_test_fastq_file(test_fq_filename)
	defer os.Remove(test_fq_filename)
	fqs, err := Load(test_fq_filename)
	if err != nil {
		t.Fail()
	}
	for fq := range fqs {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
}

func Test_Load_gz(t *testing.T) {
	filename := test_fq_filename + ".gz"
	create_test_fastq_file(filename)
	defer os.Remove(filename)
	fqs, err := Load(filename)
	if err != nil {
		t.Fail()
	}
	for fq := range fqs {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
}

func Test_Load_nil(t *testing.T) {
	_, err := Load("this is not a exit file")
	if err == nil {
		t.Fail()
	}
}

func Test_Loads_txt(t *testing.T) {
	create_test_fastq_file(test_fq_filename)
	defer os.Remove(test_fq_filename)
	fqs, err := Loads(test_fq_filename, test_fq_filename)
	if err != nil {
		t.Fail()
	}
	for fq := range fqs {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
}

func Test_Loads_gz(t *testing.T) {
	filename := test_fq_filename + ".gz"
	create_test_fastq_file(filename)
	defer os.Remove(filename)
	fqs, err := Loads(filename, filename)
	if err != nil {
		t.Fail()
	}
	for fq := range fqs {
		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
			t.Fail()
		}
	}
}

func Test_Loads_empyty_files(t *testing.T) {
	if _, err := Loads(); err == nil {
		t.Fail()
	}
}

func Test_Loads_nil(t *testing.T) {
	if _, err := Loads("no an exists file", "none file"); err == nil {
		t.Fail()
	}
}
