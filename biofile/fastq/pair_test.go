package fastq

// func Test_Pair(t *testing.T) {
// 	fq1 := &Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
// 	fq2 := &Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
// 	p := Pair{Read1: fq1, Read2: fq2}
// 	if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 		t.Fail()
// 	}
// 	if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 		t.Fail()
// 	}
// }

// func Test_Pair_String(t *testing.T) {
// 	fq1 := &Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
// 	fq2 := &Fastq{Name: test_fq_name, Seq: test_fq_seq, Qual: test_fq_qual}
// 	p := Pair{Read1: fq1, Read2: fq2}
// 	if p.String() != fmt.Sprintf("@%s\n%s\n+\n%s\n@%s\n%s\n+\n%s", test_fq_name, test_fq_seq, test_fq_qual, test_fq_name, test_fq_seq, test_fq_qual) {
// 		t.Fail()
// 	}
// }

// func Test_FastqPairFile_txt(t *testing.T) {
// 	if err := create_test_fastq_file(test_fq_filename); err != nil {
// 		t.Fail()
// 	}
// 	pfile, err := OpenPair(test_fq_filename, test_fq_filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer pfile.Close()
// 	defer os.Remove(test_fq_filename)

// 	if pfile.Name1 != test_fq_filename || pfile.Name2 != test_fq_filename {
// 		t.Fail()
// 	}

// 	for p := range pfile.Load() {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// 	if err := pfile.Close(); err != nil {
// 		t.Fail()
// 	}
// }

// func Test_FastqPairFile_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	if err := create_test_fastq_file(filename); err != nil {
// 		t.Fail()
// 	}
// 	pfile, err := OpenPair(filename, filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer pfile.Close()
// 	defer os.Remove(filename)

// 	if pfile.Name1 != filename || pfile.Name2 != filename {
// 		t.Fail()
// 	}

// 	for p := range pfile.Load() {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_FastqPairFile_nil(t *testing.T) {
// 	_, err := OpenPair("tt_inposable_FYLRMane", "yy_no_exists_file")
// 	if err == nil {
// 		t.Fail()
// 	}
// }

// func Test_FastqPairFile_nil2(t *testing.T) {
// 	create_test_fastq_file("tt_inposable_FYLRMane")
// 	defer os.Remove("tt_inposable_FYLRMane")
// 	_, err := OpenPair("tt_inposable_FYLRMane", "yy_no_exists_file")
// 	if err == nil {
// 		t.Fail()
// 	}
// }

// func Test_FastqPairFile_emptyfile(t *testing.T) {
// 	o, err := lib.Xcreate("tt")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	o.Close()
// 	defer os.Remove("tt")

// 	pfile, err := OpenPair("tt", "tt")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for _ = range pfile.Load() {
// 		t.Fail()
// 	}
// }

// func Test_LoadPair_txt(t *testing.T) {
// 	create_test_fastq_file(test_fq_filename)
// 	defer os.Remove(test_fq_filename)
// 	ps, err := LoadPair(test_fq_filename, test_fq_filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for p := range ps {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_LoadPair_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	create_test_fastq_file(filename)
// 	defer os.Remove(filename)
// 	ps, err := LoadPair(filename, filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for p := range ps {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_LoadPair_nil(t *testing.T) {
// 	_, err := LoadPair("this is not a exit file", "file not suppose exist")
// 	if err == nil {
// 		t.Fail()
// 	}
// }

// func Test_LoadPair_nil2(t *testing.T) {
// 	create_test_fastq_file("this is not a exit file")
// 	defer os.Remove("this is not a exit file")
// 	_, err := LoadPair("this is not a exit file", "file not suppose exist")
// 	if err == nil {
// 		t.Fail()
// 	}
// }

// func Test_LoadPairs_txt(t *testing.T) {
// 	create_test_fastq_file(test_fq_filename)
// 	defer os.Remove(test_fq_filename)
// 	ps, err := LoadPairs(test_fq_filename, test_fq_filename, test_fq_filename, test_fq_filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for p := range ps {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_LoadPairs_gz(t *testing.T) {
// 	filename := test_fq_filename + ".gz"
// 	create_test_fastq_file(filename)
// 	defer os.Remove(filename)
// 	ps, err := LoadPairs(filename, filename)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	for p := range ps {
// 		if p.Read1.Name != test_fq_name || p.Read1.Seq != test_fq_seq || p.Read1.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 		if p.Read2.Name != test_fq_name || p.Read2.Seq != test_fq_seq || p.Read2.Qual != test_fq_qual {
// 			t.Fail()
// 		}
// 	}
// }

// func Test_LoadPairs_empyty_files(t *testing.T) {
// 	if _, err := LoadPairs(); err == nil {
// 		t.Fail()
// 	}
// }

// func Test_LoadPairs_nil(t *testing.T) {
// 	if _, err := LoadPairs("no an exists file", "none file", " ", ""); err == nil {
// 		t.Fail()
// 	}
// }

// func Test_LoadPairs_nopair(t *testing.T) {
// 	create_test_fastq_file("this is not a exit file")
// 	defer os.Remove("this is not a exit file")
// 	_, err := LoadPairs("this is not a exit file", "file not suppose exist", "file third")
// 	if err == nil {
// 		t.Fail()
// 	}
// }
