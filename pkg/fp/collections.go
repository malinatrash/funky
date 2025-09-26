package fp

// Chunk splits slice into chunks of a given size
func Chunk[T any](slice []T, size int) [][]T {
	if slice == nil || size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// ChunkBy splits slice into chunks based on a predicate
func ChunkBy[T any](slice []T, predicate Predicate[T]) [][]T {
	if slice == nil {
		return nil
	}

	var chunks [][]T
	var currentChunk []T

	for _, item := range slice {
		if predicate(item) && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = []T{item}
		} else {
			currentChunk = append(currentChunk, item)
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// Sliding creates a sliding window of a given size
func Sliding[T any](slice []T, size int) [][]T {
	if slice == nil || size <= 0 || size > len(slice) {
		return nil
	}

	var windows [][]T
	for i := 0; i <= len(slice)-size; i++ {
		window := make([]T, size)
		copy(window, slice[i:i+size])
		windows = append(windows, window)
	}
	return windows
}

// Take takes the first n elements
func Take[T any](slice []T, n int) []T {
	if slice == nil || n <= 0 {
		return nil
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}

	result := make([]T, n)
	copy(result, slice[:n])
	return result
}

// TakeWhile takes elements while the predicate is true
func TakeWhile[T any](slice []T, predicate Predicate[T]) []T {
	if slice == nil {
		return nil
	}

	var result []T
	for _, item := range slice {
		if !predicate(item) {
			break
		}
		result = append(result, item)
	}
	return result
}

// Drop drops the first n elements
func Drop[T any](slice []T, n int) []T {
	if slice == nil || n <= 0 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	if n >= len(slice) {
		return []T{}
	}

	result := make([]T, len(slice)-n)
	copy(result, slice[n:])
	return result
}

// DropWhile drops elements while the predicate is true
func DropWhile[T any](slice []T, predicate Predicate[T]) []T {
	if slice == nil {
		return nil
	}

	start := 0
	for i, item := range slice {
		if !predicate(item) {
			start = i
			break
		}
		if i == len(slice)-1 {
			return []T{}
		}
	}

	result := make([]T, len(slice)-start)
	copy(result, slice[start:])
	return result
}

// Reverse reverses the slice
func Reverse[T any](slice []T) []T {
	if slice == nil {
		return nil
	}

	result := make([]T, len(slice))
	for i, item := range slice {
		result[len(slice)-1-i] = item
	}
	return result
}

// Zip zips two slices into a slice of pairs
func Zip[T, R any](slice1 []T, slice2 []R) []Pair[T, R] {
	if slice1 == nil || slice2 == nil {
		return nil
	}

	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]Pair[T, R], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = Pair[T, R]{First: slice1[i], Second: slice2[i]}
	}
	return result
}

// ZipWith zips two slices using a function
func ZipWith[T, R, S any](slice1 []T, slice2 []R, zipper func(T, R) S) []S {
	if slice1 == nil || slice2 == nil {
		return nil
	}

	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]S, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = zipper(slice1[i], slice2[i])
	}
	return result
}

// Unzip splits a slice of pairs into two slices
func Unzip[T, R any](pairs []Pair[T, R]) ([]T, []R) {
	if pairs == nil {
		return nil, nil
	}

	first := make([]T, len(pairs))
	second := make([]R, len(pairs))

	for i, pair := range pairs {
		first[i] = pair.First
		second[i] = pair.Second
	}

	return first, second
}

// Flatten flattens a two-dimensional slice into a one-dimensional slice
func Flatten[T any](slices [][]T) []T {
	if slices == nil {
		return nil
	}

	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Intersperse inserts an element between all elements of the slice
func Intersperse[T any](slice []T, separator T) []T {
	if slice == nil || len(slice) <= 1 {
		return slice
	}

	result := make([]T, 0, len(slice)*2-1)
	result = append(result, slice[0])

	for _, item := range slice[1:] {
		result = append(result, separator, item)
	}

	return result
}

// Transpose transposes a two-dimensional slice
func Transpose[T any](matrix [][]T) [][]T {
	if matrix == nil || len(matrix) == 0 {
		return nil
	}

	maxCols := 0
	for _, row := range matrix {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	if maxCols == 0 {
		return nil
	}

	result := make([][]T, maxCols)
	for i := range result {
		result[i] = make([]T, 0, len(matrix))
	}

	for _, row := range matrix {
		for j, item := range row {
			result[j] = append(result[j], item)
		}
	}

	return result
}

// Cartesian returns the Cartesian product of two slices
func Cartesian[T, R any](slice1 []T, slice2 []R) []Pair[T, R] {
	if slice1 == nil || slice2 == nil {
		return nil
	}

	result := make([]Pair[T, R], 0, len(slice1)*len(slice2))
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			result = append(result, Pair[T, R]{First: item1, Second: item2})
		}
	}
	return result
}

// Pair represents a pair of values
type Pair[T, R any] struct {
	First  T
	Second R
}

// Permutations generates all permutations of a slice
func Permutations[T any](slice []T) [][]T {
	if slice == nil || len(slice) == 0 {
		return [][]T{{}}
	}

	if len(slice) == 1 {
		return [][]T{{slice[0]}}
	}

	var result [][]T
	for i, item := range slice {
		// Create a slice without the current element
		remaining := make([]T, 0, len(slice)-1)
		remaining = append(remaining, slice[:i]...)
		remaining = append(remaining, slice[i+1:]...)

		// Generate permutations for the remaining elements
		subPerms := Permutations(remaining)

		// Add the current element to the beginning of each permutation
		for _, perm := range subPerms {
			newPerm := make([]T, 0, len(slice))
			newPerm = append(newPerm, item)
			newPerm = append(newPerm, perm...)
			result = append(result, newPerm)
		}
	}

	return result
}

// Combinations generates all combinations of a given size
func Combinations[T any](slice []T, r int) [][]T {
	if slice == nil || r < 0 || r > len(slice) {
		return nil
	}

	if r == 0 {
		return [][]T{{}}
	}

	if r == len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return [][]T{result}
	}

	var result [][]T

	// Recursively generate combinations
	var generate func(start, remaining int, current []T)
	generate = func(start, remaining int, current []T) {
		if remaining == 0 {
			combination := make([]T, len(current))
			copy(combination, current)
			result = append(result, combination)
			return
		}

		for i := start; i <= len(slice)-remaining; i++ {
			current = append(current, slice[i])
			generate(i+1, remaining-1, current)
			current = current[:len(current)-1]
		}
	}

	generate(0, r, []T{})
	return result
}
