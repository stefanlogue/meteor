// Package util provides utility functions for various operations.
package util

// PrependItem prepends an item to a slice of any type
func PrependItem[T any](slice []T, item T) []T {
	items := make([]T, len(slice)+1)
	copy(items, slice)
	copy(items[1:], items)
	items[0] = item
	return items
}
