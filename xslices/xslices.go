package xslices

func Index[T comparable](slice []T, x T) int {
	for i, elem := range slice {
		if x == elem {
			return i
		}
	}
	return -1
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
