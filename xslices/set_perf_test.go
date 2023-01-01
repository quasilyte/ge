package xslices

import "testing"

var sink any

func BenchmarkSetAdd4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSet[int](4)
		for j := 0; j < 4; j++ {
			s.Add(j)
		}
		sink = s
	}
}

func BenchmarkMapSetAdd5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make(map[int]struct{}, 4)
		for j := 0; j < 4; j++ {
			s[j] = struct{}{}
		}
		sink = s
	}
}

func BenchmarkSetAdd10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSet[int](10)
		for j := 0; j < 10; j++ {
			s.Add(j)
		}
		sink = s
	}
}

func BenchmarkMapSetAdd10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make(map[int]struct{}, 10)
		for j := 0; j < 10; j++ {
			s[j] = struct{}{}
		}
		sink = s
	}
}

func BenchmarkSetContains(b *testing.B) {
	s := NewSet[int](10)
	for j := 0; j < 10; j++ {
		s.Add(j)
	}

	numHits := 0

	b.Run("hit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if s.Contains(5) {
				numHits++
			}
		}
	})
	if numHits == 0 {
		b.Fail()
	}
	numHits = 0

	b.Run("miss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if s.Contains(11) {
				numHits++
			}
		}
	})
	if numHits != 0 {
		b.Fail()
	}
}

func BenchmarkMapSetContains(b *testing.B) {
	s := make(map[int]struct{}, 10)
	for j := 0; j < 10; j++ {
		s[j] = struct{}{}
	}

	numHits := 0

	b.Run("hit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := s[5]; ok {
				numHits++
			}
		}
	})
	if numHits == 0 {
		b.Fail()
	}
	numHits = 0

	b.Run("miss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := s[11]; ok {
				numHits++
			}
		}
	})
	if numHits != 0 {
		b.Fail()
	}
}
