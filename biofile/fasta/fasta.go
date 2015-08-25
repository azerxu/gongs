// Read and Parse fasta format file

package fasta

import (
	"fmt"
	"gongs/biofile"
	"gongs/scan"
	"gongs/xopen"
	"io"
	"strings"
)

const (
	LineWidth    = 60
	MaxLineWidth = 200
)

type Fasta struct {
	Name string
	Seq  []byte
}

func (fa Fasta) GetName() string {
	return fa.Name
}

func (fa Fasta) GetSeq() []byte {
	return fa.Seq
}

func (fa Fasta) GetQual() []byte {
	return nil
}

func (fa Fasta) String() string {
	if len(fa.Seq) > MaxLineWidth {
		buf := []string{}
		for start, end := 0, len(fa.Seq); start < end; start += LineWidth {
			if start+LineWidth < l {
				buf = append(buf, fa.Seq[start:start+LineWidth])
			} else {
				buf = append(buf, fa.Seq[start:end])
			}
		}
		return fmt.Sprintf(">%s\n%s", fa.Name, strings.Join(buf, "\n"))
	}
	return fmt.Sprintf(">%s\n%s", fa.Name, fa.Seq)
}

func (fa Fasta) Id() string {
	if n := strings.IndexByte(fa.Name, ' '); n >= 0 {
		return fa.Name[:n]
	}
	return fa.Name
}

func (fa *Fasta) Slice(start, end int) *Fasta {
	return &Fasta{Name: fa.Name, Seq: fa.Seq[start:end]}
}

type FastaFile struct {
	Name string
	file io.ReadCloser
	s    *scan.Scanner
	err  error
	name string
	seq  []byte
	last []byte
}

func Open(filename string) (*FastaFile, error) {
	file, err := xopen.Xopen(filename)
	if err != nil {
		return nil, err
	}

	return &FastaFile{
		Name: filename,
		file: file,
		s:    scan.New(file),
	}, nil
}

func (ff *FastaFile) Err() error {
	if ff.err == nil || ff.err == io.EOF {
		return ff.s.Err()
	}
	return ff.err
}

func (ff *FastaFile) Close() error {
	return ff.file.Close()
}

func (ff *FastaFile) Next() bool {
	if ff.err != nil {
		return false
	}
	var line []byte
	if len(ff.last) == 0 {
		for ff.s.Scan() { // get fasta record name
			if line = ff.s.Bytes(); (len(line) > 0) && (line[0] == '>') {
				ff.last = line
				break
			}
		}
	}
	if len(ff.last) == 0 { // end of file
		ff.setErr(io.EOF)
		return false
	}

	ff.name = string(ff.last[1:])
	ff.last = ff.last[:0]
	ff.seq = ff.seq[:0]
	for ff.s.Scan() { // get fasta record sequence
		line = ff.s.Bytes()
		if len(line) > 0 && line[0] == '>' {
			ff.last = line
			return true
		}
		ff.seq = append(ff.seq, line...)
	}
	ff.setErr(io.EOF)
	return true
}

func (ff *FastaFile) Fa() *Fasta {
	return &Fasta{Name: ff.name, Seq: ff.seq}
}

func (ff *FastaFile) Value() (string, []byte, []byte) {
	return ff.name, ff.seq, nil
}

func (ff *FastaFile) Iter() <-chan *Fasta {
	ch := make(chan *Fasta)
	go func(ch chan *Fasta, ff *FastaFile) {
		for ff.Next() {
			ch <- ff.Fa()
		}
		close(ch)
	}(ch, ff)
	return ch
}

func (ff *FastaFile) Seqs() <-chan biofile.Seqer {
	ch := make(chan biofile.Seqer)
	go func(ch chan biofile.Seqer, ff *FastaFile) {
		for ff.Next() {
			ch <- ff.Fa()
		}
		close(ch)
	}(ch, ff)
	return ch
}
