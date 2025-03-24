package collection

import (
	"hash/fnv"
	"iter"
	"reflect"
	"sort"
	"strconv"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](initial ...T) Set[T] {
	return NewSetWithCapacity(len(initial), initial...)
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

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}

	for value := range s {
		if !other.Contains(value) {
			return false
		}
	}

	return true
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

// Hash returns a hash value for the set with low probability of collisions
func (s Set[T]) Hash() uint64 {
	if len(s) == 0 {
		return 0
	}

	// We need to sort keys for stable hash results,
	// but since we can't sort arbitrary types, we'll calculate individual hashes and combine them
	h := fnv.New64a()

	// Calculate hash for each element
	hashes := make([]uint64, 0, len(s))
	for value := range s {
		// Hash calculation based on element type
		var elemHash uint64

		rv := reflect.ValueOf(value)
		//nolint:exhaustive // We only need to handle basic types
		switch rv.Kind() {
		case reflect.Bool:
			if rv.Bool() {
				elemHash = 1
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			elemHash = rv.Uint()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			elemHash = rv.Uint()
		case reflect.Float32, reflect.Float64:
			elemHash = uint64(rv.Float())
		case reflect.String:
			h.Reset()
			h.Write([]byte(rv.String()))
			elemHash = h.Sum64()
		default:
			h.Reset()
			h.Write([]byte(strconv.Itoa(int(reflect.ValueOf(value).Pointer()))))
			elemHash = h.Sum64()
		}

		hashes = append(hashes, elemHash)
	}

	// Sort to eliminate differences due to hash order
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})

	// Calculate final hash using FNV-1a
	h.Reset()
	for _, hash := range hashes {
		h.Write([]byte{
			byte(hash),
			byte(hash >> 8),
			byte(hash >> 16),
			byte(hash >> 24),
			byte(hash >> 32),
			byte(hash >> 40),
			byte(hash >> 48),
			byte(hash >> 56),
		})
	}

	return h.Sum64()
}
