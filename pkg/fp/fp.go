// Package fp предоставляет мощную библиотеку функционального программирования для Go
// с современными дженериками, включающую все необходимые инструменты для
// функциональной обработки данных.
//
// Основные возможности:
//   - Базовые функции высшего порядка (Map, Filter, Reduce)
//   - Композиция и каррирование функций (Pipe, Compose, Curry)
//   - Безопасная работа с nullable значениями (Optional, Result)
//   - Утилиты для работы с коллекциями (GroupBy, Chunk, Partition)
//   - Параллельная обработка данных с контекстом
//   - Stream API для ленивых вычислений
//   - Множество дополнительных утилит
//
// Пример базового использования:
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
// Пример с Pipeline:
//
//	result := fp.NewPipeline([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
//	    Filter(func(x int) bool { return x%2 == 0 }).
//	    Map(func(x int) int { return x * x }).
//	    Collect()
//	// [4, 16, 36, 64, 100]
//
// Пример с Optional:
//
//	user := fp.Some("John")
//	name := user.Map(strings.ToUpper).GetOrElse("Unknown")
//	// "JOHN"
//
// Пример с Stream:
//
//	numbers := fp.RangeStream(1, 11).
//	    Filter(func(x int) bool { return x%2 == 0 }).
//	    Map(func(x int) int { return x * x }).
//	    Take(3).
//	    Collect()
//	// [4, 16, 36]
//
package fp

// Версия библиотеки
const Version = "1.0.0"

// Информация о библиотеке
var (
	Name        = "Functional Programming Library for Go"
	Description = "Comprehensive functional programming toolkit with generics"
	Author      = "NLG.KZ Backend Team"
)

// GetVersion возвращает версию библиотеки
func GetVersion() string {
	return Version
}

// GetInfo возвращает информацию о библиотеке
func GetInfo() map[string]string {
	return map[string]string{
		"name":        Name,
		"version":     Version,
		"description": Description,
		"author":      Author,
	}
}

// Экспорт основных типов для удобства использования

// Основные функциональные типы
type (
	// Func представляет функцию одного аргумента
	Func[T, R any] func(T) R
	
	// BiFunc представляет функцию двух аргументов
	BiFunc[T, R, S any] func(T, R) S
	
	// Consumer представляет функцию-потребитель
	Consumer[T any] func(T)
	
	// Supplier представляет функцию-поставщик
	Supplier[T any] func() T
	
	// Runnable представляет функцию без аргументов и возвращаемого значения
	Runnable func()
)

// Утилиты для создания функций

// AsFunc преобразует функцию в Func
func AsFunc[T, R any](f func(T) R) Func[T, R] {
	return f
}

// AsBiFunc преобразует функцию в BiFunc
func AsBiFunc[T, R, S any](f func(T, R) S) BiFunc[T, R, S] {
	return f
}

// AsConsumer преобразует функцию в Consumer
func AsConsumer[T any](f func(T)) Consumer[T] {
	return f
}

// AsSupplier преобразует функцию в Supplier
func AsSupplier[T any](f func() T) Supplier[T] {
	return f
}

// AsRunnable преобразует функцию в Runnable
func AsRunnable(f func()) Runnable {
	return f
}

// Цепочка вызовов для Func

// AndThen создает композицию функций (сначала текущая, потом after)
func (f Func[T, R]) AndThen(after Func[R, any]) Func[T, any] {
	return func(t T) any {
		return after(f(t))
	}
}

// Compose создает композицию функций (сначала before, потом текущая)
func (f Func[T, R]) Compose(before Func[any, T]) Func[any, R] {
	return func(v any) R {
		return f(before(v))
	}
}

// Константы для часто используемых операций

var (
	// IntAdd складывает два числа
	IntAdd = func(a, b int) int { return a + b }
	
	// IntMul умножает два числа
	IntMul = func(a, b int) int { return a * b }
	
	// IntMax возвращает максимум из двух чисел
	IntMax = func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	
	// IntMin возвращает минимум из двух чисел
	IntMin = func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	
	// StringConcat объединяет две строки
	StringConcat = func(a, b string) string { return a + b }
	
	// IsPositive проверяет, что число положительное
	IsPositive = func(x int) bool { return x > 0 }
	
	// IsNegative проверяет, что число отрицательное
	IsNegative = func(x int) bool { return x < 0 }
	
	// IsEven проверяет, что число четное
	IsEven = func(x int) bool { return x%2 == 0 }
	
	// IsOdd проверяет, что число нечетное
	IsOdd = func(x int) bool { return x%2 != 0 }
)

// Быстрые операции для чисел

// Numbers предоставляет утилиты для работы с числами
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

// Strings предоставляет утилиты для работы со строками
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

// Slices предоставляет утилиты для работы со слайсами
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

// Predicates предоставляет готовые предикаты
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

// Collectors предоставляет готовые коллекторы
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

// Быстрые алиасы для часто используемых функций

// F создает функцию из лямбды (алиас для удобства)
func F[T, R any](f func(T) R) Func[T, R] {
	return f
}

// P создает предикат (алиас для удобства)
func P[T any](f func(T) bool) Predicate[T] {
	return f
}

// C создает consumer (алиас для удобства)
func C[T any](f func(T)) Consumer[T] {
	return f
}

// S создает supplier (алиас для удобства)
func S[T any](f func() T) Supplier[T] {
	return f
}

// Дополнительные утилиты

// Times выполняет функцию n раз
func Times(n int, action Runnable) {
	for i := 0; i < n; i++ {
		action()
	}
}

// TimesWithIndex выполняет функцию n раз с передачей индекса
func TimesWithIndex(n int, action func(int)) {
	for i := 0; i < n; i++ {
		action(i)
	}
}

// Benchmark измеряет время выполнения функции
func Benchmark[T any](name string, fn func() T) T {
	start := time.Now()
	result := fn()
	duration := time.Since(start)
	fmt.Printf("Benchmark %s: %v\n", name, duration)
	return result
}

// Lazy создает ленивое значение
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

// Необходимые импорты для компиляции
import (
	"fmt"
	"strings"
	"sync"
	"time"
)
