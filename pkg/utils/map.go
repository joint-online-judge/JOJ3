package utils

import "sort"

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func SortMap[K comparable, V any](m map[K]V, less func(i, j Pair[K, V]) bool) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair[K, V]{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return less(pairs[i], pairs[j])
	})
	return pairs
}
