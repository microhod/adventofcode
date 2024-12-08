package slices

func Appends[T any](a ...[]T) []T {
	var s []T
	for i := range a {
		s = append(s, a[i]...) 
	}
	return s
}
