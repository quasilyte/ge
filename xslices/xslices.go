package xslices

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Sort[T constraints.Ordered](slice []T) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
}

func Diff[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return s2
	}
	if len(s2) == 0 {
		return s1
	}
	smaller := s1
	bigger := s2
	if len(s2) < len(s1) {
		smaller = s2
		bigger = s1
	}
	var result []T
	if len(smaller) <= 4 {
		for _, v := range bigger {
			if !Contains(smaller, v) {
				result = append(result, v)
			}
		}
	} else {
		set := make(map[T]struct{}, len(smaller))
		for _, v := range smaller {
			set[v] = struct{}{}
		}
		for _, v := range bigger {
			if _, ok := set[v]; !ok {
				result = append(result, v)
			}
		}
	}
	return result
}

func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func ContainsWhere[T comparable](slice []T, pred func(T) bool) bool {
	return IndexWhere(slice, pred) != -1
}

func Contains[T comparable](slice []T, x T) bool {
	return Index(slice, x) != -1
}

func Index[T comparable](slice []T, x T) int {
	for i, elem := range slice {
		if x == elem {
			return i
		}
	}
	return -1
}

func IndexWhere[T any](slice []T, pred func(T) bool) int {
	for i, elem := range slice {
		if pred(elem) {
			return i
		}
	}
	return -1
}

func Find[T any](slice []T, pred func(*T) bool) *T {
	for i := range slice {
		if pred(&slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

func CountIf[T any](slice []T, pred func(T) bool) int {
	c := 0
	for _, elem := range slice {
		if pred(elem) {
			c++
		}
	}
	return c
}

func Remove[T comparable](slice []T, x T) []T {
	i := Index(slice, x)
	if i == -1 {
		return slice
	}
	return RemoveAt(slice, i)
}

func RemoveAt[T any](slice []T, i int) []T {
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice
}

func RemoveIf[T any](slice []T, pred func(T) bool) []T {
	keep := slice[:0]
	for _, elem := range slice {
		if !pred(elem) {
			keep = append(keep, elem)
		}
	}
	return keep
}

func Any[T any](slice []T, pred func(T) bool) bool {
	for _, elem := range slice {
		if pred(elem) {
			return true
		}
	}
	return false
}

func All[T any](slice []T, pred func(T) bool) bool {
	for _, elem := range slice {
		if !pred(elem) {
			return false
		}
	}
	return true
}

func Prepend[T any](slice []T, elems ...T) []T {
	result := make([]T, len(elems), len(slice)+len(elems))
	copy(result, elems)
	return append(result, slice...)
}
