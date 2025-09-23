package util

import (
	"reflect"
	"testing"
)

type item interface {
	string | int
}
type testcase[T item] struct {
	desc  string
	item  T
	slice []T
	want  []T
}

func TestPrependItem(t *testing.T) {
	intCases := []testcase[int]{
		{
			desc:  "prepend to non-empty slice",
			item:  1,
			slice: []int{2, 3, 4},
			want:  []int{1, 2, 3, 4},
		},
		{
			desc:  "prepend to empty slice",
			item:  5,
			slice: []int{},
			want:  []int{5},
		},
	}
	stringCases := []testcase[string]{
		{
			desc:  "prepend to non-empty slice",
			item:  "a",
			slice: []string{"b", "c", "d"},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			desc:  "prepend to empty slice",
			item:  "x",
			slice: []string{},
			want:  []string{"x"},
		},
	}
	for _, tc := range intCases {
		t.Run(tc.desc, runTestCase(tc))
	}
	for _, tc := range stringCases {
		t.Run(tc.desc, runTestCase(tc))
	}
}

func runTestCase[T item](tc testcase[T]) func(t *testing.T) {
	return func(t *testing.T) {
		got := PrependItem(tc.slice, tc.item)
		assertEqual(t, tc.want, got)
	}
}

func assertEqual[T item](t *testing.T, want, got []T) {
	t.Helper()
	if len(want) != len(got) {
		t.Error("length mismatch")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
