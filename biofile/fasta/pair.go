package fasta

import (
	"errors"
	"gongs/xopen"
	"sync"
)

var (
	ErrEmptyInputFile  = errors.New("No Input Fasta File Given")
	ErrUnPairInputFile = errors.New("Input Fasta File Not Paired")
)

type Pair struct {
	Read1 *Fasta
	Read2 *Fasta
}

type FastaPairFile struct {
	ff1 *FastaFile
	ff2 *FastaFile
}

func (pf *FastaPairFile) Close() error {
	if err := pf.ff1.Close(); err != nil {
		return err
	}
	if err := pf.ff2.Close(); err != nil {
		return err
	}
	return nil
}

func (pf *FastaPairFile) Err() error {
	if err := pf.ff1.Err(); err != nil {
		return err
	} else if err := pf.ff2.Err(); err != nil {
		return err
	}
	return nil
}

func (pf *FastaPairFile) Filenames() (string, string) {
	return pf.ff1.Name, pf.ff2.Name
}

func (pf *FastaPairFile) Next() bool {
	if pf.ff1.Next() && pf.ff2.Next() {
		return true
	}
	return false
}

func (pf *FastaPairFile) Value() *Pair {
	return &Pair{Read1: pf.ff1.Value(), Read2: pf.ff2.Value()}
}

func (pairfile *FastaPairFile) Read() chan *Pair {
	chan1 := pairfile.file1.Read()
	chan2 := pairfile.file2.Read()

	out := make(chan *Pair)
	go func(chan1, chan2 <-chan *Fasta, out chan *Pair) {
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

func OpenPair(filename1, filename2 string) (*FastaPairFile, error) {
	file1, err := xopen.Xopen(filename1)
	if err != nil {
		return nil, err
	}
	file2, err := xopen.Xopen(filename2)
	if err != nil {
		return nil, err
	}
	return &FastaPairFile{
		ff1: file1,
		ff2: file2,
	}, nil
}

func OpenPairs(filenames ...string) ([]*FastaPairFile, error) {
	if l := len(filenames); l == 0 {
		return nil, ErrEmptyInputFile
	} else if l%2 != 0 {
		return nil, ErrUnPairInputFile
	}
	pfs := make([]*FastaPairFile, len(filenames)/2)
	for i := 0; i < len(filenames); i += 2 {
		filename1 := filenames[i]
		filename2 := filenames[i+1]
		pf, err := OpenPair(filename1, filename2)
		if err != nil {
			return nil, err
		}
		pfs[i/2] = pf
	}
	return pfs, nil
}

func Iter(filenames ...string) (<-chan *Pair, error) {
	pfs, err := OpenPairs(filenames...)
	if err != nil {
		return nil, err
	}

	ch := make(chan *Pair, 4*len(pfs))
	go func(ch chan *Pair, pfs []*FastaPairFile) {
		wg := &sync.WaitGroup{}
		wg.Add(len(pfs))
		for _, pf := range pfs {
			go func(ch chan *Pair, pf *FastqPairFile, wg *sync.WaitGroup) {
				defer pf.Close()
				for pf.Next() {
					ch <- pf.Value()
				}
				wg.Done()
			}(ch, pf, wg)
		}
		wg.Wait()
		close(ch)
	}(ch, pfs)
	return ch, nil
}
