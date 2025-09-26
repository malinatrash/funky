package fp

// Filter filters elements of the slice by predicate
func Filter[T any](slice []T, predicate Predicate[T]) []T {
	if slice == nil {
		return nil
	}

	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// FilterWithIndex filters elements of the slice by predicate with index
func FilterWithIndex[T any](slice []T, predicate PredicateWithIndex[T]) []T {
	if slice == nil {
		return nil
	}

	var result []T
	for i, item := range slice {
		if predicate(item, i) {
			result = append(result, item)
		}
	}
	return result
}

// FilterNot filters elements that do not satisfy the predicate
func FilterNot[T any](slice []T, predicate Predicate[T]) []T {
	if slice == nil {
		return nil
	}

	var result []T
	for _, item := range slice {
		if !predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// FilterNotNil filters nil elements
func FilterNotNil[T any](slice []*T) []*T {
	if slice == nil {
		return nil
	}

	var result []*T
	for _, item := range slice {
		if item != nil {
			result = append(result, item)
		}
	}
	return result
}

// FilterNotZero filters zero values
func FilterNotZero[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}

	var zero T
	var result []T
	for _, item := range slice {
		if item != zero {
			result = append(result, item)
		}
	}
	return result
}

// Partition partitions the slice into two based on the predicate
func Partition[T any](slice []T, predicate Predicate[T]) ([]T, []T) {
	if slice == nil {
		return nil, nil
	}

	var truthy, falsy []T
	for _, item := range slice {
		if predicate(item) {
			truthy = append(truthy, item)
		} else {
			falsy = append(falsy, item)
		}
	}
	return truthy, falsy
}

// Find finds the first element that satisfies the predicate
func Find[T any](slice []T, predicate Predicate[T]) (*T, bool) {
	if slice == nil {
		return nil, false
	}

	for _, item := range slice {
		if predicate(item) {
			return &item, true
		}
	}
	return nil, false
}

// FindIndex finds the index of the first element that satisfies the predicate
func FindIndex[T any](slice []T, predicate Predicate[T]) int {
	if slice == nil {
		return -1
	}

	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// FindLast finds the last element that satisfies the predicate
func FindLast[T any](slice []T, predicate Predicate[T]) (*T, bool) {
	if slice == nil {
		return nil, false
	}

	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return &slice[i], true
		}
	}
	return nil, false
}

// All checks if all elements satisfy the predicate
func All[T any](slice []T, predicate Predicate[T]) bool {
	if slice == nil {
		return true
	}

	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Any checks if any element satisfies the predicate
func Any[T any](slice []T, predicate Predicate[T]) bool {
	if slice == nil {
		return false
	}

	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// None checks if no element satisfies the predicate
func None[T any](slice []T, predicate Predicate[T]) bool {
	return !Any(slice, predicate)
}

// Count counts the number of elements that satisfy the predicate
func Count[T any](slice []T, predicate Predicate[T]) int {
	if slice == nil {
		return 0
	}

	count := 0
	for _, item := range slice {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Unique returns unique elements of the slice
func Unique[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}

	seen := make(map[T]bool)
	var result []T

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// UniqueBy returns unique elements of the slice by key
func UniqueBy[T any, K comparable](slice []T, keyExtractor KeyExtractor[T, K]) []T {
	if slice == nil {
		return nil
	}

	seen := make(map[K]bool)
	var result []T

	for _, item := range slice {
		key := keyExtractor(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return result
}
