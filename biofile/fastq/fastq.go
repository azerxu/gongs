package fastq

import (
	"fmt"
	"gongs/scan"
	"gongs/xopen"
	"io"
	"strings"
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

type FastqFile struct {
	Name  string
	file  io.ReadCloser
	s     *scan.Scanner
	name  []byte
	seq   []byte
	qual  []byte
	buf   []byte
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
		buf:  make([]byte, 1024),
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
		line = ff.s.Bytes()
		switch ff.stage {
		case 0: // get fastq name
			if line[0] != '@' {
				ff.setErr(fmt.Errorf("file: %v Wrong Fastq Record Name %s at line: %d", ff.Name, string(line), ff.s.Lid()))
				return false
			}
			ff.stage++
			ff.name = line[1:]
			ff.seq = ff.seq[:0]   // clear seq
			ff.qual = ff.qual[:0] // clear qual
		case 1: // get fastq seq
			if line[0] == '+' {
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
	return false
}

func (ff *FastqFile) Value() *Fastq {
	return &Fastq{Name: string(ff.name), Seq: ff.seq, Qual: ff.qual}
}

func (ff *FastqFile) Iter() <-chan *Fastq {
	ch := make(chan *Fastq)
<<<<<<< HEAD
	go func(ch chan *Fastq, fqfiles []*FastqFile) {
		wg := &sync.WaitGroup{}
		for _, fqfile := range fqfiles {
			wg.Add(1)
			go func(ch chan *Fastq, fqfile *FastqFile, wg *sync.WaitGroup) {
				defer fqfile.Close()
				for fq := range fqfile.Load() {
					ch <- fq
				}
				wg.Done()
			}(ch, fqfile, wg)
=======
	go func(ch chan *Fastq) {
		for ff.Next() {
			ch <- ff.Value()
>>>>>>> align
		}
		close(ch)
	}(ch)
	return ch
}