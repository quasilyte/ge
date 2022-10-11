package xslices

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

func Remove[T comparable](slice []T, x T) []T {
	i := Index(slice, x)
	if i == -1 {
		return slice
	}
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
