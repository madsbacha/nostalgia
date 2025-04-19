package util

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      []int
		str    int
		expect bool
	}{
		{"ElementExists", []int{1, 2, 3}, 2, true},
		{"ElementDoesNotExist", []int{1, 2, 3}, 4, false},
		{"EmptySlice", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.s, tt.str)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name   string
		s      []int
		str    int
		expect []int
	}{
		{"RemoveExistingElement", []int{1, 2, 3}, 2, []int{1, 3}},
		{"RemoveNonExistingElement", []int{1, 2, 3}, 4, []int{1, 2, 3}},
		{"EmptySlice", []int{}, 1, []int{}},
		{"RemoveAllElements", []int{1, 1, 1}, 1, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Remove(tt.s, tt.str)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name   string
		a      []int
		b      []int
		expect []int
	}{
		{"NonEmptyIntersection", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
		{"NoIntersection", []int{1, 2, 3}, []int{4, 5, 6}, []int{}},
		{"EmptyFirstSlice", []int{}, []int{1, 2, 3}, []int{}},
		{"EmptySecondSlice", []int{1, 2, 3}, []int{}, []int{}},
		{"BothEmpty", []int{}, []int{}, []int{}},
		{"SingleElementIntersection", []int{1}, []int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersection(tt.a, tt.b)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name   string
		a      []int
		expect []int
	}{
		{"NoDuplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"WithDuplicates", []int{1, 2, 2, 3, 3}, []int{1, 2, 3}},
		{"EmptySlice", []int{}, []int{}},
		{"AllDuplicates", []int{1, 1, 1}, []int{1}},
		{"SingleElement", []int{1}, []int{1}},
		{"MultipleDuplicates", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.a)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name   string
		a      []int
		b      []int
		expect []int
	}{
		{"UnionWithOverlap", []int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"UnionWithoutOverlap", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"EmptyFirstSlice", []int{}, []int{1, 2}, []int{1, 2}},
		{"EmptySecondSlice", []int{1, 2}, []int{}, []int{1, 2}},
		{"BothEmpty", []int{}, []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Union(tt.a, tt.b)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestExcept(t *testing.T) {
	tests := []struct {
		name   string
		list   []int
		except []int
		expect []int
	}{
		{"ExceptSubset", []int{1, 2, 3}, []int{2}, []int{1, 3}},
		{"ExceptNonExisting", []int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3}},
		{"EmptyList", []int{}, []int{1, 2}, []int{}},
		{"EmptyExcept", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"AllElementsExcepted", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
		{"BothEmpty", []int{}, []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Except(tt.list, tt.except)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}
