package seq

import (
	"fmt"
	"gongs/xopen"
	"testing"
)

var test_fq_filename = "test_fq.fastq"
var test_fq_name = "test"
var test_fq_seq = []byte("ATCG")
var test_fq_qual = []byte("@AAA")

func checkSeq(name string, seq, qual []byte) bool {
	if name != test_fq_name || string(seq) != string(test_fq_seq) || string(qual) != string(test_fq_qual) {
		return false
	}
	return true
}

func createTestFastqFile(filename string) error {
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

func TestFastqFileTxt(t *testing.T) {
	if err := createTestFastqFile(test_fq_filename); err != nil {
		t.Error("Test FastqFile txt create test fastq error:", err)
	}
	seqFile, err := Open(test_fq_filename)
	if err != nil {
		t.Error("Test FastqFile error:", err)
	}
	defer seqFile.Close()

	if seqFile.Name != test_fq_filename {
		t.Error("Test FastqFile Name:", seqFile.Name, "expect:", test_fq_filename)
	}

	for seqFile.Next() {
		name, seq, qual := seqFile.Value()
		if !checkSeq(name, seq, qual) {
			t.Error("Test FastqFile Name:", []byte(name), "Qual:", string(qual), "Seq:", string(seq), "Lid:", seqFile.s.Lid())
			t.FailNow()
		}
	}

	// if err := os.Remove(test_fq_filename); err != nil {
	// 	t.Error("Test FastqFile Remove file Error:", err)
	// }
}

// func Test_FastqFile_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	if err := create_test_fastq_file(filename); err != nil {
// 		t.Error("Test FastqFile gz create test fastq error:", err)
// 	}
// 	fqfile, err := Open(filename)
// 	defer fqfile.Close()
// 	if err != nil {
// 		t.Error("Test FastqFile gz error:", err)
// 	}
// 	if fqfile.Name != filename {
// 		t.Error("Test FastqFile gz Name:", fqfile.Name, "expect:", test_fq_filename)
// 	}

// 	for fq := range fqfile.Iter() {
// 		if !checkFq(fq) {
// 			t.Error("Test FastqFile gz fq:", fq)
// 		}
// 	}
// 	if err := os.Remove(filename); err != nil {
// 		t.Error("Test FastqFile Remove file Error:", err)
// 	}
// }

// func Test_FastqFile_nil(t *testing.T) {
// 	_, err := Open("tt_inposable_FYLRMane")
// 	if err == nil {
// 		t.Error("Test FastqFile nil error:", err)
// 	}
// }

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

// func Test_FastqFile_emptyfile(t *testing.T) {
// 	o, err := xopen.Xcreate("tt")
// 	if err != nil {
// 		t.Error("Test FastqFile emptyfile Xcreate Error:", err)
// 	}
// 	o.Close()
// 	defer os.Remove("tt")

// 	fqfile, err := Open("tt")
// 	if err != nil {
// 		t.Error("Test FastqFile emptyfile Open Error:", err)
// 	}
// 	for fq := range fqfile.Iter() {
// 		t.Error("Test FastqFile emptyfile read data Error:", fq)
// 	}
// }

// func Test_Load_txt(t *testing.T) {
// 	create_test_fastq_file(test_fq_filename)
// 	defer os.Remove(test_fq_filename)
// 	fqch, errch := Load(test_fq_filename)
// 	for {
// 		select {
// 		case fq, ok := <-fqch:
// 			if !ok {
// 				fqch = nil
// 			} else if !checkFq(fq) {
// 				t.Error("Test Load txt fq error", fq)
// 			}
// 			t.Log(fq)
// 		case err, ok := <-errch:
// 			if !ok {
// 				errch = nil
// 			} else {
// 				t.Error("Test Load txt error:", err)
// 			}
// 		}
// 		if fqch == nil && errch == nil {
// 			break
// 		}
// 	}
// }

// func Test_Load_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	create_test_fastq_file(filename)
// 	defer os.Remove(filename)
// 	fqs, err := Load(filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for fq := range fqs {
// 		if checkFq(fq) {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_Load_nil(t *testing.T) {
// 	_, err := Load("this is not a exit file")
// 	if err == nil {
// 		t.Fail()
// 	}
// }

// func Test_Loads_txt(t *testing.T) {
// 	create_test_fastq_file(test_fq_filename)
// 	defer os.Remove(test_fq_filename)
// 	fqs, err := Loads(test_fq_filename, test_fq_filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for fq := range fqs {
// 		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_Loads_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	create_test_fastq_file(filename)
// 	defer os.Remove(filename)
// 	fqs, err := Loads(filename, filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for fq := range fqs {
// 		if fq.Name != test_fq_name || fq.Seq != test_fq_seq || fq.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_Loads_empyty_files(t *testing.T) {
// 	if _, err := Loads(); err == nil {
// 		t.Fail()
// 	}
// }

// func Test_Loads_nil(t *testing.T) {
// 	if _, err := Loads("no an exists file", "none file"); err == nil {
// 		t.Fail()
// 	}
// }
