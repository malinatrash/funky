package fp

// Reduce reduces a slice to a single value
func Reduce[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	if slice == nil {
		return initial
	}

	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// ReduceWithIndex reduces a slice to a single value with index
func ReduceWithIndex[T, R any](slice []T, reducer ReducerWithIndex[T, R], initial R) R {
	if slice == nil {
		return initial
	}

	result := initial
	for i, item := range slice {
		result = reducer(result, item, i)
	}
	return result
}

// ReduceRight reduces a slice to a single value from right to left
func ReduceRight[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	if slice == nil {
		return initial
	}

	result := initial
	for i := len(slice) - 1; i >= 0; i-- {
		result = reducer(result, slice[i])
	}
	return result
}

// Fold - alias for Reduce (more functional name)
func Fold[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return Reduce(slice, reducer, initial)
}

// FoldLeft - alias for Reduce
func FoldLeft[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return Reduce(slice, reducer, initial)
}

// FoldRight - alias for ReduceRight
func FoldRight[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return ReduceRight(slice, reducer, initial)
}

// Sum sums numeric values
func Sum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	var sum T
	for _, item := range slice {
		sum += item
	}
	return sum
}

// Product multiplies numeric values
func Product[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	if len(slice) == 0 {
		return 0
	}

	var product T = 1
	for _, item := range slice {
		product *= item
	}
	return product
}

// Min finds the minimum value
func Min[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	min := slice[0]
	for _, item := range slice[1:] {
		if item < min {
			min = item
		}
	}
	return min, true
}

// Max finds the maximum value
func Max[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	max := slice[0]
	for _, item := range slice[1:] {
		if item > max {
			max = item
		}
	}
	return max, true
}

// MinBy finds the minimum element by comparator
func MinBy[T any](slice []T, comparator Comparator[T]) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	min := slice[0]
	for _, item := range slice[1:] {
		if comparator(item, min) < 0 {
			min = item
		}
	}
	return min, true
}

// MaxBy finds the maximum element by comparator
func MaxBy[T any](slice []T, comparator Comparator[T]) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	max := slice[0]
	for _, item := range slice[1:] {
		if comparator(item, max) > 0 {
			max = item
		}
	}
	return max, true
}

// GroupBy groups elements by key
func GroupBy[T any, K comparable](slice []T, keyExtractor KeyExtractor[T, K]) map[K][]T {
	if slice == nil {
		return nil
	}

	groups := make(map[K][]T)
	for _, item := range slice {
		key := keyExtractor(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// GroupByCount groups and counts elements by key
func GroupByCount[T any, K comparable](slice []T, keyExtractor KeyExtractor[T, K]) map[K]int {
	if slice == nil {
		return nil
	}

	counts := make(map[K]int)
	for _, item := range slice {
		key := keyExtractor(item)
		counts[key]++
	}
	return counts
}

// Join joins strings with a separator
func Join(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}

	totalLen := len(separator) * (len(slice) - 1)
	for _, s := range slice {
		totalLen += len(s)
	}

	result := make([]byte, 0, totalLen)
	result = append(result, slice[0]...)

	for _, s := range slice[1:] {
		result = append(result, separator...)
		result = append(result, s...)
	}

	return string(result)
}

// JoinBy joins elements into a string using a mapper function
func JoinBy[T any](slice []T, mapper func(T) string, separator string) string {
	if len(slice) == 0 {
		return ""
	}

	strings := make([]string, len(slice))
	for i, item := range slice {
		strings[i] = mapper(item)
	}

	return Join(strings, separator)
}

// Aggregate executes aggregation with multiple accumulators
func Aggregate[T any, R any](slice []T, initial R, aggregator func(R, T) R) R {
	return Reduce(slice, aggregator, initial)
}
