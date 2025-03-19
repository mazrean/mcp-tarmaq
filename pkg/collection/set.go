package collection

import "iter"

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](initial ...T) Set[T] {
	return NewSetWithCapacity[T](len(initial), initial...)
}

func NewSetWithCapacity[T comparable](capacity int, initial ...T) Set[T] {
	if capacity < len(initial) {
		capacity = len(initial)
	}

	set := make(Set[T], capacity)
	for _, v := range initial {
		set.Add(v)
	}

	return set
}

func (s Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range s {
			if !yield(value) {
				return
			}
		}
	}
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := NewSet[T]()
	for value := range s {
		if other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := NewSet[T]()
	for value := range s {
		if !other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

func (s Set[T]) Subset(other Set[T]) bool {
	for value := range s {
		if !other.Contains(value) {
			return false
		}
	}

	return true
}
