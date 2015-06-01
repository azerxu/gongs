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

type Fastq struct {
	Name string
	Seq  string
	Qual string
}

func (fq Fastq) String() string {
	return fmt.Sprintf("@%s\n%s\n+\n%s", fq.Name, fq.Seq, fq.Qual)
}

type FastqFile struct {
	Filename string
	file     io.ReadCloser
}

func (fqfile *FastqFile) Read() <-chan *Fastq {
	ch := make(chan *Fastq)
	go run(fqfile.Filename, fqfile.file, ch)
	return ch
}

func (fqfile *FastqFile) Close() error {
	return fqfile.file.Close()
}

func Open(filename string) (*FastqFile, error) {
	file, err := lib.Xopen(filename)
	if err != nil {
		return nil, err
	}

	if filename == "-" || filename == "" {
		filename = "STDIN"
	}

	fqfile := &FastqFile{
		Filename: filename,
		file:     file,
	}
	return fqfile, nil
}

func Read(filename string) (<-chan *Fastq, error) {
	fqfile, err := Open(filename)
	if err != nil {
		return nil, err
	}
	return fqfile.Read(), nil
}

// Parse(file io.Reader, filename string)
func Parse(args ...interface{}) <-chan *Fastq {
	var file io.Reader
	var filename string
	switch l := len(args); {
	case l >= 2:
		filename = args[1].(string)
		fallthrough
	case l == 1:
		file = args[0].(io.Reader)
	default:
		file = os.Stdin
	}
	ch := make(chan *Fastq)
	go run(filename, file, ch)
	return ch
}

func run(filename string, file io.Reader, ch chan *Fastq) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
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

func SimpleRead(filename string) (<-chan *Fastq, error) {
	fqfile, err := Open(filename)
	if err != nil {
		return nil, err
	}
	ch := make(chan *Fastq)
	go func(fqfile *FastqFile, ch chan *Fastq) {
		ln, name, seq := 0, "", ""
		handle := bufio.NewScanner(fqfile.file)
		for handle.Scan() {
			line := strings.TrimSpace(handle.Text())
			if len(line) == 0 {
				continue
			}

			switch ln % 4 {
			case 0:
				name = line[1:]
			case 1:
				seq = line
			case 3:
				ch <- &Fastq{Name: name, Seq: seq, Qual: line}
			}
			ln++
		}
		close(ch)
	}(fqfile, ch)
	return ch, nil
}

type Pair struct {
	Read1 *Fastq
	Read2 *Fastq
}

type FastqPairFile struct {
	Filename1 string
	Filename2 string
	file1     *FastqFile
	file2     *FastqFile
}

func (pairfile *FastqPairFile) Close() error {
	if err := pairfile.file1.Close(); err != nil {
		return err
	}
	if err := pairfile.file2.Close(); err != nil {
		return err
	}
	return nil
}

func (pairfile *FastqPairFile) Read() chan *Pair {
	chan1 := pairfile.file1.Read()
	chan2 := pairfile.file2.Read()

	out := make(chan *Pair)
	go func(chan1, chan2 <-chan *Fastq, out chan *Pair) {
		for {
			read1, ok1 := <-chan1
			read2, ok2 := <-chan2
			if ok1 && ok2 {
				out <- &Pair{Read1: read1, Read2: read2}
			} else {
				break
			}
		}
		close(out)
	}(chan1, chan2, out)
	return out
}

func OpenPair(filename1, filename2 string) (*FastqPairFile, error) {
	file1, err := Open(filename1)
	if err != nil {
		return nil, err
	}
	file2, err := Open(filename2)
	if err != nil {
		return nil, err
	}
	return &FastqPairFile{
		Filename1: filename1,
		Filename2: filename2,
		file1:     file1,
		file2:     file2,
	}, nil
}

func Load(filenames ...string) (chan *Fastq, error) {
	l := len(filenames)
	if l == 0 {
		return nil, fmt.Errorf("no fastq file given%s", "!!!")
	}

	fqfiles := []*FastqFile{}
	for i := 0; i < l; i++ {
		filename := filenames[i]
		fqfile, err := Open(filename)
		if err != nil {
			return nil, err
		}
		fqfiles = append(fqfiles, fqfile)
	}

	ch := make(chan *Fastq)
	go func(ch chan *Fastq, fqfiles []*FastqFile) {
		wg := &sync.WaitGroup{}
		wg.Add(len(fqfiles))
		for _, fqfile := range fqfiles {
			go func(ch chan *Fastq, fqfile *FastqFile, wg *sync.WaitGroup) {
				defer fqfile.Close()
				for fq := range fqfile.Read() {
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

func LoadPair(filenames ...string) (chan *Pair, error) {
	if l := len(filenames); l == 0 {
		return nil, fmt.Errorf("no fastq file given%s", "!!!")
	} else if l%2 != 0 {
		return nil, fmt.Errorf("input fastq file not paired, given %d files.", l)
	}

	pairfiles := []*FastqPairFile{}
	for i := 0; i < len(filenames); i += 2 {
		filename1 := filenames[i]
		filename2 := filenames[i+1]
		pairfile, err := OpenPair(filename1, filename2)
		if err != nil {
			return nil, err
		}
		pairfiles = append(pairfiles, pairfile)
	}

	ch := make(chan *Pair)
	go func(ch chan *Pair, pairfiles []*FastqPairFile) {
		wg := &sync.WaitGroup{}
		wg.Add(len(pairfiles))
		for _, pairfile := range pairfiles {
			go func(ch chan *Pair, pairfile *FastqPairFile, wg *sync.WaitGroup) {
				defer pairfile.Close()
				for pair := range pairfile.Read() {
					ch <- pair
				}
				wg.Done()
			}(ch, pairfile, wg)
		}
		wg.Wait()
		close(ch)
	}(ch, pairfiles)
	return ch, nil
}
