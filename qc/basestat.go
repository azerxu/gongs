package qc

type Base struct {
	stat map[byte]int // A,C,G,T,N
}

func (b *Base) Count(seq string) {
	for _, nt := range seq {
		b.stat[nt]++
	}
}

func (b *Base) GC() float64 {
	gc := b.G() + b.C()
	tot := b.A() + b.T() + gc
	return float64(gc*100) / float64(tot)
}

func (b *Base) Total() int {
	for key, val := range b.stat {
		if key == 'A' || key == 'a' || key == 'C' || key == 'c' || key == 'G' || key == 'g' || key == 'T' || key == 't' {
			tot += val
		}
	}
	return tot
}

func (b *Base) TotalAll() int {
	tot := 0
	for _, val := range b.stat {
		tot += val
	}
	return tot
}

func (b *Base) A() int {
	return b.stat['A'] + b.stat['a']
}

func (b *Base) C() int {
	return b.stat['C'] + b.stat['c']
}

func (b *Base) G() int {
	return b.stat['G'] + b.stat['g']
}

func (b *Base) T() int {
	return b.stat['T'] + b.stat['t']
}

func (b *Base) N() int {
	return b.stat['N'] + b.stat['n']
}
