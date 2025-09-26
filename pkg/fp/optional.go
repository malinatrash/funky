package fp

import "fmt"

// Optional представляет значение, которое может отсутствовать
type Optional[T any] struct {
	value   T
	present bool
}

// Some создает Optional с значением
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// Empty создает пустой Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// Of создает Optional из значения (nil становится Empty)
func Of[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return Some(*value)
}

// OfNillable создает Optional из значения, которое может быть nil
func OfNillable[T any](value T, isNil bool) Optional[T] {
	if isNil {
		return Empty[T]()
	}
	return Some(value)
}

// IsPresent проверяет, есть ли значение
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty проверяет, пустой ли Optional
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get возвращает значение или панику, если значения нет
func (o Optional[T]) Get() T {
	if !o.present {
		panic("Optional is empty")
	}
	return o.value
}

// GetOrElse возвращает значение или значение по умолчанию
func (o Optional[T]) GetOrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// GetOrElseGet возвращает значение или результат функции
func (o Optional[T]) GetOrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// OrElse возвращает текущий Optional или другой, если текущий пустой
func (o Optional[T]) OrElse(other Optional[T]) Optional[T] {
	if o.present {
		return o
	}
	return other
}

// OrElseGet возвращает текущий Optional или результат функции
func (o Optional[T]) OrElseGet(supplier func() Optional[T]) Optional[T] {
	if o.present {
		return o
	}
	return supplier()
}

// Map применяет функцию к значению, если оно есть
func (o Optional[T]) Map(mapper func(T) T) Optional[T] {
	if !o.present {
		return Empty[T]()
	}
	return Some(mapper(o.value))
}

// MapTo применяет функцию и возвращает Optional другого типа
func MapTo[T, R any](o Optional[T], mapper func(T) R) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return Some(mapper(o.value))
}

// FlatMap применяет функцию, возвращающую Optional
func (o Optional[T]) FlatMap(mapper func(T) Optional[T]) Optional[T] {
	if !o.present {
		return Empty[T]()
	}
	return mapper(o.value)
}

// FlatMapTo применяет функцию, возвращающую Optional другого типа
func FlatMapTo[T, R any](o Optional[T], mapper func(T) Optional[R]) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return mapper(o.value)
}

// Filter фильтрует значение по предикату
func (o Optional[T]) Filter(predicate Predicate[T]) Optional[T] {
	if !o.present || !predicate(o.value) {
		return Empty[T]()
	}
	return o
}

// IfPresent выполняет функцию, если значение есть
func (o Optional[T]) IfPresent(consumer func(T)) {
	if o.present {
		consumer(o.value)
	}
}

// IfPresentOrElse выполняет одну из функций в зависимости от наличия значения
func (o Optional[T]) IfPresentOrElse(consumer func(T), emptyAction func()) {
	if o.present {
		consumer(o.value)
	} else {
		emptyAction()
	}
}

// ToPointer возвращает указатель на значение или nil
func (o Optional[T]) ToPointer() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

// ToSlice возвращает слайс с одним элементом или пустой слайс
func (o Optional[T]) ToSlice() []T {
	if !o.present {
		return []T{}
	}
	return []T{o.value}
}

// String возвращает строковое представление Optional
func (o Optional[T]) String() string {
	if !o.present {
		return "Optional.empty"
	}
	return fmt.Sprintf("Optional[%v]", o.value)
}

// Equals сравнивает два Optional
func (o Optional[T]) Equals(other Optional[T], equals Equality[T]) bool {
	if o.present != other.present {
		return false
	}
	if !o.present {
		return true // оба пустые
	}
	return equals(o.value, other.value)
}

// Result представляет результат операции, которая может завершиться ошибкой
type Result[T any] struct {
	value T
	err   error
}

// Ok создает успешный Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err создает Result с ошибкой
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: err}
}

// TryFrom создает Result из функции, которая может вернуть ошибку
func TryFrom[T any](fn TryFunc[T]) Result[T] {
	value, err := fn()
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// IsOk проверяет, успешен ли результат
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr проверяет, есть ли ошибка
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap возвращает значение или панику при ошибке
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("called Unwrap on an Err value: %v", r.err))
	}
	return r.value
}

// UnwrapOr возвращает значение или значение по умолчанию при ошибке
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse возвращает значение или результат функции при ошибке
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.err != nil {
		return fn(r.err)
	}
	return r.value
}

// Error возвращает ошибку
func (r Result[T]) Error() error {
	return r.err
}

// Map применяет функцию к значению, если нет ошибки
func (r Result[T]) Map(mapper func(T) T) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return Ok(mapper(r.value))
}

// MapTo применяет функцию и возвращает Result другого типа
func MapToResult[T, R any](r Result[T], mapper func(T) R) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return Ok(mapper(r.value))
}

// FlatMap применяет функцию, возвращающую Result
func (r Result[T]) FlatMap(mapper func(T) Result[T]) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return mapper(r.value)
}

// FlatMapTo применяет функцию, возвращающую Result другого типа
func FlatMapToResult[T, R any](r Result[T], mapper func(T) Result[R]) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return mapper(r.value)
}

// MapErr применяет функцию к ошибке
func (r Result[T]) MapErr(mapper func(error) error) Result[T] {
	if r.err != nil {
		return Err[T](mapper(r.err))
	}
	return r
}

// ToOptional преобразует Result в Optional (игнорирует ошибку)
func (r Result[T]) ToOptional() Optional[T] {
	if r.err != nil {
		return Empty[T]()
	}
	return Some(r.value)
}

// FromOptional создает Result из Optional
func FromOptional[T any](opt Optional[T], err error) Result[T] {
	if opt.IsEmpty() {
		return Err[T](err)
	}
	return Ok(opt.Get())
}

// Utility functions for working with slices of Optional and Result

// FilterSome фильтрует только непустые Optional
func FilterSome[T any](optionals []Optional[T]) []T {
	var result []T
	for _, opt := range optionals {
		if opt.IsPresent() {
			result = append(result, opt.Get())
		}
	}
	return result
}

// FilterOk фильтрует только успешные Result
func FilterOk[T any](results []Result[T]) []T {
	var result []T
	for _, res := range results {
		if res.IsOk() {
			result = append(result, res.value)
		}
	}
	return result
}

// FilterErr фильтрует только ошибки из Result
func FilterErr[T any](results []Result[T]) []error {
	var errors []error
	for _, res := range results {
		if res.IsErr() {
			errors = append(errors, res.err)
		}
	}
	return errors
}

// Sequence преобразует слайс Optional в Optional слайса
func Sequence[T any](optionals []Optional[T]) Optional[[]T] {
	result := make([]T, 0, len(optionals))
	for _, opt := range optionals {
		if opt.IsEmpty() {
			return Empty[[]T]()
		}
		result = append(result, opt.Get())
	}
	return Some(result)
}

// SequenceResults преобразует слайс Result в Result слайса
func SequenceResults[T any](results []Result[T]) Result[[]T] {
	result := make([]T, 0, len(results))
	for _, res := range results {
		if res.IsErr() {
			return Err[[]T](res.err)
		}
		result = append(result, res.value)
	}
	return Ok(result)
}
