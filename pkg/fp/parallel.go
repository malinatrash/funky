package fp

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// ParallelConfig configuration for parallel processing
type ParallelConfig struct {
	WorkerCount int
	BufferSize  int
}

// DefaultParallelConfig returns the default configuration
func DefaultParallelConfig() ParallelConfig {
	return ParallelConfig{
		WorkerCount: runtime.NumCPU(),
		BufferSize:  100,
	}
}

// MapParallelWithConfig parallel Map with configuration
func MapParallelWithConfig[T, R any](slice []T, mapper Mapper[T, R], config ParallelConfig) []R {
	if slice == nil || len(slice) == 0 {
		return nil
	}

	if len(slice) < config.WorkerCount {
		return Map(slice, mapper)
	}

	result := make([]R, len(slice))
	jobs := make(chan int, config.BufferSize)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				result[idx] = mapper(slice[idx])
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := range slice {
			jobs <- i
		}
	}()

	wg.Wait()
	return result
}

// FilterParallel parallel filtering
func FilterParallel[T any](slice []T, predicate Predicate[T], config ParallelConfig) []T {
	if slice == nil || len(slice) == 0 {
		return nil
	}

	if len(slice) < config.WorkerCount {
		return Filter(slice, predicate)
	}

	type result struct {
		index int
		item  T
		keep  bool
	}

	jobs := make(chan int, config.BufferSize)
	results := make(chan result, config.BufferSize)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				item := slice[idx]
				keep := predicate(item)
				results <- result{index: idx, item: item, keep: keep}
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := range slice {
			jobs <- i
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sort results by index and filter
	resultMap := make(map[int]result)
	for res := range results {
		resultMap[res.index] = res
	}

	var filtered []T
	for i := 0; i < len(slice); i++ {
		if res, exists := resultMap[i]; exists && res.keep {
			filtered = append(filtered, res.item)
		}
	}

	return filtered
}

// ReduceParallel parallel reduction (for associative operations)
func ReduceParallel[T any](slice []T, reducer func(T, T) T, identity T, config ParallelConfig) T {
	if slice == nil || len(slice) == 0 {
		return identity
	}

	if len(slice) == 1 {
		return slice[0]
	}

	if len(slice) < config.WorkerCount {
		result := identity
		for _, item := range slice {
			result = reducer(result, item)
		}
		return result
	}

	// Split into chunks
	chunkSize := len(slice) / config.WorkerCount
	if chunkSize == 0 {
		chunkSize = 1
	}

	chunks := Chunk(slice, chunkSize)
	results := make([]T, len(chunks))
	var wg sync.WaitGroup

	// Process each chunk in parallel
	for i, chunk := range chunks {
		wg.Add(1)
		go func(idx int, ch []T) {
			defer wg.Done()
			result := identity
			for _, item := range ch {
				result = reducer(result, item)
			}
			results[idx] = result
		}(i, chunk)
	}

	wg.Wait()

	// Fold chunk results
	finalResult := identity
	for _, result := range results {
		finalResult = reducer(finalResult, result)
	}

	return finalResult
}

// Pipeline represents a data processing pipeline
type Pipeline[T any] struct {
	data   []T
	config ParallelConfig
}

// NewPipeline creates a new pipeline
func NewPipeline[T any](data []T) *Pipeline[T] {
	return &Pipeline[T]{
		data:   data,
		config: DefaultParallelConfig(),
	}
}

// WithConfig sets the configuration
func (p *Pipeline[T]) WithConfig(config ParallelConfig) *Pipeline[T] {
	p.config = config
	return p
}

// Map applies a transformation function
func (p *Pipeline[T]) Map(mapper Mapper[T, T]) *Pipeline[T] {
	p.data = MapParallelWithConfig(p.data, mapper, p.config)
	return p
}

// Filter applies a filtering function
func (p *Pipeline[T]) Filter(predicate Predicate[T]) *Pipeline[T] {
	p.data = FilterParallel(p.data, predicate, p.config)
	return p
}

// Reduce reduces the data
func (p *Pipeline[T]) Reduce(reducer func(T, T) T, identity T) T {
	return ReduceParallel(p.data, reducer, identity, p.config)
}

// Collect returns the result
func (p *Pipeline[T]) Collect() []T {
	return p.data
}

// ForEachParallel executes a function for each element in parallel
func ForEachParallel[T any](slice []T, action func(T), config ParallelConfig) {
	if slice == nil || len(slice) == 0 {
		return
	}

	jobs := make(chan T, config.BufferSize)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				action(item)
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for _, item := range slice {
			jobs <- item
		}
	}()

	wg.Wait()
}

