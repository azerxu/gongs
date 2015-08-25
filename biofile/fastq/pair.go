package fastq

import (
	"errors"
	"fmt"
	"gongs/biofile"
	"sync"
)

var (
	ErrUnPairInputFile = errors.New("Input Fastq File Not Paired")
)

type Pair struct {
	Read1 *Fastq
	Read2 *Fastq
}

func (p Pair) GetRead1() biofile.Seqer {
	return p.Read1
}

func (p Pair) GetRead2() biofile.Seqer {
	return p.Read2
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

func (pf *FastqPairFile) Value() (biofile.Seqer, biofile.Seqer) {
	return pf.ff1.Fq(), pf.ff2.Fq()
}

func (pf *FastqPairFile) Pair() *Pair {
	return &Pair{Read1: pf.ff1.Fq(), Read2: pf.ff2.Fq()}
}

func (pf *FastqPairFile) Iter() <-chan *Pair {
	out := make(chan *Pair)
	go func(pf *FastqPairFile, out chan *Pair) {
		for pf.Next() {
			out <- pf.Pair()
		}
		close(out)
	}(pf, out)
	return out
}

func (pf *FastqPairFile) Pairs() <-chan biofile.PairSeqer {
	out := make(chan biofile.PairSeqer)
	go func(pf *FastqPairFile, out chan biofile.PairSeqer) {
		for pf.Next() {
			out <- pf.Pair()
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

func OpenPairs(filenames ...string) ([]*FastqPairFile, error) {
	n := len(filenames)
	if n == 0 {
		return nil, ErrEmptyInputFile
	} else if n%2 != 0 {
		return nil, ErrUnPairInputFile
	}
	pfs := make([]*FastqPairFile, n%2)
	for i := 0; i < n; i += 2 {
		pf, err := OpenPair(filenames[i], filenames[i+1])
		if err != nil {
			return nil, err
		}
		pfs[i/2] = pf
	}
	return pfs, nil
}

func LoadPair(filenames ...string) (<-chan *Pair, <-chan error) {
	pChan := make(chan *Pair, len(filenames))
	errChan := make(chan error, 1)

	pfs, err := OpenPairs(filenames...)
	if err != nil {
		errChan <- err
		return nil, errChan
	}

	go func(pfs []*FastqPairFile, pChan chan *Pair, errChan chan error) {
		for _, pf := range pfs {
			defer pf.Close()
			for pf.Next() {
				pChan <- pf.Pair()
			}
			if err := pf.Err(); err != nil {
				errChan <- err
			}
		}
		close(pChan)
		close(errChan)
	}(pfs, pChan, errChan)
	return pChan, errChan
}

func LoadPairMix(filenames ...string) (<-chan *Pair, <-chan error) {
	pChan := make(chan *Pair, len(filenames))
	errChan := make(chan error, 1)

	pfs, err := OpenPairs(filenames...)
	if err != nil {
		errChan <- err
		return nil, errChan
	}

	go func(pfs []*FastqPairFile, pChan chan *Pair, errChan chan error) {
		wg := &sync.WaitGroup{}
		for _, pf := range pfs {
			wg.Add(1)
			go func(pf *FastqPairFile, wg *sync.WaitGroup, pChan chan *Pair, errChan chan error) {
				defer pf.Close()
				for pf.Next() {
					pChan <- pf.Pair()
				}
				if err := pf.Err(); err != nil {
					errChan <- err
				}
				wg.Done()
			}(pf, wg, pChan, errChan)
		}
		wg.Wait()
		close(pChan)
		close(errChan)
	}(pfs, pChan, errChan)
	return pChan, errChan
}
