package xmaps

import (
	"sort"

	"github.com/quasilyte/ge/tuple"
	"golang.org/x/exp/constraints"
)

func KeysAndValues[Key comparable, Value any](m map[Key]Value) []tuple.Value2[Key, Value] {
	pairs := make([]tuple.Value2[Key, Value], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, tuple.New2(k, v))
	}
	return pairs
}

func Keys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func KeysSortedByValue[Key comparable, Value any](m map[Key]Value, less func(v1, v2 Value) bool) []Key {
	keysAndValues := KeysAndValues(m)
	sort.Slice(keysAndValues, func(i, j int) bool {
		return less(keysAndValues[i].Second, keysAndValues[j].Second)
	})
	keys := make([]Key, len(keysAndValues))
	for i := range keys {
		keys[i] = keysAndValues[i].First
	}
	return keys
}

func KeysSorted[Key constraints.Ordered, Value any](m map[Key]Value) []Key {
	keys := Keys(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
