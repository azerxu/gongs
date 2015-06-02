package align

import (
	"fmt"
)

type LocalAligner struct {
	seq       string
	match     int
	mismatch  int
	gap       int
	wild      int
	errorRate float64
}

func (la LocalAligner) Align(target string) *AlignResult {
	return Local(la.seq, target, la.match, la.mismatch, la.gap, la.wild, la.errorRate)
}

func (la LocalAligner) AlignTo(query string) *AlignResult {
	return Local(query, la.seq, la.match, la.mismatch, la.gap, la.wild, la.errorRate)
}

func (la LocalAligner) String() string {
	return fmt.Sprintf("LocalAligner(Seq:%s, match:%d, mismatch:%d, gap:%d, wild:%d, errRate:%0.3f)",
		la.seq, la.match, la.mismatch, la.gap, la.wild, la.errorRate)
}

/**************************************** ALIGNMENT ***************************
                target (j)
            -----------------> n
           |
           |
query (i)  |
           |
           |
           V
           m

  query:        query string
  target:       target string
  method:       one of the LOCAL, GLOCAL, GLOBAL
  match:        match score [1]
  mismatch:     mismatch score [-1]
  gap:          gap penity score [-2]
  wild:         match N or not
  errorRate:    error / align_length [0.05]
******************************************************************************/

func Local(query, target string, match, mismatch, gap int, wild int, errorRate float64) *AlignResult {
	var temp cell
	var is_match, diag, up, left, score, errors, matchs, qstart, tstart int

	best_align := AlignResult{}
	qlen := len(query)
	tlen := len(target)
	rows := make([]cell, qlen+1)

	// init [0] row qstart
	for i := 0; i < qlen+1; i++ {
		rows[i] = cell{qstart: i}
	}

	for j := 1; j < tlen+1; j++ { // align each col one by one
		temp = rows[0]
		rows[0] = cell{tstart: j}

		for i := 1; i < qlen+1; i++ {
			if query[i-1] == target[j-1] || (wild&WILD_QUERY != 0 && query[i-1] == WILD_CHAR) || (wild&WILD_TARGET != 0 && target[j-1] == WILD_CHAR) {
				is_match = 1
				diag = temp.score + match
			} else {
				is_match = 0
				diag = temp.score + mismatch
			}

			up = rows[i-1].score + gap
			left = rows[i].score + gap

			if diag >= up && diag >= left { // match or mismatch
				score = diag
				matchs = temp.matchs + is_match
				errors = temp.errors + 1 - is_match
				qstart = temp.qstart
				tstart = temp.tstart
			} else if up >= left { // insert
				score = up
				errors = rows[i-1].errors + 1
				matchs = rows[i-1].matchs
				qstart = rows[i-1].qstart
				tstart = rows[i-1].tstart
			} else { // delete
				score = left
				errors = rows[i].errors + 1
				matchs = rows[i].matchs
				qstart = rows[i].qstart
				tstart = rows[i].tstart
			}

			temp = rows[i] // record current row for next diag compare

			// update current row
			if score < 0 { // reset alignment start if align score too low
				rows[i].score = 0
				rows[i].matchs = 0
				rows[i].errors = 0
				rows[i].qstart = i
				rows[i].tstart = j
			} else {
				rows[i].score = score
				rows[i].matchs = matchs
				rows[i].errors = errors
				rows[i].qstart = qstart
				rows[i].tstart = tstart

				// if float64(errors) <= float64(j-tstart)*errorRate &&
				if float64(errors) <= float64(j-tstart)*errorRate &&
					(score > best_align.Score || (score == best_align.Score && matchs >= best_align.Matchs)) {
					best_align.Matchs = matchs
					best_align.Score = score
					best_align.Errors = errors
					best_align.Qstart = qstart
					best_align.Qend = i
					best_align.Tstart = tstart
					best_align.Tend = j
				}
			}
		}
	}
	return &best_align
}
