package fp

// Predicate functions
type Predicate[T any] func(T) bool
type PredicateWithIndex[T any] func(T, int) bool

// Mapper functions
type Mapper[T, R any] func(T) R
type MapperWithIndex[T, R any] func(T, int) R

// Reducer functions
type Reducer[T, R any] func(R, T) R
type ReducerWithIndex[T, R any] func(R, T, int) R

// Comparator functions
type Comparator[T any] func(T, T) int
type Equality[T any] func(T, T) bool

// Key extractor function
type KeyExtractor[T any, K comparable] func(T) K

// Error handler function
type ErrorHandler func(error)
type TryFunc[T any] func() (T, error)

// Comparison constants
const (
	Less    = -1
	Equal   = 0
	Greater = 1
)

// Common predicates
func IsNil[T any](item *T) bool {
	return item == nil
}

func IsNotNil[T any](item *T) bool {
	return item != nil
}

func IsZero[T comparable](item T) bool {
	var zero T
	return item == zero
}

func IsNotZero[T comparable](item T) bool {
	return !IsZero(item)
}

// Identity function
func Identity[T any](item T) T {
	return item
}

// Constant function
func Const[T, R any](value R) func(T) R {
	return func(T) R {
		return value
	}
}
