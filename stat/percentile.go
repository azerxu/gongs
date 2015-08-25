// Refer from http://code.activestate.com/recipes/511478/

package stat

import (
	"math"
	"sort"
)

type IntSlice struct {
	Data []int
}

func NewIntSlice(data []int) *IntSlice {
	return &IntSlice{Data: data}
}

// Sort sort data
func (s *IntSlice) Sort() {
	sort.Ints(s.Data)
}

func (s *IntSlice) Percentile(p float64) float64 {
	if len(s.Data) == 0 {
		return 0
	}

	k := float64(len(s.Data)-1) * p
	f := math.Floor(k)
	c := math.Ceil(k)
	if f == c {
		return float64(s.Data[int(k)])
	}
	d0 := float64(s.Data[int(f)]) * (c - k)
	d1 := float64(s.Data[int(c)]) * (k - f)
	return d0 + d1
}

func (s *IntSlice) Sum() int {
	sum := 0
	for i, n := 0, len(s.Data); i < n; i++ {
		sum += s.Data[i]
	}
	return sum
}

func (s *IntSlice) Mean() float64 {
	n := len(s.Data)
	if n == 0 {
		return 0
	}
	return float64(s.Sum()) / float64(n)
}

func (s *IntSlice) Median() float64 {
	return s.Percentile(0.5)
}

type IntMap struct {
	Data map[int]int
	keys []int
	vals []int
}

func NewIntMap(data map[int]int) *IntMap {
	return &IntMap{Data: data}
}

// Keys return sorted keys
func (m *IntMap) Keys() []int {
	if m.keys != nil {
		return m.keys
	}
	keys := []int{}
	for key := range m.Data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	m.keys = keys
	return keys
}

func (m *IntMap) Vals() []int {
	if m.vals != nil {
		return m.vals
	}
	vals := make([]int, len(m.Data))
	for i, key := range m.Keys() {
		vals[i] = m.Data[key]
	}
	m.vals = vals
	return vals
}

func (m *IntMap) Items() int {
	n := 0 // total items number
	for _, v := range m.Data {
		n += v
	}
	return n
}

func (m *IntMap) percentKey(i int) int {
	count := -1
	for _, key := range m.Keys() {
		count += m.Data[key]
		if count >= i {
			return key
		}
	}
	return 0
}

func (m *IntMap) Percentile(p float64) float64 {
	if len(m.Data) == 0 {
		return 0
	}
	n := m.Items()
	k := float64(n-1) * p
	f := math.Floor(k)
	c := math.Ceil(k)
	if f == c {
		return float64(m.percentKey(int(k)))
	}

	d0 := float64(m.percentKey(int(f))) * (c - k)
	d1 := float64(m.percentKey(int(c))) * (k - f)
	return d0 + d1
}

func (m *IntMap) Sum() int {
	sum := 0
	for _, key := range m.Keys() {
		sum += key * m.Data[key]
	}
	return sum
}

func (m *IntMap) Mean() float64 {
	if len(m.Data) == 0 {
		return 0
	}
	return float64(m.Sum()) / float64(m.Items())
}

func (m *IntMap) Median() float64 {
	return m.Percentile(0.5)
}
