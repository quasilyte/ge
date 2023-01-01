package xslices

// Set implements a map[T]struct{} like data structure using slices.
//
// Most operations have a linear time complexity.
//
// This structure is useful for small sets of values when
// you want to re-use the memory. It also consumes less memory
// than a real map-based implementation.
type Set[T comparable] struct {
	elems []T
}

// NewSet creates a fresh empty slice-based set.
// The capacity argument is used to initialize the underlying slice.
func NewSet[T comparable](capacity int) *Set[T] {
	return &Set[T]{elems: make([]T, 0, capacity)}
}

// Elem returns the underlying slice of elements.
// The returned slice is not a copy, so you should treat
// it as a readonly slice; it's better not to retain this slice on the caller side.
func (s *Set[T]) Elems() []T { return s.elems }

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int { return len(s.elems) }

func (s *Set[T]) Reset() {
	s.elems = s.elems[:0]
}

// Contains reports whether this set contains x.
func (s *Set[T]) Contains(x T) bool {
	return Contains(s.elems, x)
}

// Add tries to insert x to a set.
// It returns true if the element was not in the set and it was added.
// If false is returned, the set already had that element.
func (s *Set[T]) Add(x T) bool {
	if Contains(s.elems, x) {
		return false
	}
	s.elems = append(s.elems, x)
	return true
}
