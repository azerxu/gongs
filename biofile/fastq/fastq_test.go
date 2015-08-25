package fastq

import (
	"fmt"
	"gongs/xopen"
	"os"
	"testing"
)

var test_fq_filename = "test_fq.fastq"
var test_fq_name = "test"
var test_fq_seq = []byte("ATCG")
var test_fq_qual = []byte("@AAA")

func checkFq(fq *Fastq) bool {
	if fq.Name != test_fq_name || string(fq.Seq) != string(test_fq_seq) || string(fq.Qual) != string(test_fq_qual) {
		return false
	}
	return true
}

func create_test_fastq_file(filename string) error {
	o, err := xopen.Xcreate(filename, "w")
	if err != nil {
		return err
	}
	defer o.Close()

	for i := 0; i < 1000; i++ {
		fmt.Fprintf(o, "@%s\n%s\n+\n%s\n", test_fq_name, string(test_fq_seq), string(test_fq_qual))
	}
	return nil
}

func Test_Fastq(t *testing.T) {
	fq := &Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
	if !checkFq(fq) {
		t.Error("Test Fastq Failed!")
	}
}

func Test_Fastq_String(t *testing.T) {
	fq := Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
	if fq.String() != fmt.Sprintf("@%s\n%s\n+\n%s", test_fq_name, string(test_fq_seq), string(test_fq_qual)) {
		t.Error("Test Fastq String Failed")
	}
}

func Test_FastqFile_stdin(t *testing.T) {
	fqfile, err := Open("-")
	if err != nil {
		t.Error("Test FastqFile Stdin Error:", err)
	}
	defer fqfile.Close()
	if fqfile.Name != "STDIN" {
		t.Errorf("fqfile Name: %s", fqfile.Name)
	}
}

func Test_FastqFile_txt(t *testing.T) {
	if err := create_test_fastq_file(test_fq_filename); err != nil {
		t.Error("Test FastqFile txt create test fastq error:", err)
	}
	fqfile, err := Open(test_fq_filename)
	if err != nil {
		t.Error("Test FastqFile error:", err)
	}
	defer fqfile.Close()

	if fqfile.Name != test_fq_filename {
		t.Error("Test FastqFile Name:", fqfile.Name, "expect:", test_fq_filename)
	}

	for fq := range fqfile.Iter() {
		if !checkFq(fq) {
			t.Error("Test FastqFile Name:", []byte(fq.Name), "Qual:", string(fq.Qual), "Seq:", string(fq.Seq), "Lid:", fqfile.s.Lid())
			t.FailNow()
		}
	}
	if err := os.Remove(test_fq_filename); err != nil {
		t.Error("Test FastqFile Remove file Error:", err)
	}
}

func Test_FastqFile_gz(t *testing.T) {
	filename := test_fq_filename + ".gz"
	if err := create_test_fastq_file(filename); err != nil {
		t.Error("Test FastqFile gz create test fastq error:", err)
	}
	fqfile, err := Open(filename)
	defer fqfile.Close()
	if err != nil {
		t.Error("Test FastqFile gz error:", err)
	}
	if fqfile.Name != filename {
		t.Error("Test FastqFile gz Name:", fqfile.Name, "expect:", test_fq_filename)
	}

	for fq := range fqfile.Iter() {
		if !checkFq(fq) {
			t.Error("Test FastqFile gz fq:", fq)
		}
	}
	if err := os.Remove(filename); err != nil {
		t.Error("Test FastqFile Remove file Error:", err)
	}
}

func Test_FastqFile_nil(t *testing.T) {
	_, err := Open("tt_inposable_FYLRMane")
	if err == nil {
		t.Error("Test FastqFile nil error:", err)
	}
}

