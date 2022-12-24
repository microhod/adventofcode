package set

type Set[T comparable] map[T]bool

func NewSet[T comparable](vals ...T) Set[T] {
	set := Set[T]{}
	for _, v := range vals {
		set.Add(v)
	}
	return set
}

func (s Set[T]) Add(vals ...T) {
	for _, v := range vals {
		s[v] = true
	}
}

func (s Set[T]) Remove(vals ...T) {
	for _, v := range vals {
		delete(s, v)
	}
}

func (s Set[T]) Contains(v T) bool {
	return s[v]
}

func (s Set[T]) ToSlice() []T {
	var slice []T
	for v := range s {
		slice = append(slice, v)
	}
	return slice
}

func Union[T comparable](sets ...Set[T]) Set[T] {
	union := Set[T]{}
	for _, s := range sets {
		for r := range s {
			union.Add(r)
		}
	}
	return union
}

func Intersect[T comparable](sets ...Set[T]) Set[T] {
	intersection := Set[T]{}
	for key := range Union(sets...) {
		if allContain(key, sets...) {
			intersection.Add(key)
		}
	}
	return intersection
}

func allContain[T comparable](key T, sets ...Set[T]) bool {
	contains := true
	for _, s := range sets {
		contains = contains && s[key]
	}
	return contains
}
