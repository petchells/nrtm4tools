package util

// HashMap is a generic definition of a map
type HashMap[K comparable, E any] map[K]E

// KeyFunction is a function that extracts a map key from an entity
type KeyFunction[K comparable, E any] func(E) K

// SliceToMap creates a map from a slice or array. fn should return the key for the map
func SliceToMap[K comparable, E any](fn KeyFunction[K, E], vals []E) HashMap[K, E] {
	m := HashMap[K, E]{}
	for _, v := range vals {
		m[fn(v)] = v
	}
	return m
}