func Test_FastqFile_errformat(t *testing.T) {
	o, err := xopen.Xcreate("tt")
	if err != nil {
		t.Fail()
	}
	fmt.Fprintln(o, "@tt")
	fmt.Fprintln(o, "aaaa")
	fmt.Fprintln(o, " ")
	fmt.Fprintln(o, "+")
	fmt.Fprintln(o, "aaa")

	o.Close()
	defer os.Remove("tt")

	fqfile, err := Open("tt")
	if err != nil {
		t.Fail()
	}

	for _ = range fqfile.Iter() {
		t.Fail()
	}
	if err := fqfile.Err(); err == nil {
		t.Fail()
	}
}

func Test_FastqFile_emptyfile(t *testing.T) {
	o, err := xopen.Xcreate("tt")
	if err != nil {
		t.Error("Test FastqFile emptyfile Xcreate Error:", err)
	}
	o.Close()
	defer os.Remove("tt")

	fqfile, err := Open("tt")
	if err != nil {
		t.Error("Test FastqFile emptyfile Open Error:", err)
	}
	for fq := range fqfile.Iter() {
		t.Error("Test FastqFile emptyfile read data Error:", fq)
	}
}

func Test_Load_txt(t *testing.T) {
	create_test_fastq_file(test_fq_filename)
	defer os.Remove(test_fq_filename)
	fqch, errch := Load(test_fq_filename)
	for {
		select {
		case fq, ok := <-fqch:
			if !ok {
				fqch = nil
			} else if !checkFq(fq) {
				t.Error("Test Load txt fq error", fq)
			}
			t.Log(fq)
		case err, ok := <-errch:
			if !ok {
				errch = nil
			} else {
				t.Error("Test Load txt error:", err)
			}
		}
		if fqch == nil && errch == nil {
			break
		}
	}
}

func Test_Load_gz(t *testing.T) {
	filename := test_fq_filename + ".gz"
	create_test_fastq_file(filename)
	defer os.Remove(filename)
	fqch, errch := Load(filename)
	for {
		select {
		case fq, ok := <-fqch:
			if !ok {
				fqch = nil
			} else if !checkFq(fq) {
				t.Error("Test Load gz fq error", fq)
			}
			t.Log(fq)
		case err, ok := <-errch:
			if !ok {
				errch = nil
			} else {
				t.Error("Test Load gz error:", err)
			}
		}
		if fqch == nil && errch == nil {
			break
		}
	}
}

func Test_Load_nil(t *testing.T) {
	_, errch := Load("this is not a exit file")

	if err := <-errch; err == nil {
		t.Fail()
	}
}

func Test_Load_txts(t *testing.T) {
	create_test_fastq_file(test_fq_filename)
	defer os.Remove(test_fq_filename)
	fqch, errch := Load(test_fq_filename, test_fq_filename)
	for {
		select {
		case fq, ok := <-fqch:
			if !ok {
				fqch = nil
			} else if !checkFq(fq) {
				t.Error("Test Load txts fq error", fq)
			}
			t.Log(fq)
		case err, ok := <-errch:
			if !ok {
				errch = nil
			} else {
				t.Error("Test Load txts error:", err)
			}
		}
		if fqch == nil && errch == nil {
			break
		}
	}
}

func Test_Loads_gzs(t *testing.T) {
	filename := test_fq_filename + ".gz"
	create_test_fastq_file(filename)
	defer os.Remove(filename)
	fqch, errch := Load(filename, filename)
	for {
		select {
		case fq, ok := <-fqch:
			if !ok {
				fqch = nil
			} else if !checkFq(fq) {
				t.Error("Test Load gzs fq error", fq)
			}
			t.Log(fq)
		case err, ok := <-errch:
			if !ok {
				errch = nil
			} else {
				t.Error("Test Load gzs error:", err)
			}
		}
		if fqch == nil && errch == nil {
			break
		}
	}
}

func Test_Load_empyty_files(t *testing.T) {
	_, errch := Load()
	if err := <-errch; err != ErrEmptyInputFile {
		t.Fail()
	}
}

func Test_Loads_nil(t *testing.T) {
	_, errch := Load("no an exists file", "none file")
	if err := <-errch; err == nil {
		t.Fail()
	}
}
