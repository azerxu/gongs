// package for read line
// To fix bufio.Scanner meeting error when fasta seq length too long

package scan

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrNegativeAdvance = errors.New("Liner: split returns negative advance count")
	ErrAdvanceTooFar   = errors.New("Liner: split returns advance count beyond input")
)

type Scanner struct {
	r     io.Reader
	token []byte
	buf   []byte
	start int
	end   int
	err   error
	lid   int
}

func New(r io.Reader) *Scanner {
	return &Scanner{
		r:   r,
		buf: make([]byte, 4096),
	}
}

func (s Scanner) Bytes() []byte {
	return s.token
}

func (s Scanner) Line() string {
	return string(s.token)
}

func (s Scanner) Err() error {
	if s.err == nil || s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s Scanner) Lid() int {
	return s.lid
}

func (s *Scanner) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

func (s *Scanner) Scan() bool {
	for {
		// check buf
		if s.end > s.start || s.err != nil {
			advance, token, err := splitLine(s.buf[s.start:s.end], s.err != nil)
			if err != nil {
				s.setErr(err)
				return false
			}
			if !s.checkAdvance(advance) {
				return false
			}
			s.token = token
			if token != nil {
				s.lid++
				return true
			}
		}
		// end of file or meet an error
		if s.err != nil {
			s.start = 0
			s.end = 0
			return false
		}

		// now need read more date
		// first shift data to the beginning of buf
		if s.start > 0 && (s.end == len(s.buf)) || s.start > len(s.buf)/2 {
			copy(s.buf, s.buf[s.start:s.end])
			s.end -= s.start
			s.start = 0
		}

		// check buf is full, if so, double size buf
		if s.end == len(s.buf) {
			newSize := 2 * len(s.buf)
			newbuf := make([]byte, newSize)
			copy(newbuf, s.buf[s.start:s.end])
			s.buf = newbuf
			s.end -= s.start
			s.start = 0
		}

		n, err := s.r.Read(s.buf[s.end:len(s.buf)])
		s.end += n
		if err != nil {
			s.setErr(err)
			return false
		}
	}
}

func (s *Scanner) checkAdvance(n int) bool {
	if n < 0 {
		s.setErr(ErrNegativeAdvance)
		return false
	}
	if n > s.end-s.start {
		s.setErr(ErrAdvanceTooFar)
		return false
	}
	s.start += n
	return true
}

func dropCR(buf []byte) []byte {
	if len(buf) > 0 && buf[len(buf)-1] == '\r' {
		return buf[0 : len(buf)-1]
	}
	return buf
}

func splitLine(buf []byte, eof bool) (int, []byte, error) {
	if eof && len(buf) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(buf, '\n'); i >= 0 {
		return i + 1, dropCR(buf[0:i]), nil
	}

	if eof {
		return len(buf), dropCR(buf), nil
	}
	return 0, nil, nil
}
