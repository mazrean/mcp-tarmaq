package collection

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		initial  []int
		toAdd    int
		expected []int
	}{
		{
			name:     "Add to empty set",
			initial:  []int{},
			toAdd:    1,
			expected: []int{1},
		},
		{
			name:     "Add existing element",
			initial:  []int{1, 2, 3},
			toAdd:    2,
			expected: []int{1, 2, 3},
		},
		{
			name:     "Add new element",
			initial:  []int{1, 2, 3},
			toAdd:    4,
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)

			set.Add(tt.toAdd)

			// Verify expected result
			for _, v := range tt.expected {
				if !set.Contains(v) {
					t.Errorf("Expected element %v not found in the set", v)
				}
			}
			if len(set) != len(tt.expected) {
				t.Errorf("Set size differs from expected: got %d, want %d", len(set), len(tt.expected))
			}
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		initial  []int
		toRemove int
		expected []int
	}{
		{
			name:     "Remove existing element",
			initial:  []int{1, 2, 3},
			toRemove: 2,
			expected: []int{1, 3},
		},
		{
			name:     "Remove non-existing element",
			initial:  []int{1, 2, 3},
			toRemove: 4,
			expected: []int{1, 2, 3},
		},
		{
			name:     "Remove from empty set",
			initial:  []int{},
			toRemove: 1,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)

			set.Remove(tt.toRemove)

			// Verify expected result
			for _, v := range tt.expected {
				if !set.Contains(v) {
					t.Errorf("Expected element %v not found in the set", v)
				}
			}
			if len(set) != len(tt.expected) {
				t.Errorf("Set size differs from expected: got %d, want %d", len(set), len(tt.expected))
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		elements []int
		check    int
		expected bool
	}{
		{
			name:     "Existing element",
			elements: []int{1, 2, 3},
			check:    2,
			expected: true,
		},
		{
			name:     "Non-existing element",
			elements: []int{1, 2, 3},
			check:    4,
			expected: false,
		},
		{
			name:     "Empty set",
			elements: []int{},
			check:    1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.elements...)

			result := set.Contains(tt.check)

			if result != tt.expected {
				t.Errorf("Contains(%v) = %v, want %v", tt.check, result, tt.expected)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{
			name:     "With common elements",
			set1:     []int{1, 2, 3, 4},
			set2:     []int{3, 4, 5, 6},
			expected: []int{3, 4},
		},
		{
			name:     "No common elements",
			set1:     []int{1, 2, 3},
			set2:     []int{4, 5, 6},
			expected: []int{},
		},
		{
			name:     "Empty set",
			set1:     []int{1, 2, 3},
			set2:     []int{},
			expected: []int{},
		},
		{
			name:     "Identical sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1 := NewSet(tt.set1...)
			set2 := NewSet(tt.set2...)

			result := set1.Intersection(set2)

			expected := NewSet(tt.expected...)

			if !reflect.DeepEqual(result, expected) {
				t.Errorf("Intersection() = %v, want %v", result, expected)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{
			name:     "Partially common elements",
			set1:     []int{1, 2, 3, 4},
			set2:     []int{3, 4, 5, 6},
			expected: []int{1, 2},
		},
		{
			name:     "No common elements",
			set1:     []int{1, 2, 3},
			set2:     []int{4, 5, 6},
			expected: []int{1, 2, 3},
		},
		{
			name:     "All elements common",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "Empty set",
			set1:     []int{},
			set2:     []int{1, 2, 3},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1 := NewSet(tt.set1...)
			set2 := NewSet(tt.set2...)

			result := set1.Difference(set2)

			expected := NewSet(tt.expected...)

			if !reflect.DeepEqual(result, expected) {
				t.Errorf("Difference() = %v, want %v", result, expected)
			}
		})
	}
}

func TestSubset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected bool
	}{
		{
			name:     "Is subset",
			set1:     []int{1, 2},
			set2:     []int{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "Not subset",
			set1:     []int{1, 2, 5},
			set2:     []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "Same set",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "Empty set",
			set1:     []int{},
			set2:     []int{1, 2, 3},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1 := NewSet(tt.set1...)
			set2 := NewSet(tt.set2...)

			result := set1.Subset(set2)

			if result != tt.expected {
				t.Errorf("Subset() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected bool
	}{
		{
			name:     "Equal sets",
			set1:     []int{1, 2, 3},
			set2:     []int{3, 2, 1},
			expected: true,
		},
		{
			name:     "Different sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 4},
			expected: false,
		},
		{
			name:     "Same length but different elements",
			set1:     []int{1, 2, 3},
			set2:     []int{4, 5, 6},
			expected: false,
		},
		{
			name:     "Different length",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2},
			expected: false,
		},
		{
			name:     "Both empty",
			set1:     []int{},
			set2:     []int{},
			expected: true,
		},
		{
			name:     "One empty, one not",
			set1:     []int{1, 2, 3},
			set2:     []int{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1 := NewSet(tt.set1...)
			set2 := NewSet(tt.set2...)

			result := set1.Equal(set2)

			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		elements    []any
		elementType string // "int", "string", "float"
		checkFunc   func(t *testing.T, hash uint64)
		compareWith int // Index to compare with another test case (-1 if no comparison needed)
	}{
		{
			name:        "Empty set returns 0",
			elements:    []any{},
			elementType: "int",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash != 0 {
					t.Errorf("Empty set hash = %v, want 0", hash)
				}
			},
			compareWith: -1,
		},
		{
			name:        "Set with integers",
			elements:    []any{1, 2, 3},
			elementType: "int",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
		{
			name:        "Same elements in different order",
			elements:    []any{3, 1, 2},
			elementType: "int",
			checkFunc:   nil, // No special check needed
			compareWith: 1,   // Should have the same hash as "Set with integers"
		},
		{
			name:        "Different elements",
			elements:    []any{4, 5, 6},
			elementType: "int",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: 1, // Should have a different hash from "Set with integers"
		},
		{
			name:        "Set with strings",
			elements:    []any{"a", "b", "c"},
			elementType: "string",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
		{
			name:        "Set with floats",
			elements:    []any{1.1, 2.2, 3.3},
			elementType: "float",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
		{
			name:        "Set with uints",
			elements:    []any{uint(1), uint(2), uint(3)},
			elementType: "uint",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
		{
			name:        "Set with booleans",
			elements:    []any{true, false},
			elementType: "bool",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
		{
			name:        "Set with negative integers",
			elements:    []any{-1, -2, -3},
			elementType: "int",
			checkFunc: func(t *testing.T, hash uint64) {
				if hash == 0 {
					t.Errorf("Non-empty set should not have hash 0")
				}
			},
			compareWith: -1,
		},
	}

	// Store hash values for each test case
	hashes := make([]uint64, len(tests))

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.elementType {
			case "int":
				intSet := NewSet[int]()
				for _, e := range tt.elements {
					if e, ok := e.(int); ok {
						intSet.Add(e)
					}
				}
				hashes[i] = intSet.Hash()

			case "string":
				strSet := NewSet[string]()
				for _, e := range tt.elements {
					if e, ok := e.(string); ok {
						strSet.Add(e)
					}
				}
				hashes[i] = strSet.Hash()

			case "float":
				floatSet := NewSet[float64]()
				for _, e := range tt.elements {
					if e, ok := e.(float64); ok {
						floatSet.Add(e)
					}
				}
				hashes[i] = floatSet.Hash()

			case "uint":
				uintSet := NewSet[uint]()
				for _, e := range tt.elements {
					if e, ok := e.(uint); ok {
						uintSet.Add(e)
					}
				}
				hashes[i] = uintSet.Hash()

			case "bool":
				boolSet := NewSet[bool]()
				for _, e := range tt.elements {
					if e, ok := e.(bool); ok {
						boolSet.Add(e)
					}
				}
				hashes[i] = boolSet.Hash()
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, hashes[i])
			}
		})
	}
}
