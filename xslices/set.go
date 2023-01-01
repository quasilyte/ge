package xslices

// Set implements a map[T]struct{} like data structure using slices.
//
// Most operations have a linear time complexity.
//
// This structure is useful for small sets of values when
// you want to re-use the memory. It also consumes less memory
// than a real map-based implementation.
//
// To get the set length, use `len(*s)`.
// To iterate over the set, use `for _, x := range *s`.
type Set[T comparable] []T

// NewSet creates a fresh empty slice-based set.
// The capacity argument is used to initialize the underlying slice.
func NewSet[T comparable](capacity int) *Set[T] {
	set := Set[T](make([]T, 0, capacity))
	return &set
}

func (s *Set[T]) Reset() {
	(*s) = (*s)[:0]
}

// Contains reports whether this set contains x.
func (s *Set[T]) Contains(x T) bool {
	return Contains(*s, x)
}

// Add tries to insert x to a set.
// It returns true if the element was not in the set and it was added.
// If false is returned, the set already had that element.
func (s *Set[T]) Add(x T) bool {
	if Contains(*s, x) {
		return false
	}
	*s = append(*s, x)
	return true
}
