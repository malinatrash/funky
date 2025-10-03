package fp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Try executes a function and returns Result
func Try[T any](fn func() (T, error)) Result[T] {
	return TryFrom(fn)
}

// TryVoid executes a function without a return value
func TryVoid(fn func() error) Result[struct{}] {
	err := fn()
	if err != nil {
		return Err[struct{}](err)
	}
	return Ok(struct{}{})
}

// Tap executes a side effect and returns the original value
func Tap[T any](value T, sideEffect func(T)) T {
	sideEffect(value)
	return value
}

// TapIf executes a side effect only if the condition is true
func TapIf[T any](value T, condition bool, sideEffect func(T)) T {
	if condition {
		sideEffect(value)
	}
	return value
}

// TapWhen executes a side effect only if the predicate is true
func TapWhen[T any](value T, predicate Predicate[T], sideEffect func(T)) T {
	if predicate(value) {
		sideEffect(value)
	}
	return value
}

// Let allows executing an operation on a value
func Let[T, R any](value T, transform func(T) R) R {
	return transform(value)
}

// Also executes an operation and returns the original value
func Also[T any](value T, operation func(T)) T {
	operation(value)
	return value
}

// TakeIf returns the value if the predicate is true, otherwise nil
func TakeIf[T any](value T, predicate Predicate[T]) *T {
	if predicate(value) {
		return &value
	}
	return nil
}

// TakeUnless returns the value if the predicate is false, otherwise nil
func TakeUnless[T any](value T, predicate Predicate[T]) *T {
	if !predicate(value) {
		return &value
	}
	return nil
}

// Coalesce returns the first non-nil value
func Coalesce[T any](values ...*T) *T {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// CoalesceFunc returns the first non-nil value from functions
func CoalesceFunc[T any](suppliers ...func() *T) *T {
	for _, supplier := range suppliers {
		if value := supplier(); value != nil {
			return value
		}
	}
	return nil
}

// DefaultIfNil returns the default value if the pointer is nil
func DefaultIfNil[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// Ptr creates a pointer to a value
func Ptr[T any](value T) *T {
	return &value
}

// Deref dereferences a pointer with a default value
func Deref[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// SafeDeref safely dereferences a pointer
func SafeDeref[T any](ptr *T) (T, bool) {
	if ptr == nil {
		var zero T
		return zero, false
	}
	return *ptr, true
}

// Contains checks if an element exists in a slice
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsBy checks if an element exists in a slice by predicate
func ContainsBy[T any](slice []T, predicate Predicate[T]) bool {
	return Any(slice, predicate)
}

// IndexOf returns the index of an element in a slice
func IndexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index of an element in a slice
func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// Remove removes the first occurrence of an element
func Remove[T comparable](slice []T, item T) []T {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// RemoveAll removes all occurrences of an element
func RemoveAll[T comparable](slice []T, item T) []T {
	return Filter(slice, func(v T) bool { return v != item })
}

// RemoveAt removes an element by index
func RemoveAt[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// Insert inserts an element at the specified index
func Insert[T any](slice []T, index int, item T) []T {
	if index < 0 {
		index = 0
	}
	if index > len(slice) {
		index = len(slice)
	}

	slice = append(slice, item)
	copy(slice[index+1:], slice[index:])
	slice[index] = item
	return slice
}

// Prepend adds an element to the beginning
func Prepend[T any](slice []T, item T) []T {
	return append([]T{item}, slice...)
}

// Append adds an element to the end (alias for built-in append)
func Append[T any](slice []T, item T) []T {
	return append(slice, item)
}

// Concat concatenates slices
func Concat[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Repeat repeats an element n times
func Repeat[T any](item T, count int) []T {
	if count <= 0 {
		return []T{}
	}

	result := make([]T, count)
	for i := range result {
		result[i] = item
	}
	return result
}

// Range creates a slice of numbers from start to end (not including end)
func Range(start, end int) []int {
	if start >= end {
		return []int{}
	}

	result := make([]int, end-start)
	for i := range result {
		result[i] = start + i
	}
	return result
}

// RangeStep creates a slice of numbers with a step
func RangeStep(start, end, step int) []int {
	if step == 0 || (step > 0 && start >= end) || (step < 0 && start <= end) {
		return []int{}
	}

	var result []int
	if step > 0 {
		for i := start; i < end; i += step {
			result = append(result, i)
		}
	} else {
		for i := start; i > end; i += step {
			result = append(result, i)
		}
	}
	return result
}

// SortBy sorts a slice by a key extraction function
func SortBy[T any, K any](slice []T, keyExtractor func(T) K, less func(K, K) bool) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	sort.Slice(result, func(i, j int) bool {
		return less(keyExtractor(result[i]), keyExtractor(result[j]))
	})

	return result
}

// SortByComparable sorts a slice by a comparable key
func SortByComparable[T any, K comparable](slice []T, keyExtractor func(T) K) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	sort.Slice(result, func(i, j int) bool {
		ki, kj := keyExtractor(result[i]), keyExtractor(result[j])
		return fmt.Sprintf("%v", ki) < fmt.Sprintf("%v", kj)
	})

	return result
}

// Shuffle shuffles a slice (simple implementation)
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	// Simple Fisher-Yates implementation without crypto/rand
	for i := len(result) - 1; i > 0; i-- {
		j := (i*7 + 13) % (i + 1) // Simple pseudorandomness
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// Sample returns a random element from the slice
func Sample[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	// Simple pseudorandomness
	index := (len(slice)*17 + 23) % len(slice)
	return slice[index], true
}

// SampleN returns n random elements
func SampleN[T any](slice []T, n int) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}

	if n >= len(slice) {
		return Shuffle(slice)
	}

	shuffled := Shuffle(slice)
	return shuffled[:n]
}

// ToMap converts a slice to a map
func ToMap[T any, K comparable, V any](slice []T, keyExtractor func(T) K, valueExtractor func(T) V) map[K]V {
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		key := keyExtractor(item)
		value := valueExtractor(item)
		result[key] = value
	}
	return result
}

// ToMapBy converts a slice to a map with elements as values
func ToMapBy[T any, K comparable](slice []T, keyExtractor func(T) K) map[K]T {
	return ToMap(slice, keyExtractor, Identity[T])
}

// FromMap creates a slice from a map
func FromMap[K comparable, V any, T any](m map[K]V, transformer func(K, V) T) []T {
	result := make([]T, 0, len(m))
	for k, v := range m {
		result = append(result, transformer(k, v))
	}
	return result
}

// Keys returns the keys of a map
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values returns the values of a map
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Entries returns the key-value pairs of a map
func Entries[K comparable, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{First: k, Second: v})
	}
	return result
}

