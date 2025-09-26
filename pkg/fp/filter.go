package fp

// Filter фильтрует элементы слайса по предикату
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

// FilterWithIndex фильтрует элементы слайса по предикату с индексом
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

// FilterNot фильтрует элементы, НЕ удовлетворяющие предикату
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

// FilterNotNil фильтрует nil элементы
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

// FilterNotZero фильтрует нулевые значения
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

// Partition разделяет слайс на два по предикату
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

// Find находит первый элемент, удовлетворяющий предикату
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

// FindIndex находит индекс первого элемента, удовлетворяющего предикату
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

// FindLast находит последний элемент, удовлетворяющий предикату
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

// All проверяет, что все элементы удовлетворяют предикату
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

// Any проверяет, что хотя бы один элемент удовлетворяет предикату
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

// None проверяет, что ни один элемент не удовлетворяет предикату
func None[T any](slice []T, predicate Predicate[T]) bool {
	return !Any(slice, predicate)
}

// Count подсчитывает количество элементов, удовлетворяющих предикату
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

// Unique возвращает уникальные элементы слайса
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

// UniqueBy возвращает уникальные элементы по ключу
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