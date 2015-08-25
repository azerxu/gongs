// fastq package read fastq file

package fastq

import (
	"bytes"
	"errors"
	"fmt"
	"gongs/biofile"
	"gongs/scan"
	"gongs/xopen"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	ErrEmptyInputFile = errors.New("No Input Fastq File Given")
)

// Fastq (Name, Seq, Qual)
type Fastq struct {
	Name string
	Seq  []byte
	Qual []byte
}

func (fq Fastq) GetName() string {
	return fq.Name
}

func (fq Fastq) GetSeq() []byte {
	return fq.Seq
}

func (fq Fastq) GetQual() []byte {
	return fq.Qual
}

func (fq Fastq) String() string {
	return fmt.Sprintf("@%s\n%s\n+\n%s", fq.Name, string(fq.Seq), string(fq.Qual))
}

func (fq Fastq) IsFilter() bool {
	return strings.Contains(fq.Name, ":Y:")
}

func (fq Fastq) Id() string {
	if n := strings.IndexByte(fq.Name, ' '); n >= 0 {
		return fq.Name[:n]
	}

	// for old solexa data format
	if n := strings.IndexByte(fq.Name, '#'); n >= 0 {
		return fq.Name[:n]
	}

	return fq.Name
}

type FastqFile struct {
	Name  string
	file  io.ReadCloser
	s     *scan.Scanner
	name  string
	seq   []byte
	qual  []byte
	err   error
	stage int
}

func Open(filename string) (*FastqFile, error) {
	file, err := xopen.Xopen(filename)
	if err != nil {
		return nil, err
	}
	if filename == "-" {
		filename = "STDIN"
	}

	return &FastqFile{
		Name: filename,
		s:    scan.New(file),
		file: file,
	}, nil
}

func (ff *FastqFile) Close() error {
	return ff.file.Close()
}

func (ff *FastqFile) Err() error {
	if ff.err == nil || ff.err == io.EOF {
		if err := ff.s.Err(); err != nil {
			return err
		}
		return nil
	}
	return ff.err
}

func (ff *FastqFile) setErr(err error) {
	if ff.err == nil {
		ff.err = err
	}
}

func (ff *FastqFile) Next() bool {
	if ff.err != nil {
		return false
	}

	var line []byte
	for ff.s.Scan() {
		line = bytes.TrimSpace(ff.s.Bytes())
		if len(line) == 0 { // ingore empty line
			continue
		}
		switch ff.stage {
		case 0: // get fastq name
			if len(line) > 0 && line[0] != '@' {
				ff.setErr(fmt.Errorf("file: %v Wrong Fastq Record Name %s at line: %d", ff.Name, string(line), ff.s.Lid()))
				return false
			}
			ff.stage++
			ff.name = string(line[1:])
			ff.seq = ff.seq[:0]   // clear seq
			ff.qual = ff.qual[:0] // clear qual
		case 1: // get fastq seq
			if len(line) > 0 && line[0] == '+' {
				ff.stage += 2
				break
			}
			ff.seq = append(ff.seq, line...)
		case 2: // get + line
		case 3: // get fastq qual
			ff.qual = append(ff.qual, line...)
			if len(ff.qual) == len(ff.seq) {
				ff.stage = 0
				return true
			} else if len(ff.qual) > len(ff.seq) {
				ff.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual (%s) length (%d) != seq (%s) length (%d) at line: %d",
					ff.Name, string(ff.name), ff.qual, len(ff.qual), ff.seq, len(ff.seq), ff.s.Lid()))
				return false
			}
		}
	}
	if len(ff.qual) < len(ff.seq) {
		ff.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual length (%d) != seq length (%d) at line: %d",
			ff.Name, string(ff.name), len(ff.qual), len(ff.seq), ff.s.Lid()))
	}
	ff.setErr(io.EOF)
	return false
}

func (ff *FastqFile) Fq() *Fastq {
	return &Fastq{Name: ff.name, Seq: ff.seq, Qual: ff.qual}
}

func (ff *FastqFile) Value() (string, []byte, []byte) {
	return ff.name, ff.seq, ff.qual
}

func (ff *FastqFile) Iter() <-chan *Fastq {
	ch := make(chan *Fastq)
	go func(ch chan *Fastq) {
		for ff.Next() {
			ch <- ff.Fq()
		}
		close(ch)
	}(ch)
	return ch
}

func (ff *FastqFile) Seqs() <-chan biofile.Seqer {
	ch := make(chan biofile.Seqer)
	go func(ch chan biofile.Seqer) {
		for ff.Next() {
			ch <- ff.Fq()
		}
		close(ch)
	}(ch)
	return ch
}

func Opens(filenames ...string) ([]*FastqFile, error) {
	if len(filenames) == 0 {
		return nil, ErrEmptyInputFile
	}

	fqfiles := make([]*FastqFile, len(filenames))
	for i, filename := range filenames {
		fqfile, err := Open(filename)
		if err != nil {
			return nil, err
		}
		fqfiles[i] = fqfile
	}
	return fqfiles, nil
}

func Load(filenames ...string) (<-chan *Fastq, <-chan error) {
	fqChan := make(chan *Fastq, 2*len(filenames))
	errChan := make(chan error, 1)

	fqfiles, err := Opens(filenames...)
	if err != nil {
		errChan <- err
		return nil, errChan
	}

	go func(fqChan chan *Fastq, errChan chan error, fqfiles []*FastqFile) {
		for _, fqfile := range fqfiles {
			defer fqfile.Close()
			for fqfile.Next() {
				fqChan <- fqfile.Fq()
			}
			if err := fqfile.Err(); err != nil {
				errChan <- err
			}
		}
		close(fqChan)
		close(errChan)
	}(fqChan, errChan, fqfiles)
	return fqChan, errChan
}

func LoadMix(filenames ...string) (<-chan *Fastq, <-chan error) {
	fqChan := make(chan *Fastq, 2*len(filenames))
	errChan := make(chan error, 1)

	fqfiles, err := Opens(filenames...)
	if err != nil {
		errChan <- err
		return nil, errChan
	}

	go func(fqfiles []*FastqFile, fqChan chan *Fastq, errChan chan error) {
		wg := &sync.WaitGroup{}
		for _, fqfile := range fqfiles {
			wg.Add(1)
			go func(wg *sync.WaitGroup, fqChan chan *Fastq, errChan chan error, fqfile *FastqFile) {
				defer fqfile.Close()
				for fqfile.Next() {
					fmt.Fprintln(os.Stderr, fqfile.Fq())
					fqChan <- fqfile.Fq()
				}
				if err := fqfile.Err(); err != nil {
					errChan <- err
				}
				wg.Done()
			}(wg, fqChan, errChan, fqfile)
		}
		wg.Wait()
		close(fqChan)
		close(errChan)
	}(fqfiles, fqChan, errChan)
	return fqChan, errChan
}
