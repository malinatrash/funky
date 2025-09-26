package fp

// Chunk разбивает слайс на чанки заданного размера
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

// ChunkBy разбивает слайс на чанки по предикату
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

// Sliding создает скользящее окно заданного размера
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

// Take берет первые n элементов
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

// TakeWhile берет элементы пока предикат истинен
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

// Drop пропускает первые n элементов
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

// DropWhile пропускает элементы пока предикат истинен
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

// Reverse переворачивает слайс
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

// Zip объединяет два слайса в слайс пар
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

// ZipWith объединяет два слайса с помощью функции
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

// Unzip разделяет слайс пар на два слайса
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

// Flatten сплющивает двумерный слайс в одномерный
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

// Intersperse вставляет элемент между всеми элементами слайса
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

// Transpose транспонирует двумерный слайс
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

// Cartesian возвращает декартово произведение двух слайсов
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

// Pair представляет пару значений
type Pair[T, R any] struct {
	First  T
	Second R
}

// Permutations генерирует все перестановки слайса
func Permutations[T any](slice []T) [][]T {
	if slice == nil || len(slice) == 0 {
		return [][]T{{}}
	}
	
	if len(slice) == 1 {
		return [][]T{{slice[0]}}
	}
	
	var result [][]T
	for i, item := range slice {
		// Создаем слайс без текущего элемента
		remaining := make([]T, 0, len(slice)-1)
		remaining = append(remaining, slice[:i]...)
		remaining = append(remaining, slice[i+1:]...)
		
		// Генерируем перестановки для оставшихся элементов
		subPerms := Permutations(remaining)
		
		// Добавляем текущий элемент в начало каждой перестановки
		for _, perm := range subPerms {
			newPerm := make([]T, 0, len(slice))
			newPerm = append(newPerm, item)
			newPerm = append(newPerm, perm...)
			result = append(result, newPerm)
		}
	}
	
	return result
}

// Combinations генерирует все сочетания заданного размера
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
	
	// Рекурсивно генерируем сочетания
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
