package stat

import "testing"

func TestIntSlice(t *testing.T) {
	s := NewIntSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	s.Sort()
	if sum := s.Sum(); sum != 45 {
		t.Errorf("Sum expect: %v get: %v", 45, sum)
	}
	if mean := s.Mean(); mean != 5 {
		t.Errorf("Mean expect: %v get: %v", 5, mean)
	}
	if median := s.Median(); median != 5 {
		t.Errorf("Median expect: %v get: %v", 5, median)
	}
	if q1 := s.Percentile(.25); q1 != 3 {
		t.Errorf("quantile1 expect: %v get: %v", 3, q1)
	}
	if q2 := s.Percentile(.5); q2 != 5 {
		t.Errorf("quantile2 expect: %v get: %v", 5, q2)
	}
	if q3 := s.Percentile(.75); q3 != 7 {
		t.Errorf("quantile3 expect: %v get: %v", 7, q3)
	}
	if p95 := s.Percentile(.95); p95 != 8.60 {
		t.Errorf("percentile95 expect: %v get: %v", 8.60, p95)
	}
}

func TestIntMap(t *testing.T) {
	m := NewIntMap(map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1})

	if sum := m.Sum(); sum != 45 {
		t.Errorf("Sum expect: %v get: %v", 45, sum)
	}
	if mean := m.Mean(); mean != 5 {
		t.Errorf("Mean expect: %v get: %v", 5, mean)
	}
	if median := m.Median(); median != 5 {
		t.Errorf("Median expect: %v get: %v", 5, median)
	}
	if q1 := m.Percentile(.25); q1 != 3 {
		t.Errorf("quantile1 expect: %v get: %v", 3, q1)
	}
	if q2 := m.Percentile(.5); q2 != 5 {
		t.Errorf("quantile2 expect: %v get: %v", 5, q2)
	}
	if q3 := m.Percentile(.75); q3 != 7 {
		t.Errorf("quantile3 expect: %v get: %v", 7, q3)
	}
	if p95 := m.Percentile(.95); p95 != 8.60 {
		t.Errorf("percentile95 expect: %v get: %v", 8.60, p95)
	}

}

func TestIntSlice2(t *testing.T) {
	s := NewIntSlice([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6,
		7, 7, 7, 7, 7, 7, 7, 8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 9, 9, 9})

	s.Sort()
	if sum := s.Sum(); sum != 285 {
		t.Errorf("Sum expect: %v get: %v", 285, sum)
	}
	if mean := s.Mean(); mean != 6.3333333333333333 {
		t.Errorf("Mean expect: %v get: %v", 6.33333333333333, mean)
	}
	if median := s.Median(); median != 7 {
		t.Errorf("Median expect: %v get: %v", 7, median)
	}
	if q1 := s.Percentile(.25); q1 != 5 {
		t.Errorf("quantile1 expect: %v get: %v", 5, q1)
	}
	if q2 := s.Percentile(.5); q2 != 7 {
		t.Errorf("quantile2 expect: %v get: %v", 7, q2)
	}
	if q3 := s.Percentile(.75); q3 != 8 {
		t.Errorf("quantile3 expect: %v get: %v", 8, q3)
	}
	if p95 := s.Percentile(.95); p95 != 9 {
		t.Errorf("percentile95 expect: %v get: %v", 9, p95)
	}
}

func TestIntMap2(t *testing.T) {
	m := NewIntMap(map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9})

	if sum := m.Sum(); sum != 285 {
		t.Errorf("Sum expect: %v get: %v", 285, sum)
	}
	if mean := m.Mean(); mean != 6.333333333333333 {
		t.Errorf("Mean expect: %v get: %v", 6.3333333333333, mean)
	}
	if median := m.Median(); median != 7 {
		t.Errorf("Median expect: %v get: %v", 7, median)
	}
	if q1 := m.Percentile(.25); q1 != 5 {
		t.Errorf("quantile1 expect: %v get: %v", 5, q1)
	}
	if q2 := m.Percentile(.5); q2 != 7 {
		t.Errorf("quantile2 expect: %v get: %v", 7, q2)
	}
	if q3 := m.Percentile(.75); q3 != 8 {
		t.Errorf("quantile3 expect: %v get: %v", 8, q3)
	}
	if p95 := m.Percentile(.95); p95 != 9 {
		t.Errorf("percentile95 expect: %v get: %v", 9, p95)
	}
}

func TestIntSlice3(t *testing.T) {
	s := NewIntSlice([]int{43, 54, 56, 61, 62, 66, 68, 69, 69, 70, 71, 72, 77, 78, 79, 85, 87, 88, 89, 93, 95, 96, 98, 99, 99})

	s.Sort()
	if sum := s.Sum(); sum != 1924 {
		t.Errorf("Sum expect: %v get: %v", 1924, sum)
	}
	if mean := s.Mean(); mean != 76.959999999999994 {
		t.Errorf("Mean expect: %v get: %v", 76.959999999999994, mean)
	}
	if median := s.Median(); median != 77 {
		t.Errorf("Median expect: %v get: %v", 77, median)
	}
	if q1 := s.Percentile(.25); q1 != 68 {
		t.Errorf("quantile1 expect: %v get: %v", 68, q1)
	}
	if p20 := s.Percentile(.20); p20 != 65.2 {
		t.Errorf("percentile20 expect: %v get: %v", 65.2, p20)
	}
	if p90 := s.Percentile(.90); p90 != 97.2 {
		t.Errorf("percentile90 expect: %v get: %v", 97.2, p90)
	}
	if p95 := s.Percentile(.95); p95 != 98.799999999999997 {
		t.Errorf("percentile95 expect: %v get: %v", 98.799999999999997, p95)
	}
}

func TestIntMap3(t *testing.T) {
	m := NewIntMap(map[int]int{43: 1, 54: 1, 56: 1, 61: 1, 62: 1, 66: 1, 68: 1, 69: 2, 70: 1, 71: 1, 72: 1, 77: 1, 78: 1, 79: 1, 85: 1, 87: 1, 88: 1, 89: 1, 93: 1, 95: 1, 96: 1, 98: 1, 99: 2})
	if sum := m.Sum(); sum != 1924 {
		t.Errorf("Sum expect: %v get: %v", 1924, sum)
	}
	if mean := m.Mean(); mean != 76.959999999999994 {
		t.Errorf("Mean expect: %v get: %v", 76.959999999999994, mean)
	}
	if median := m.Median(); median != 77 {
		t.Errorf("Median expect: %v get: %v", 77, median)
	}
	if q1 := m.Percentile(.25); q1 != 68 {
		t.Errorf("quantile1 expect: %v get: %v", 68, q1)
	}
	if p20 := m.Percentile(.20); p20 != 65.2 {
		t.Errorf("percentile20 expect: %v get: %v", 65.2, p20)
	}
	if p90 := m.Percentile(.90); p90 != 97.2 {
		t.Errorf("percentile90 expect: %v get: %v", 97.2, p90)
	}
	if p95 := m.Percentile(.95); p95 != 98.799999999999997 {
		t.Errorf("percentile95 expect: %v get: %v", 98.799999999999997, p95)
	}
}
