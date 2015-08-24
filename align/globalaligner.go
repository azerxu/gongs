package align

import "fmt"

type GlobalAligner struct {
	seq       string
	match     int
	mismatch  int
	gap       int
	wild      int
	errorRate float64
}

func (ga GlobalAligner) Align(target string) *AlignResult {
	return Global(ga.seq, target, ga.match, ga.mismatch, ga.gap, ga.wild, ga.errorRate)
}

func (ga GlobalAligner) AlignTo(query string) *AlignResult {
	return Global(query, ga.seq, ga.match, ga.mismatch, ga.gap, ga.wild, ga.errorRate)
}

func (ga GlobalAligner) String() string {
	return fmt.Sprintf("GlobalAligner(Seq:%s, match:%d, mismatch:%d, gap:%d, wild:%d, errRate:%0.3f)",
		ga.seq, ga.match, ga.mismatch, ga.gap, ga.wild, ga.errorRate)
}

func Glocal(query, target string, match, mismatch, gap int, wild int, errorRate float64) *AlignResult {
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

			rows[i].score = score
			rows[i].matchs = matchs
			rows[i].errors = errors
			rows[i].qstart = qstart
			rows[i].tstart = tstart
		}
		// check the last row for the max row alignment
		if (float64(errors) <= errorRate*float64(j-tstart)) &&
			(score > best_align.Score || (score == best_align.Score && matchs >= best_align.Matchs)) {
			best_align.Matchs = matchs
			best_align.Score = score
			best_align.Errors = errors
			best_align.Qstart = qstart
			best_align.Qend = qlen
			best_align.Tstart = tstart
			best_align.Tend = j
		}
	}

	// check the last col for the max row alignment
	for i := 1; i < qlen; i++ {
		if (float64(rows[i].errors) <= errorRate*float64(tlen-rows[i].tstart)) &&
			(rows[i].score > best_align.Score || (rows[i].score == best_align.Score && rows[i].matchs >= best_align.Matchs)) {
			best_align.Matchs = rows[i].matchs
			best_align.Score = rows[i].score
			best_align.Errors = rows[i].errors
			best_align.Qstart = rows[i].qstart
			best_align.Qend = i
			best_align.Tstart = rows[i].tstart
			best_align.Tend = tlen
		}
	}
	return &best_align
}
