package qc

type Qual struct {
	stat map[byte]int
	pos  map[int]map[byte]int
	max  byte
	min  byte
}

func NewQual() *Qual {
	return &Qual{
		stat: make(map[byte]int),
		pos:  make(map[int]map[byte]int),
		min:  127,
	}
}

func (q *Qual) Count(qual string) {
	for i, qu := range qual {
		if max < qu {
			max = qu
		}
		if min > qu {
			min = qu
		}
		q.stat[qu]++
		mm, ok := q.pos[i]
		if !ok {
			mm = make(map[byte]int)
			q.pos[i] = mm
		}
		mm[qu]++
	}
}

func (q *Qual) Max() int {
	return int(q.max)
}

func (q *Qual) Min() int {
	return int(q.min)
}
