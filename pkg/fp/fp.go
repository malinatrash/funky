// Package fp provides a powerful functional programming library for Go
// with generics, including all necessary tools for functional data processing.
//
// Main features:
//   - Basic higher-order functions (Map, Filter, Reduce)
//   - Function composition and currying (Pipe, Compose, Curry)
//   - Safe nullable value handling (Optional, Result)
//   - Collection utilities (GroupBy, Chunk, Partition)
//   - Parallel data processing with context
//   - Lazy evaluation with Stream API
//   - A set of additional utilities
//
// Example of basic usage:
//
//	import "github.com/malinatrash/funky/pkg/fp"
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := fp.Map(numbers, func(x int) int { return x * 2 })
//	// [2, 4, 6, 8, 10]
//
//	evens := fp.Filter(numbers, func(x int) bool { return x%2 == 0 })
//	// [2, 4]
//
//	sum := fp.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
//	// 15
//
// Example with Pipeline:
//
//	result := fp.NewPipeline([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
//	    Filter(func(x int) bool { return x%2 == 0 }).
//	    Map(func(x int) int { return x * x }).
//	    Collect()
//	// [4, 16, 36, 64, 100]
//
// Example with Optional:
//
//	user := fp.Some("John")
//	name := user.Map(strings.ToUpper).GetOrElse("Unknown")
//	// "JOHN"
//
// Example with Stream:
//
//	numbers := fp.RangeStream(1, 11).
//	    Filter(func(x int) bool { return x%2 == 0 }).
//	    Map(func(x int) int { return x * x }).
//	    Take(3).
//	    Collect()
//	// [4, 16, 36]
package fp

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Basic functional types
type (
	// Func represents a function with one argument
	Func[T, R any] func(T) R

	// BiFunc represents a function with two arguments
	BiFunc[T, R, S any] func(T, R) S

	// Consumer represents a function that consumes a value
	Consumer[T any] func(T)

	// Supplier represents a function that provides a value
	Supplier[T any] func() T

	// Runnable represents a function without arguments and return value
	Runnable func()
)

// Utilities for creating functions

// AsFunc converts a function to Func
func AsFunc[T, R any](f func(T) R) Func[T, R] {
	return f
}

// AsBiFunc converts a function to BiFunc
func AsBiFunc[T, R, S any](f func(T, R) S) BiFunc[T, R, S] {
	return f
}

// AsConsumer converts a function to Consumer
func AsConsumer[T any](f func(T)) Consumer[T] {
	return f
}

// AsSupplier converts a function to Supplier
func AsSupplier[T any](f func() T) Supplier[T] {
	return f
}

// AsRunnable converts a function to Runnable
func AsRunnable(f func()) Runnable {
	return f
}

// AndThen creates a function composition (first current, then after)
func (f Func[T, R]) AndThen(after Func[R, any]) Func[T, any] {
	return func(t T) any {
		return after(f(t))
	}
}

// Compose creates a function composition (first before, then current)
func (f Func[T, R]) Compose(before Func[any, T]) Func[any, R] {
	return func(v any) R {
		return f(before(v))
	}
}

// Constants for frequently used operations
var (
	// IntAdd adds two numbers
	IntAdd = func(a, b int) int { return a + b }

	// IntMul multiplies two numbers
	IntMul = func(a, b int) int { return a * b }

	// IntMax returns the maximum of two numbers
	IntMax = func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// IntMin returns the minimum of two numbers
	IntMin = func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// StringConcat concatenates two strings
	StringConcat = func(a, b string) string { return a + b }

	// IsPositive checks if a number is positive
	IsPositive = func(x int) bool { return x > 0 }

	// IsNegative checks if a number is negative
	IsNegative = func(x int) bool { return x < 0 }

	// IsEven checks if a number is even
	IsEven = func(x int) bool { return x%2 == 0 }

	// IsOdd checks if a number is odd
	IsOdd = func(x int) bool { return x%2 != 0 }
)

// Fast operations for numbers

// Numbers provides utilities for working with numbers
var Numbers = struct {
	Add      func(a, b int) int
	Subtract func(a, b int) int
	Multiply func(a, b int) int
	Divide   func(a, b int) int
	Mod      func(a, b int) int
	Abs      func(x int) int
	Square   func(x int) int
	Cube     func(x int) int
	Double   func(x int) int
	Half     func(x int) int
}{
	Add:      func(a, b int) int { return a + b },
	Subtract: func(a, b int) int { return a - b },
	Multiply: func(a, b int) int { return a * b },
	Divide:   func(a, b int) int { return a / b },
	Mod:      func(a, b int) int { return a % b },
	Abs: func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	},
	Square: func(x int) int { return x * x },
	Cube:   func(x int) int { return x * x * x },
	Double: func(x int) int { return x * 2 },
	Half:   func(x int) int { return x / 2 },
}