// ToString converts a value to a string
func ToString[T any](value T) string {
	return fmt.Sprintf("%v", value)
}

// ToJSON converts a value to a JSON string
func ToJSON[T any](value T) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON parses a JSON string into a value
func FromJSON[T any](jsonStr string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// ToInt converts a string to an int
func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ToFloat converts a string to a float64
func ToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ToBool converts a string to a bool
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// IsEmpty checks if a slice is empty
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// IsNotEmpty checks if a slice is not empty
func IsNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}

// Size returns the size of a slice
func Size[T any](slice []T) int {
	return len(slice)
}

// Head returns the first element of a slice
func Head[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Tail returns all elements except the first
func Tail[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	result := make([]T, len(slice)-1)
	copy(result, slice[1:])
	return result
}

// Last returns the last element of a slice
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Init returns all elements except the last
func Init[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	result := make([]T, len(slice)-1)
	copy(result, slice[:len(slice)-1])
	return result
}

// IsEqual checks if two slices are equal
func IsEqual[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}
	return true
}

// IsEqualBy checks if two slices are equal by a comparison function
func IsEqualBy[T any](slice1, slice2 []T, equals Equality[T]) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, v := range slice1 {
		if !equals(v, slice2[i]) {
			return false
		}
	}
	return true
}

// DeepEqual checks if two values are equal using reflection
func DeepEqual[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// Clone creates a shallow copy of a slice
func Clone[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

// Words splits a string into words
func Words(s string) []string {
	return strings.Fields(s)
}

// Lines splits a string into lines
func Lines(s string) []string {
	return strings.Split(s, "\n")
}

// Chars splits a string into characters
func Chars(s string) []string {
	return strings.Split(s, "")
}

// Runes splits a string into runes
func Runes(s string) []rune {
	return []rune(s)
	// TODO: implement Runes
}

// StringJoin joins a slice of strings with a separator
func StringJoin(strs []string, separator string) string {
	return strings.Join(strs, separator)
}

// Conditional returns one of two values based on a condition
func Conditional[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// ConditionalFunc returns the result of one of two functions based on a condition
func ConditionalFunc[T any](condition bool, ifTrue, ifFalse func() T) T {
	if condition {
		return ifTrue()
	}
	return ifFalse()
}
