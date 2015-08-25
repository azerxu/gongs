package seq

import (
	"fmt"
	"gongs/biofile"
	"gongs/scan"
	"gongs/xopen"
	"io"
)

type Seq struct {
	Name string
	Seq  []byte
	Qual []byte
}

func (s Seq) String() string {
	if seq.Qual == nil { // seq is fasta record
		return fmt.Sprintf(">%s\n%s", s.Name, string(s.Seq))
	}
	// seq is fastq record
	return fmt.Sprintf("@%s\n%s\n+\n%s", s.Name, string(s.Seq), string(s.Qual))
}

func (s Seq) GetName() string {
	return s.Name
}

func (s Seq) GetSeq() []byte {
	return s.Seq
}

func (s Seq) GetQual() []byte {
	return nil
}

// SeqFile seq file for fasta or fastq
type SeqFile struct {
	Name string // record filename
	file io.ReadCloser
	s    *scan.Scanner
	last []byte // record last line for read name
	name string // record seq name
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

func (sf *SeqFile) Seq() *Seq {
	return &Seq{Name: sf.name, Seq: sf.seq, qual: sf.qual}
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
	sf.name = string(sf.last[1:])
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
	//	qual length  < seq length
	sf.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual length (%d) longer than seq length (%d) at line: %d",
		sf.Name, string(sf.name), len(sf.qual), len(sf.seq), sf.s.Lid()))
	return false
}

func (sf *SeqFile) Value() (string, []byte, []byte) {
	return sf.name, sf.seq, sf.qual
}

func (sf *SeqFile) Items() <-chan biofile.Seqer {
	ch := make(chan *Seq)
	go func(sf *SeqFile, ch chan *Seq) {
		for sf.Next() {
			ch <- sf.Seq()
		}
		close(ch)
	}(sf, ch)
	return ch
}
