package fp

// Reduce сворачивает слайс в одно значение
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

// ReduceWithIndex сворачивает слайс в одно значение с индексом
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

// ReduceRight сворачивает слайс справа налево
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

// Fold - алиас для Reduce (более функциональное название)
func Fold[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return Reduce(slice, reducer, initial)
}

// FoldLeft - алиас для Reduce
func FoldLeft[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return Reduce(slice, reducer, initial)
}

// FoldRight - алиас для ReduceRight
func FoldRight[T, R any](slice []T, reducer Reducer[T, R], initial R) R {
	return ReduceRight(slice, reducer, initial)
}

// Sum суммирует числовые значения
func Sum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	var sum T
	for _, item := range slice {
		sum += item
	}
	return sum
}

// Product перемножает числовые значения
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

// Min находит минимальное значение
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

// Max находит максимальное значение
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

// MinBy находит минимальный элемент по функции сравнения
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

// MaxBy находит максимальный элемент по функции сравнения
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

// GroupBy группирует элементы по ключу
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

// GroupByCount группирует и подсчитывает элементы по ключу
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

// Join объединяет строки с разделителем
func Join(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	
	// Подсчитываем общую длину для оптимизации
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

// JoinBy объединяет элементы в строку через функцию преобразования
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

// Aggregate выполняет агрегацию с несколькими аккумуляторами
func Aggregate[T any, R any](slice []T, initial R, aggregator func(R, T) R) R {
	return Reduce(slice, aggregator, initial)
}