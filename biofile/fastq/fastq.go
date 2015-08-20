package fastq

import (
	"bytes"
	"fmt"
	"gongs/scan"
	"gongs/xopen"
	"io"
	"os"
	"strings"
	"sync"
)

type Fastq struct {
	Name string
	Seq  []byte
	Qual []byte
}

func (fq Fastq) String() string {
	return fmt.Sprintf("@%v\n%v\n+\n%v", fq.Name, string(fq.Seq), string(fq.Qual))
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

func (fq Fastq) IsFilter() bool {
	return strings.Contains(fq.Name, ":Y:")
}

type FastqFile struct {
	Name  string
	file  io.ReadCloser
	s     *scan.Scanner
	name  []byte
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
			ff.name = line[1:]
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
				ff.setErr(fmt.Errorf("file: %v Fastq Record (%s) qual length (%d) != seq length (%d) at line: %d",
					ff.Name, string(ff.name), len(ff.qual), len(ff.seq), ff.s.Lid()))
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

func (ff *FastqFile) Value() *Fastq {
	return &Fastq{Name: string(ff.name), Seq: ff.seq, Qual: ff.qual}
}

func (ff *FastqFile) Iter() <-chan *Fastq {
	ch := make(chan *Fastq)
	go func(ch chan *Fastq) {
		for ff.Next() {
			ch <- ff.Value()
		}
		close(ch)
	}(ch)
	return ch
}

func Opens(filenames ...string) ([]*FastqFile, error) {
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
				fqChan <- fqfile.Value()
			}
			if err := fqfile.Err(); err != nil {
				errChan <- err
			}
		}
		close(fqChan)
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
					fmt.Fprintln(os.Stderr, fqfile.Value())
					fqChan <- fqfile.Value()
				}
				if err := fqfile.Err(); err != nil {
					errChan <- err
				}
				wg.Done()
			}(wg, fqChan, errChan, fqfile)
		}
		wg.Wait()
		close(fqChan)
	}(fqfiles, fqChan, errChan)
	return fqChan, errChan
}
