package fp

import (
	"context"
	"sync"
)

// Stream представляет поток данных для ленивых вычислений
type Stream[T any] struct {
	source   func() <-chan T
	pipeline []func(<-chan T) <-chan T
}

// NewStream создает новый поток из слайса
func NewStream[T any](slice []T) *Stream[T] {
	return &Stream[T]{
		source: func() <-chan T {
			ch := make(chan T, len(slice))
			go func() {
				defer close(ch)
				for _, item := range slice {
					ch <- item
				}
			}()
			return ch
		},
		pipeline: []func(<-chan T) <-chan T{},
	}
}

// NewStreamFromChannel создает поток из канала
func NewStreamFromChannel[T any](ch <-chan T) *Stream[T] {
	return &Stream[T]{
		source:   func() <-chan T { return ch },
		pipeline: []func(<-chan T) <-chan T{},
	}
}

// NewStreamFromFunc создает поток из функции-генератора
func NewStreamFromFunc[T any](generator func() <-chan T) *Stream[T] {
	return &Stream[T]{
		source:   generator,
		pipeline: []func(<-chan T) <-chan T{},
	}
}

// Map применяет функцию преобразования к потоку
func (s *Stream[T]) Map(mapper Mapper[T, T]) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			for item := range input {
				output <- mapper(item)
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Filter применяет фильтрацию к потоку
func (s *Stream[T]) Filter(predicate Predicate[T]) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			for item := range input {
				if predicate(item) {
					output <- item
				}
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Take берет первые n элементов из потока
func (s *Stream[T]) Take(n int) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			count := 0
			for item := range input {
				if count >= n {
					break
				}
				output <- item
				count++
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Skip пропускает первые n элементов потока
func (s *Stream[T]) Skip(n int) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			count := 0
			for item := range input {
				if count >= n {
					output <- item
				}
				count++
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Distinct удаляет дубликаты из потока
func (s *Stream[T]) Distinct(equals Equality[T]) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			var seen []T
			for item := range input {
				isDuplicate := false
				for _, seenItem := range seen {
					if equals(item, seenItem) {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					seen = append(seen, item)
					output <- item
				}
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// DistinctComparable удаляет дубликаты для comparable типов
func (s *Stream[T]) DistinctComparable() *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			seen := make(map[interface{}]bool)
			for item := range input {
				if !seen[item] {
					seen[item] = true
					output <- item
				}
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Parallel применяет параллельную обработку к потоку
func (s *Stream[T]) Parallel(workerCount int, processor func(T) T) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		
		go func() {
			defer close(output)
			
			var wg sync.WaitGroup
			jobs := make(chan T, workerCount*2)
			results := make(chan T, workerCount*2)
			
			// Запускаем воркеры
			for i := 0; i < workerCount; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for job := range jobs {
						results <- processor(job)
					}
				}()
			}
			
			// Отправляем задачи
			go func() {
				defer close(jobs)
				for item := range input {
					jobs <- item
				}
			}()
			
			// Собираем результаты
			go func() {
				wg.Wait()
				close(results)
			}()
			
			for result := range results {
				output <- result
			}
		}()
		
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Buffer буферизует поток
func (s *Stream[T]) Buffer(size int) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T, size)
		go func() {
			defer close(output)
			for item := range input {
				output <- item
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// WithContext добавляет контекст к потоку
func (s *Stream[T]) WithContext(ctx context.Context) *Stream[T] {
	newPipeline := append(s.pipeline, func(input <-chan T) <-chan T {
		output := make(chan T)
		go func() {
			defer close(output)
			for {
				select {
				case item, ok := <-input:
					if !ok {
						return
					}
					select {
					case output <- item:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		return output
	})
	
	return &Stream[T]{
		source:   s.source,
		pipeline: newPipeline,
	}
}

// Collect собирает все элементы потока в слайс
func (s *Stream[T]) Collect() []T {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	var result []T
	for item := range ch {
		result = append(result, item)
	}
	
	return result
}

// CollectToChannel собирает поток в канал
func (s *Stream[T]) CollectToChannel() <-chan T {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	return ch
}

// ForEach выполняет функцию для каждого элемента потока
func (s *Stream[T]) ForEach(action func(T)) {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	for item := range ch {
		action(item)
	}
}

// Reduce сворачивает поток в одно значение
func (s *Stream[T]) Reduce(reducer Reducer[T, T], initial T) T {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	result := initial
	for item := range ch {
		result = reducer(result, item)
	}
	
	return result
}

// Count подсчитывает количество элементов в потоке
func (s *Stream[T]) Count() int {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	count := 0
	for range ch {
		count++
	}
	
	return count
}

// AnyMatch проверяет, есть ли элемент, удовлетворяющий предикату
func (s *Stream[T]) AnyMatch(predicate Predicate[T]) bool {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	for item := range ch {
		if predicate(item) {
			return true
		}
	}
	
	return false
}

// AllMatch проверяет, что все элементы удовлетворяют предикату
func (s *Stream[T]) AllMatch(predicate Predicate[T]) bool {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	for item := range ch {
		if !predicate(item) {
			return false
		}
	}
	
	return true
}

// FindFirst находит первый элемент, удовлетворяющий предикату
func (s *Stream[T]) FindFirst(predicate Predicate[T]) Optional[T] {
	ch := s.source()
	
	// Применяем все этапы pipeline
	for _, stage := range s.pipeline {
		ch = stage(ch)
	}
	
	for item := range ch {
		if predicate(item) {
			return Some(item)
		}
	}
	
	return Empty[T]()
}

// StreamBuilder помогает создавать потоки
type StreamBuilder[T any] struct {
	items []T
}

// NewStreamBuilder создает новый builder
func NewStreamBuilder[T any]() *StreamBuilder[T] {
	return &StreamBuilder[T]{items: []T{}}
}

// Add добавляет элемент
func (sb *StreamBuilder[T]) Add(item T) *StreamBuilder[T] {
	sb.items = append(sb.items, item)
	return sb
}

// AddAll добавляет все элементы из слайса
func (sb *StreamBuilder[T]) AddAll(items []T) *StreamBuilder[T] {
	sb.items = append(sb.items, items...)
	return sb
}

// Build создает поток
func (sb *StreamBuilder[T]) Build() *Stream[T] {
	return NewStream(sb.items)
}

// InfiniteStream создает бесконечный поток
func InfiniteStream[T any](generator func() T) *Stream[T] {
	return &Stream[T]{
		source: func() <-chan T {
			ch := make(chan T)
			go func() {
				for {
					ch <- generator()
				}
			}()
			return ch
		},
		pipeline: []func(<-chan T) <-chan T{},
	}
}

// RangeStream создает поток чисел от start до end
func RangeStream(start, end int) *Stream[int] {
	return &Stream[int]{
		source: func() <-chan int {
			ch := make(chan int)
			go func() {
				defer close(ch)
				for i := start; i < end; i++ {
					ch <- i
				}
			}()
			return ch
		},
		pipeline: []func(<-chan int) <-chan int{},
	}
}

// RepeatStream создает поток, повторяющий значение n раз
func RepeatStream[T any](value T, count int) *Stream[T] {
	return &Stream[T]{
		source: func() <-chan T {
			ch := make(chan T)
			go func() {
				defer close(ch)
				for i := 0; i < count; i++ {
					ch <- value
				}
			}()
			return ch
		},
		pipeline: []func(<-chan T) <-chan T{},
	}
}
