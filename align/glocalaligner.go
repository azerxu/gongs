package align

import "fmt"

type GlocalAligner struct {
	seq       string
	match     int
	mismatch  int
	gap       int
	wild      int
	errorRate float64
}

func (ca GlocalAligner) Align(target string) *AlignResult {
	return Glocal(ca.seq, target, ca.match, ca.mismatch, ca.gap, ca.wild, ca.errorRate)
}

func (ca GlocalAligner) AlignTo(query string) *AlignResult {
	return Glocal(query, ca.seq, ca.match, ca.mismatch, ca.gap, ca.wild, ca.errorRate)
}

func (ca GlocalAligner) String() string {
	return fmt.Sprintf("GlocalAligner(Seq:%s, match:%d, mismatch:%d, gap:%d, wild:%d, errRate:%0.3f)",
		ca.seq, ca.match, ca.mismatch, ca.gap, ca.wild, ca.errorRate)
}

func Global(query, target string, match, mismatch, gap int, wild int, errorRate float64) *AlignResult {
	var temp cell
	var is_match, diag, up, left, score, errors, matchs, qstart, tstart int

	best_align := AlignResult{}
	qlen := len(query)
	tlen := len(target)
	rows := make([]cell, qlen+1)

	// init [0] row qstart
	for i := 0; i < qlen+1; i++ {
		rows[i] = cell{qstart: i, score: i * gap}
	}

	for j := 1; j < tlen+1; j++ { // align each col one by one
		temp = rows[0]
		rows[0] = cell{tstart: j, score: j * gap}

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
			rows[i].score = score
			rows[i].matchs = matchs
			rows[i].errors = errors
			rows[i].qstart = qstart
			rows[i].tstart = tstart
		}
	}
	// check the last cell
	if float64(rows[qlen].errors) <= float64(tlen)*errorRate {
		best_align.Matchs = matchs
		best_align.Score = score
		best_align.Errors = errors
		best_align.Qstart = 0
		best_align.Qend = qlen
		best_align.Tstart = 0
		best_align.Tend = tlen
	}

	return &best_align
}
