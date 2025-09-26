package fp

// Предикатные функции
type Predicate[T any] func(T) bool
type PredicateWithIndex[T any] func(T, int) bool

// Трансформирующие функции
type Mapper[T, R any] func(T) R
type MapperWithIndex[T, R any] func(T, int) R

// Редуцирующие функции
type Reducer[T, R any] func(R, T) R
type ReducerWithIndex[T, R any] func(R, T, int) R

// Функции сравнения
type Comparator[T any] func(T, T) int
type Equality[T any] func(T, T) bool

// Функции группировки
type KeyExtractor[T any, K comparable] func(T) K

// Функции для работы с ошибками
type ErrorHandler func(error)
type TryFunc[T any] func() (T, error)

// Константы для сравнения
const (
	Less    = -1
	Equal   = 0
	Greater = 1
)

// Общие предикаты
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

// Identity функция
func Identity[T any](item T) T {
	return item
}

// Константная функция
func Const[T, R any](value R) func(T) R {
	return func(T) R {
		return value
	}
}
