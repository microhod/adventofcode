package copy

// Map creates a shallow copy
func Map[K comparable, V any](m map[K]V) map[K]V {
	copy := map[K]V{}
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

// Slice creates a shallow copy
func Slice[V any](s []V) []V {
	return append([]V{}, s...)
}
