package fp

import "fmt"

// Optional represents a value that may be absent
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates Optional with a value
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// Empty creates an empty Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// Of creates Optional from a value (nil becomes Empty)
func Of[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return Some(*value)
}

// OfNillable creates Optional from a value that may be nil
func OfNillable[T any](value T, isNil bool) Optional[T] {
	if isNil {
		return Empty[T]()
	}
	return Some(value)
}

// IsPresent checks if the Optional has a value
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty checks if the Optional is empty
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value or panics if the Optional is empty
func (o Optional[T]) Get() T {
	if !o.present {
		panic("Optional is empty")
	}
	return o.value
}

// GetOrElse returns the value or a default value
func (o Optional[T]) GetOrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// GetOrElseGet returns the value or the result of a function
func (o Optional[T]) GetOrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// OrElse returns the current Optional or another Optional if the current is empty
func (o Optional[T]) OrElse(other Optional[T]) Optional[T] {
	if o.present {
		return o
	}
	return other
}

// OrElseGet returns the current Optional or the result of a function if the current is empty
func (o Optional[T]) OrElseGet(supplier func() Optional[T]) Optional[T] {
	if o.present {
		return o
	}
	return supplier()
}

// Map applies a function to the value if it is present
func (o Optional[T]) Map(mapper func(T) T) Optional[T] {
	if !o.present {
		return Empty[T]()
	}
	return Some(mapper(o.value))
}

// MapTo applies a function and returns an Optional of another type
func MapTo[T, R any](o Optional[T], mapper func(T) R) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return Some(mapper(o.value))
}

// FlatMap applies a function that returns an Optional
func (o Optional[T]) FlatMap(mapper func(T) Optional[T]) Optional[T] {
	if !o.present {
		return Empty[T]()
	}
	return mapper(o.value)
}

// FlatMapTo applies a function that returns an Optional of another type
func FlatMapTo[T, R any](o Optional[T], mapper func(T) Optional[R]) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return mapper(o.value)
}

// Filter filters the value by a predicate
func (o Optional[T]) Filter(predicate Predicate[T]) Optional[T] {
	if !o.present || !predicate(o.value) {
		return Empty[T]()
	}
	return o
}

// IfPresent executes a function if the Optional has a value
func (o Optional[T]) IfPresent(consumer func(T)) {
	if o.present {
		consumer(o.value)
	}
}

// IfPresentOrElse executes one of the functions depending on the presence of a value
func (o Optional[T]) IfPresentOrElse(consumer func(T), emptyAction func()) {
	if o.present {
		consumer(o.value)
	} else {
		emptyAction()
	}
}

// ToPointer returns a pointer to the value or nil
func (o Optional[T]) ToPointer() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

// ToSlice returns a slice with one element or an empty slice
func (o Optional[T]) ToSlice() []T {
	if !o.present {
		return []T{}
	}
	return []T{o.value}
}

// String returns a string representation of the Optional
func (o Optional[T]) String() string {
	if !o.present {
		return "Optional.empty"
	}
	return fmt.Sprintf("Optional[%v]", o.value)
}

// Equals compares two Optional
func (o Optional[T]) Equals(other Optional[T], equals Equality[T]) bool {
	if o.present != other.present {
		return false
	}
	if !o.present {
		return true // both are empty
	}
	return equals(o.value, other.value)
}

// Result represents the result of an operation that may fail
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a Result with an error
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: err}
}

// TryFrom creates a Result from a function that may return an error
func TryFrom[T any](fn TryFunc[T]) Result[T] {
	value, err := fn()
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// IsOk checks if the Result is successful
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr checks if the Result has an error
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value or panics on error
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("called Unwrap on an Err value: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the value or a default value on error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value or the result of a function on error
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.err != nil {
		return fn(r.err)
	}
	return r.value
}

// Error returns the error
func (r Result[T]) Error() error {
	return r.err
}

// Map applies a function to the value if there is no error
func (r Result[T]) Map(mapper func(T) T) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return Ok(mapper(r.value))
}

// MapTo applies a function and returns a Result of another type
func MapToResult[T, R any](r Result[T], mapper func(T) R) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return Ok(mapper(r.value))
}

// FlatMap applies a function that returns a Result
func (r Result[T]) FlatMap(mapper func(T) Result[T]) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return mapper(r.value)
}

// FlatMapTo applies a function that returns a Result of another type
func FlatMapToResult[T, R any](r Result[T], mapper func(T) Result[R]) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return mapper(r.value)
}

// MapErr applies a function to the error
func (r Result[T]) MapErr(mapper func(error) error) Result[T] {
	if r.err != nil {
		return Err[T](mapper(r.err))
	}
	return r
}

// ToOptional converts Result to Optional (ignores error)
func (r Result[T]) ToOptional() Optional[T] {
	if r.err != nil {
		return Empty[T]()
	}
	return Some(r.value)
}

// FromOptional creates a Result from an Optional
func FromOptional[T any](opt Optional[T], err error) Result[T] {
	if opt.IsEmpty() {
		return Err[T](err)
	}
	return Ok(opt.Get())
}

// Utility functions for working with slices of Optional and Result

// FilterSome filters only non-empty Optionals
func FilterSome[T any](optionals []Optional[T]) []T {
	var result []T
	for _, opt := range optionals {
		if opt.IsPresent() {
			result = append(result, opt.Get())
		}
	}
	return result
}

// FilterOk filters only successful Results
func FilterOk[T any](results []Result[T]) []T {
	var result []T
	for _, res := range results {
		if res.IsOk() {
			result = append(result, res.value)
		}
	}
	return result
}

// FilterErr filters only errors from Result
func FilterErr[T any](results []Result[T]) []error {
	var errors []error
	for _, res := range results {
		if res.IsErr() {
			errors = append(errors, res.err)
		}
	}
	return errors
}

// Sequence converts a slice of Optionals to an Optional of a slice
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

// SequenceResults converts a slice of Results to a Result of a slice
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
