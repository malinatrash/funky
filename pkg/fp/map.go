package fp

// Map executes a function on each element of a slice and returns a new slice
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

// MapWithIndex executes a function on each element of a slice with index and returns a new slice
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

// MapParallel executes a function on each element of a slice in parallel and returns a new slice
func MapParallel[T, R any](slice []T, mapper Mapper[T, R]) []R {
	if slice == nil {
		return nil
	}

	if len(slice) < 100 { // For small slices use regular Map
		return Map(slice, mapper)
	}

	result := make([]R, len(slice))

	// Use worker pool for parallel processing
	numWorkers := min(len(slice), 10)
	chunkSize := len(slice) / numWorkers

	type job struct {
		start, end int
	}

	jobs := make(chan job, numWorkers)
	done := make(chan bool, numWorkers)

	// Start workers
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

	// Send jobs
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(slice) // Last chunk takes all remaining elements
		}
		jobs <- job{start, end}
	}
	close(jobs)

	// Wait for all workers to finish
	for w := 0; w < numWorkers; w++ {
		<-done
	}

	return result
}

// MapNotNil applies a function only to non-nil elements
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

// FlatMap applies a function and combines the results into a single slice
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

// MapKeys applies a function to the keys of a map and returns a new map
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

// MapValues applies a function to the values of a map and returns a new map
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
