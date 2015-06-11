package fastq

import (
	"bufio"
	"fmt"
	"gongs/lib"
	"io"
	"os"
	"strings"
	"sync"
)

// Fastq (Name, Seq, Qual)
type Fastq struct {
	Name string
	Seq  string
	Qual string
}

func (fq Fastq) String() string {
	return fmt.Sprintf("@%s\n%s\n+\n%s", fq.Name, fq.Seq, fq.Qual)
}

// File (filename)
type File struct {
	Name string
	f    io.ReadCloser
}

// Load return *Fastq chan
func (fqfile *File) Load() <-chan *Fastq {
	ch := make(chan *Fastq)
	go run(fqfile.Name, fqfile.f, ch)
	return ch
}

// Close File
func (fqfile *File) Close() error {
	return fqfile.f.Close()
}

// Open open fastqfile return *File, err
func Open(filename string) (*File, error) {
	file, err := lib.Xopen(filename)
	if err != nil {
		return nil, err
	}

	if filename == "-" || filename == "" {
		filename = "STDIN"
	}

	fqfile := &File{Name: filename, f: file}
	return fqfile, nil
}

// // Parse(file io.Reader, filename string)
// func Parse(args ...interface{}) <-chan *Fastq {
// 	var file io.Reader
// 	var filename string
// 	switch l := len(args); {
// 	case l >= 2:
// 		filename = args[1].(string)
// 		fallthrough
// 	case l == 1:
// 		file = args[0].(io.Reader)
// 	default:
// 		file = os.Stdin
// 	}
// 	ch := make(chan *Fastq)
// 	go run(filename, file, ch)
// 	return ch
// }

func run(filename string, file io.Reader, ch chan *Fastq) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	line, name, seq, qual, slen, qlen := "", "", []string{}, []string{}, 0, 0
	handle := bufio.NewScanner(file)
	for handle.Scan() {
		line = strings.TrimSpace(handle.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if line[0] != '@' {
			panic(fmt.Errorf("%s is not a fastq format file", filename))
		}
		break
	}

	if len(line) == 0 { // nothing input, empty fastq file
		close(ch)
		return
	}

	name = line[1:]
	isSeqBlock := true

	for handle.Scan() {
		line = strings.TrimSpace(handle.Text()) // trim space from both side

		if len(line) == 0 { // ignore empty line
			continue
		}

		if isSeqBlock {
			if line[0] == '+' {
				isSeqBlock = false
			} else {
				seq = append(seq, line)
				slen += len(line)
			}
		} else {
			if qlen > slen {
				panic(fmt.Errorf("Error: while Parsing fastq Record(%s) at file(%s)", name, filename))
			}

			if line[0] == '@' { // at beginning of next fastq
				if slen == qlen {
					ch <- &Fastq{Name: name, Seq: strings.Join(seq, ""), Qual: strings.Join(qual, "")}
					name, seq, qual, slen, qlen = line[1:], []string{}, []string{}, 0, 0
					isSeqBlock = true
				} else { // just a qual line begin with @
					qual = append(qual, line)
					qlen += len(line)
				}
			} else { // quality block
				qual = append(qual, line)
				qlen += len(line)
			}
		}
	}

	if len(name) != 0 || len(seq) != 0 { // check the last record
		if slen != qlen {
			panic(fmt.Errorf("Error: while Parsing fastq Record(%s) at file(%s)", name, filename))
		}
		ch <- &Fastq{Name: name, Seq: strings.Join(seq, ""), Qual: strings.Join(qual, "")}
	}
	close(ch)
}

// Load return *fastq chan, error
func Load(filename string) (<-chan *Fastq, error) {
	fqfile, err := Open(filename)
	if err != nil {
		return nil, err
	}
	return fqfile.Load(), nil
}

// Loads load multi-fastq file return *fastq chan, error
func Loads(filenames ...string) (<-chan *Fastq, error) {
	l := len(filenames)
	if l == 0 {
		return nil, fmt.Errorf("no fastq file given%s", "!!!")
	}

	fqfiles := []*File{}
	for i := 0; i < l; i++ {
		filename := filenames[i]
		fqfile, err := Open(filename)
		if err != nil {
			return nil, err
		}
		fqfiles = append(fqfiles, fqfile)
	}

	ch := make(chan *Fastq)
	go func(ch chan *Fastq, fqfiles []*File) {
		wg := &sync.WaitGroup{}
		for _, fqfile := range fqfiles {
			wg.Add(1)
			go func(ch chan *Fastq, fqfile *File, wg *sync.WaitGroup) {
				defer fqfile.Close()
				for fq := range fqfile.Load() {
					ch <- fq
				}
				wg.Done()
			}(ch, fqfile, wg)
		}
		wg.Wait()
		close(ch)
	}(ch, fqfiles)
	return ch, nil
}
