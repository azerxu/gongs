package fastq

import (
	"fmt"
	"sync"
)

type Pair struct {
	Read1 *Fastq
	Read2 *Fastq
}

func (p Pair) String() string {
	return fmt.Sprintf("%s\n%s", p.Read1, p.Read2)
}

type FastqPairFile struct {
	Name1 string
	Name2 string
	file1 *FastqFile
	file2 *FastqFile
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

func (pairfile *FastqPairFile) Load() <-chan *Pair {
	chan1 := pairfile.file1.Load()
	chan2 := pairfile.file2.Load()

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
		Name1: filename1,
		Name2: filename2,
		file1: file1,
		file2: file2,
	}, nil
}

func LoadPair(filename1, filename2 string) (<-chan *Pair, error) {
	pairfile, err := OpenPair(filename1, filename2)
	if err != nil {
		return nil, err
	}
	return pairfile.Load(), nil
}

func LoadPairs(filenames ...string) (<-chan *Pair, error) {
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
				for pair := range pairfile.Load() {
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