// MapWithContext executes a function for each element in parallel with a context
func MapWithContext[T, R any](ctx context.Context, slice []T, mapper func(context.Context, T) (R, error), config ParallelConfig) ([]R, error) {
	if slice == nil || len(slice) == 0 {
		return nil, nil
	}

	result := make([]R, len(slice))
	jobs := make(chan int, config.BufferSize)
	errors := make(chan error, len(slice))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case idx, ok := <-jobs:
					if !ok {
						return
					}

					res, err := mapper(ctx, slice[idx])
					if err != nil {
						errors <- err
						return
					}
					result[idx] = res

				case <-ctx.Done():
					errors <- ctx.Err()
					return
				}
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := range slice {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for completion or error
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		select {
		case err := <-errors:
			return nil, err
		default:
			return result, nil
		}
	case err := <-errors:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// BatchProcessor processes data in batches
type BatchProcessor[T, R any] struct {
	batchSize   int
	processor   func([]T) ([]R, error)
	parallelism int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor[T, R any](batchSize int, processor func([]T) ([]R, error)) *BatchProcessor[T, R] {
	return &BatchProcessor[T, R]{
		batchSize:   batchSize,
		processor:   processor,
		parallelism: runtime.NumCPU(),
	}
}

// WithParallelism sets the level of parallelism
func (bp *BatchProcessor[T, R]) WithParallelism(parallelism int) *BatchProcessor[T, R] {
	bp.parallelism = parallelism
	return bp
}

// Process processes data in batches
func (bp *BatchProcessor[T, R]) Process(ctx context.Context, data []T) ([]R, error) {
	if len(data) == 0 {
		return nil, nil
	}

	batches := Chunk(data, bp.batchSize)
	results := make([][]R, len(batches))

	jobs := make(chan int, len(batches))
	errors := make(chan error, len(batches))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < min(bp.parallelism, len(batches)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case batchIdx, ok := <-jobs:
					if !ok {
						return
					}

					result, err := bp.processor(batches[batchIdx])
					if err != nil {
						errors <- err
						return
					}
					results[batchIdx] = result

				case <-ctx.Done():
					errors <- ctx.Err()
					return
				}
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := range batches {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for completion
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		select {
		case err := <-errors:
			return nil, err
		default:
			// Fold results
			var finalResult []R
			for _, batch := range results {
				finalResult = append(finalResult, batch...)
			}
			return finalResult, nil
		}
	case err := <-errors:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RateLimiter limits the execution speed
type RateLimiter struct {
	tokens chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rps),
	}

	// Fill tokens
	for i := 0; i < rps; i++ {
		rl.tokens <- struct{}{}
	}

	// Refill tokens with the given rate
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rps))
		defer ticker.Stop()

		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Channel is full, skip
			}
		}
	}()

	return rl
}

// Wait waits for an available token
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// MapWithRateLimit applies a function with a rate limit
func MapWithRateLimit[T, R any](ctx context.Context, slice []T, mapper func(T) (R, error), rps int) ([]R, error) {
	if slice == nil || len(slice) == 0 {
		return nil, nil
	}

	limiter := NewRateLimiter(rps)
	result := make([]R, len(slice))

	for i, item := range slice {
		if err := limiter.Wait(ctx); err != nil {
			return nil, err
		}

		res, err := mapper(item)
		if err != nil {
			return nil, err
		}
		result[i] = res
	}

	return result, nil
}
