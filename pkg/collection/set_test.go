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