// Strings provides utilities for working with strings
var Strings = struct {
	ToUpper    func(s string) string
	ToLower    func(s string) string
	Trim       func(s string) string
	Length     func(s string) int
	IsEmpty    func(s string) bool
	IsNotEmpty func(s string) bool
	Reverse    func(s string) string
}{
	ToUpper:    func(s string) string { return strings.ToUpper(s) },
	ToLower:    func(s string) string { return strings.ToLower(s) },
	Trim:       func(s string) string { return strings.TrimSpace(s) },
	Length:     func(s string) int { return len(s) },
	IsEmpty:    func(s string) bool { return len(s) == 0 },
	IsNotEmpty: func(s string) bool { return len(s) > 0 },
	Reverse: func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	},
}

// Slices provides utilities for working with slices
var Slices = struct {
	IsEmpty    func(slice []any) bool
	IsNotEmpty func(slice []any) bool
	Length     func(slice []any) int
	First      func(slice []any) (any, bool)
	Last       func(slice []any) (any, bool)
}{
	IsEmpty:    func(slice []any) bool { return len(slice) == 0 },
	IsNotEmpty: func(slice []any) bool { return len(slice) > 0 },
	Length:     func(slice []any) int { return len(slice) },
	First: func(slice []any) (any, bool) {
		if len(slice) == 0 {
			return nil, false
		}
		return slice[0], true
	},
	Last: func(slice []any) (any, bool) {
		if len(slice) == 0 {
			return nil, false
		}
		return slice[len(slice)-1], true
	},
}

// Predicates provides ready-to-use predicates
var Predicates = struct {
	True     func(any) bool
	False    func(any) bool
	NotNil   func(any) bool
	IsNil    func(any) bool
	IsZero   func(any) bool
	NotZero  func(any) bool
	IsEmpty  func(string) bool
	NotEmpty func(string) bool
}{
	True:     func(any) bool { return true },
	False:    func(any) bool { return false },
	NotNil:   func(v any) bool { return v != nil },
	IsNil:    func(v any) bool { return v == nil },
	IsZero:   func(v any) bool { return v == nil || v == 0 || v == "" },
	NotZero:  func(v any) bool { return v != nil && v != 0 && v != "" },
	IsEmpty:  func(s string) bool { return len(s) == 0 },
	NotEmpty: func(s string) bool { return len(s) > 0 },
}

// Collectors provides ready-to-use collectors
var Collectors = struct {
	ToSlice   func() func([]any) []any
	ToSet     func() func([]any) map[any]bool
	Counting  func() func([]any) int
	Joining   func(delimiter string) func([]string) string
	Averaging func() func([]int) float64
}{
	ToSlice: func() func([]any) []any {
		return func(slice []any) []any { return slice }
	},
	ToSet: func() func([]any) map[any]bool {
		return func(slice []any) map[any]bool {
			set := make(map[any]bool)
			for _, item := range slice {
				set[item] = true
			}
			return set
		}
	},
	Counting: func() func([]any) int {
		return func(slice []any) int { return len(slice) }
	},
	Joining: func(delimiter string) func([]string) string {
		return func(strings []string) string {
			return Join(strings, delimiter)
		}
	},
	Averaging: func() func([]int) float64 {
		return func(numbers []int) float64 {
			if len(numbers) == 0 {
				return 0
			}
			sum := Sum(numbers)
			return float64(sum) / float64(len(numbers))
		}
	},
}

// Fast provides aliases for frequently used functions

// F creates a function from a lambda (alias for convenience)
func F[T, R any](f func(T) R) Func[T, R] {
	return f
}

// P creates a predicate (alias for convenience)
func P[T any](f func(T) bool) Predicate[T] {
	return f
}

// C creates a consumer (alias for convenience)
func C[T any](f func(T)) Consumer[T] {
	return f
}

// S creates a supplier (alias for convenience)
func S[T any](f func() T) Supplier[T] {
	return f
}

// Additional utilities

// Times executes a function n times
func Times(n int, action Runnable) {
	for i := 0; i < n; i++ {
		action()
	}
}

// TimesWithIndex executes a function n times with index
func TimesWithIndex(n int, action func(int)) {
	for i := 0; i < n; i++ {
		action(i)
	}
}

// Benchmark measures the execution time of a function
func Benchmark[T any](name string, fn func() T) T {
	start := time.Now()
	result := fn()
	duration := time.Since(start)
	fmt.Printf("Benchmark %s: %v\n", name, duration)
	return result
}

// Lazy creates a lazy value
func Lazy[T any](supplier Supplier[T]) func() T {
	var value T
	var computed bool
	var mutex sync.Mutex

	return func() T {
		mutex.Lock()
		defer mutex.Unlock()

		if !computed {
			value = supplier()
			computed = true
		}
		return value
	}
}
