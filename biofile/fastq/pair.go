package fastq

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrEmptyInputFile  = errors.New("No Input Fastq File Given")
	ErrUnPairInputFile = errors.New("Input Fastq File Not Paired")
)

type Pair struct {
	Read1 *Fastq
	Read2 *Fastq
}

func (p Pair) String() string {
	return fmt.Sprintf("%s\n%s", p.Read1, p.Read2)
}

type FastqPairFile struct {
	ff1 *FastqFile
	ff2 *FastqFile
	err error
}

func (pf *FastqPairFile) Filenames() (string, string) {
	return pf.ff1.Name, pf.ff2.Name
}

func (pf *FastqPairFile) Err() error {
	if err := pf.ff1.Err(); err != nil {
		return err
	} else if err := pf.ff2.Err(); err != nil {
		return err
	}
	return pf.err
}

func (pf *FastqPairFile) Close() error {
	if err := pf.ff1.Close(); err != nil {
		return err
	}
	if err := pf.ff2.Close(); err != nil {
		return err
	}
	return nil
}

func (pf *FastqPairFile) Next() bool {
	if pf.ff1.Next() && pf.ff2.Next() {
		return true
	}
	return false
}

func (pf *FastqPairFile) Value() *Pair {
	return &Pair{
		Read1: pf.ff1.Value(),
		Read2: pf.ff2.Value(),
	}
}

func (pf *FastqPairFile) Iter() <-chan *Pair {
	out := make(chan *Pair)
	go func(pf *FastqPairFile, out chan *Pair) {
		for pf.Next() {
			out <- pf.Value()
		}
		close(out)
	}(pf, out)
	return out
}

func OpenPair(filename1, filename2 string) (*FastqPairFile, error) {
	ff1, err := Open(filename1)
	if err != nil {
		return nil, err
	}
	ff2, err := Open(filename2)
	if err != nil {
		return nil, err
	}
	return &FastqPairFile{
		ff1: ff1,
		ff2: ff2,
	}, nil
}

func LoadPair(filename1, filename2 string) (<-chan *Pair, error) {
	pf, err := OpenPair(filename1, filename2)
	if err != nil {
		return nil, err
	}
	return pf.Iter(), nil
}

func OpenPairs(filenames ...string) ([]*FastqPairFile, error) {
	if l := len(filenames); l == 0 {
		return nil, ErrEmptyInputFile
	} else if l%2 != 0 {
		return nil, ErrUnPairInputFile
	}
	pfs := make([]*FastqPairFile, len(filenames)/2)
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
	go func(ch chan *Pair, pfs []*FastqPairFile) {
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
