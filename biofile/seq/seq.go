package seq

import (
	"fmt"
	"gongs/scan"
	"gongs/xopen"
	"io"
)

// SeqFile seq file for fasta or fastq
type SeqFile struct {
	Name string
	file io.ReadCloser
	s    *scan.Scanner
	last []byte // record last line for read name
	name []byte
	seq  []byte
	qual []byte
	err  error
}

func Open(filename string) (*SeqFile, error) {
	file, err := xopen.Xopen(filename)
	if err != nil {
		return nil, err
	}

	return &SeqFile{
		Name: filename,
		file: file,
		s:    scan.New(file),
	}, nil
}

func (sf *SeqFile) Close() error {
	return sf.file.Close()
}

func (sf *SeqFile) Err() error {
	return sf.err
}

func (sf *SeqFile) setErr(err error) {
	if sf.err == nil || sf.err == io.EOF {
		sf.err = err
	}
}

func (sf *SeqFile) Next() bool {
	if sf.err != nil {
		return false
	}

	var line []byte
	if len(sf.last) == 0 {
		for sf.s.Scan() {
			line = sf.s.Bytes()
			if line[0] == '>' || line[0] == '@' {
				sf.last = line
				break
			}
		}
	}
	if len(sf.last) == 0 { // end of file, no record found
		sf.setErr(io.EOF)
		return false
	}
	sf.name = sf.last[1:]
	sf.last = sf.last[:0]

	// scan sequence
	sf.seq = sf.seq[:0]
	for sf.s.Scan() {
		line = sf.s.Bytes()
		if line[0] == '>' || line[0] == '+' || line[0] == '@' {
			sf.last = line
			break
		}
		sf.seq = append(sf.seq, line...)
	}
	if len(sf.last) == 0 || sf.last[0] != '+' { // fasta file
		return true
	}

	// scan fastq quality
	sf.qual = sf.qual[:0]
	sf.last = sf.last[:0]
	for sf.s.Scan() {
		sf.qual = append(sf.qual, sf.s.Bytes()...)
		if len(sf.qual) == len(sf.seq) {
			return true
		} else if len(sf.qual) > len(sf.seq) {
			sf.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual length (%d) longer than seq length (%d) at line: %d",
				sf.Name, string(sf.name), len(sf.qual), len(sf.seq), sf.s.Lid()))
			return false
		}
	}
	//	qual length  < seq length {
	sf.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual length (%d) longer than seq length (%d) at line: %d",
		sf.Name, string(sf.name), len(sf.qual), len(sf.seq), sf.s.Lid()))
	return false
}

func (sf *SeqFile) Value() (string, []byte, []byte) {
	return string(sf.name), sf.seq, sf.qual
}
