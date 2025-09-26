package fp

// Map применяет функцию к каждому элементу слайса и возвращает новый слайс
func Map[T, R any](slice []T, mapper Mapper[T, R]) []R {
	if slice == nil {
		return nil
	}
	
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// MapWithIndex применяет функцию к каждому элементу слайса с индексом
func MapWithIndex[T, R any](slice []T, mapper MapperWithIndex[T, R]) []R {
	if slice == nil {
		return nil
	}
	
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item, i)
	}
	return result
}

// MapParallel применяет функцию параллельно (для больших слайсов)
func MapParallel[T, R any](slice []T, mapper Mapper[T, R]) []R {
	if slice == nil {
		return nil
	}
	
	if len(slice) < 100 { // Для маленьких слайсов используем обычный Map
		return Map(slice, mapper)
	}
	
	result := make([]R, len(slice))
	
	// Используем воркер пул для параллельной обработки
	numWorkers := min(len(slice), 10)
	chunkSize := len(slice) / numWorkers
	
	type job struct {
		start, end int
	}
	
	jobs := make(chan job, numWorkers)
	done := make(chan bool, numWorkers)
	
	// Запускаем воркеры
	for w := 0; w < numWorkers; w++ {
		go func() {
			for j := range jobs {
				for i := j.start; i < j.end; i++ {
					result[i] = mapper(slice[i])
				}
			}
			done <- true
		}()
	}
	
	// Отправляем задачи
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(slice) // Последний чанк берет все оставшиеся элементы
		}
		jobs <- job{start, end}
	}
	close(jobs)
	
	// Ждем завершения всех воркеров
	for w := 0; w < numWorkers; w++ {
		<-done
	}
	
	return result
}

// MapNotNil применяет функцию только к не-nil элементам
func MapNotNil[T, R any](slice []*T, mapper Mapper[T, R]) []R {
	if slice == nil {
		return nil
	}
	
	var result []R
	for _, item := range slice {
		if item != nil {
			result = append(result, mapper(*item))
		}
	}
	return result
}

// FlatMap применяет функцию и объединяет результаты в один слайс
func FlatMap[T, R any](slice []T, mapper func(T) []R) []R {
	if slice == nil {
		return nil
	}
	
	var result []R
	for _, item := range slice {
		result = append(result, mapper(item)...)
	}
	return result
}

// MapKeys применяет функцию к ключам мапы
func MapKeys[K1 comparable, K2 comparable, V any](m map[K1]V, mapper func(K1) K2) map[K2]V {
	if m == nil {
		return nil
	}
	
	result := make(map[K2]V, len(m))
	for k, v := range m {
		result[mapper(k)] = v
	}
	return result
}

// MapValues применяет функцию к значениям мапы
func MapValues[K comparable, V1, V2 any](m map[K]V1, mapper func(V1) V2) map[K]V2 {
	if m == nil {
		return nil
	}
	
	result := make(map[K]V2, len(m))
	for k, v := range m {
		result[k] = mapper(v)
	}
	return result
}

// min helper function для Go версий < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}