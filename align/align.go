// ****************************** align using smith-waterman algor to alignment two sequence **************************
// align package contain function: Align, Local, Glocal, Global
// require parameters::
//   query      string
//   target     string
//   method     int (only in Align)
//   match      int
//   mismatch   int
//   gap        int
//   wild       int
//   errorRate  float64

package align

import "fmt"

// ****************************** Default Const Setting ***********************

const DIR_NONE int = 0
const DIR_LEFT int = 1
const DIR_UP int = 2
const DIR_DIAG int = 3

const WILD_CHAR byte = 'N'
const WILD_NONE int = 0
const WILD_QUERY int = 1
const WILD_TARGET int = 2
const WILD_ALL int = 3

// const LOCAL int = 0
// const GLOCAL int = 1
// const GLOBAL int = 2

// ****************************** Data Structure Setting **********************

type cell struct { // record each cell in matrix
	qstart int // query start position
	tstart int // target start position
	score  int // alignment score
	matchs int // alignment match number
	errors int // error number (mismatch, insert, delete)
}

type AlignResult struct { // record alignment result
	Qstart int // start position at query
	Qend   int // end position at query
	Tstart int // start position at target
	Tend   int // end position at target
	Score  int // the score value of the alignment
	Matchs int // match numbers of the alignment
	Errors int // error numbers of the alignment
}

func (ar AlignResult) String() string {
	return fmt.Sprintf("AlignResult(qstart:%d, qend:%d, tstart:%d, tend:%d, score:%d, matchs:%d, errors:%d)",
		ar.Qstart, ar.Qend, ar.Tstart, ar.Tend, ar.Score, ar.Matchs, ar.Errors)
}

type Aligner interface {
	Align(string) *AlignResult
	AlignTo(string) *AlignResult
	String() string
}

func Align(method, query, target string, match, mismatch, gap int, wild int, errorRate float64) *AlignResult {
	switch method {
	case "LOCAL":
		return Local(query, target, match, mismatch, gap, wild, errorRate)
	case "GLOCAL":
		return Glocal(query, target, match, mismatch, gap, wild, errorRate)
	case "GLOBAL":
		return Global(query, target, match, mismatch, gap, wild, errorRate)
	}
	return nil
}

// init Aligner(name, match, mismatch, gap, wild, errRate)
func New(name, seq string, args ...int) (Aligner, error) {
	if seq == "" {
		return nil, fmt.Errorf("Alinger Sequence is empyty%s", "!")
	}

	match := 1
	mismatch := -1
	gap := -2
	wild := WILD_QUERY
	errRate := 1
	switch l := len(args); {
	case l > 4:
		errRate = args[4]
		fallthrough
	case l > 3:
		wild = args[3]
		fallthrough
	case l > 2:
		gap = args[2]
		fallthrough
	case l > 1:
		mismatch = args[1]
		fallthrough
	case l > 0:
		match = args[0]
	}

	switch name {
	case "local":
		return &LocalAligner{
			seq:      seq,
			match:    match,
			mismatch: mismatch,
			gap:      gap, wild: wild,
			errorRate: 0.01 * float64(errRate)}, nil
	case "global":
		return &GlobalAligner{
			seq:       seq,
			match:     match,
			mismatch:  mismatch,
			gap:       gap,
			wild:      wild,
			errorRate: 0.01 * float64(errRate)}, nil
	case "glocal":
		return &GlocalAligner{
			seq:       seq,
			match:     match,
			mismatch:  mismatch,
			gap:       gap,
			wild:      wild,
			errorRate: 0.01 * float64(errRate)}, nil
	}
	return nil, fmt.Errorf("Unkown Aligner Name: %s", name)
}
