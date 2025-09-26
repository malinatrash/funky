package fp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Try выполняет функцию и возвращает Result
func Try[T any](fn func() (T, error)) Result[T] {
	return TryFrom(fn)
}

// TryVoid выполняет функцию без возвращаемого значения
func TryVoid(fn func() error) Result[struct{}] {
	err := fn()
	if err != nil {
		return Err[struct{}](err)
	}
	return Ok(struct{}{})
}

// Tap выполняет побочный эффект и возвращает исходное значение
func Tap[T any](value T, sideEffect func(T)) T {
	sideEffect(value)
	return value
}

// TapIf выполняет побочный эффект только если условие истинно
func TapIf[T any](value T, condition bool, sideEffect func(T)) T {
	if condition {
		sideEffect(value)
	}
	return value
}

// TapWhen выполняет побочный эффект только если предикат истинен
func TapWhen[T any](value T, predicate Predicate[T], sideEffect func(T)) T {
	if predicate(value) {
		sideEffect(value)
	}
	return value
}

// Let позволяет выполнить операцию над значением
func Let[T, R any](value T, transform func(T) R) R {
	return transform(value)
}

// Also выполняет операцию и возвращает исходное значение
func Also[T any](value T, operation func(T)) T {
	operation(value)
	return value
}

// TakeIf возвращает значение если предикат истинен, иначе nil
func TakeIf[T any](value T, predicate Predicate[T]) *T {
	if predicate(value) {
		return &value
	}
	return nil
}

// TakeUnless возвращает значение если предикат ложен, иначе nil
func TakeUnless[T any](value T, predicate Predicate[T]) *T {
	if !predicate(value) {
		return &value
	}
	return nil
}

// Coalesce возвращает первое не-nil значение
func Coalesce[T any](values ...*T) *T {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// CoalesceFunc возвращает первое не-nil значение из функций
func CoalesceFunc[T any](suppliers ...func() *T) *T {
	for _, supplier := range suppliers {
		if value := supplier(); value != nil {
			return value
		}
	}
	return nil
}

// DefaultIfNil возвращает значение по умолчанию если указатель nil
func DefaultIfNil[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// Ptr создает указатель на значение
func Ptr[T any](value T) *T {
	return &value
}

// Deref разыменовывает указатель с значением по умолчанию
func Deref[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// SafeDeref безопасно разыменовывает указатель
func SafeDeref[T any](ptr *T) (T, bool) {
	if ptr == nil {
		var zero T
		return zero, false
	}
	return *ptr, true
}

// Contains проверяет наличие элемента в слайсе
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsBy проверяет наличие элемента по предикату
func ContainsBy[T any](slice []T, predicate Predicate[T]) bool {
	return Any(slice, predicate)
}

// IndexOf возвращает индекс элемента в слайсе
func IndexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// LastIndexOf возвращает последний индекс элемента в слайсе
func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// Remove удаляет первое вхождение элемента
func Remove[T comparable](slice []T, item T) []T {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// RemoveAll удаляет все вхождения элемента
func RemoveAll[T comparable](slice []T, item T) []T {
	return Filter(slice, func(v T) bool { return v != item })
}

// RemoveAt удаляет элемент по индексу
func RemoveAt[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// Insert вставляет элемент по индексу
func Insert[T any](slice []T, index int, item T) []T {
	if index < 0 {
		index = 0
	}
	if index > len(slice) {
		index = len(slice)
	}
	
	slice = append(slice, item) // увеличиваем размер
	copy(slice[index+1:], slice[index:])
	slice[index] = item
	return slice
}

// Prepend добавляет элемент в начало
func Prepend[T any](slice []T, item T) []T {
	return append([]T{item}, slice...)
}

// Append добавляет элемент в конец (алиас для встроенного append)
func Append[T any](slice []T, item T) []T {
	return append(slice, item)
}

// Concat объединяет слайсы
func Concat[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Repeat повторяет элемент n раз
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

// Range создает слайс чисел от start до end (не включая end)
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

// RangeStep создает слайс чисел с шагом
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

// SortBy сортирует слайс по функции извлечения ключа
func SortBy[T any, K any](slice []T, keyExtractor func(T) K, less func(K, K) bool) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	
	sort.Slice(result, func(i, j int) bool {
		return less(keyExtractor(result[i]), keyExtractor(result[j]))
	})
	
	return result
}

// SortByComparable сортирует слайс по comparable ключу
func SortByComparable[T any, K comparable](slice []T, keyExtractor func(T) K) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	
	sort.Slice(result, func(i, j int) bool {
		ki, kj := keyExtractor(result[i]), keyExtractor(result[j])
		return fmt.Sprintf("%v", ki) < fmt.Sprintf("%v", kj)
	})
	
	return result
}

// Shuffle перемешивает слайс (простая реализация)
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	
	// Простая реализация Fisher-Yates без crypto/rand
	for i := len(result) - 1; i > 0; i-- {
		j := (i * 7 + 13) % (i + 1) // Простая псевдослучайность
		result[i], result[j] = result[j], result[i]
	}
	
	return result
}

// Sample возвращает случайный элемент из слайса
func Sample[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	
	// Простая псевдослучайность
	index := (len(slice) * 17 + 23) % len(slice)
	return slice[index], true
}

// SampleN возвращает n случайных элементов
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

// ToMap преобразует слайс в мапу
func ToMap[T any, K comparable, V any](slice []T, keyExtractor func(T) K, valueExtractor func(T) V) map[K]V {
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		key := keyExtractor(item)
		value := valueExtractor(item)
		result[key] = value
	}
	return result
}

// ToMapBy преобразует слайс в мапу с элементами как значениями
func ToMapBy[T any, K comparable](slice []T, keyExtractor func(T) K) map[K]T {
	return ToMap(slice, keyExtractor, Identity[T])
}

// FromMap создает слайс из мапы
func FromMap[K comparable, V any, T any](m map[K]V, transformer func(K, V) T) []T {
	result := make([]T, 0, len(m))
	for k, v := range m {
		result = append(result, transformer(k, v))
	}
	return result
}

// Keys возвращает ключи мапы
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values возвращает значения мапы
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Entries возвращает пары ключ-значение
func Entries[K comparable, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{First: k, Second: v})
	}
	return result
}

// ToString преобразует значение в строку
func ToString[T any](value T) string {
	return fmt.Sprintf("%v", value)
}

// ToJSON преобразует значение в JSON строку
func ToJSON[T any](value T) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON парсит JSON строку в значение
func FromJSON[T any](jsonStr string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// ToInt преобразует строку в int
func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ToFloat преобразует строку в float64
func ToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ToBool преобразует строку в bool
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// IsEmpty проверяет, пуст ли слайс
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// IsNotEmpty проверяет, не пуст ли слайс
func IsNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}

// Size возвращает размер слайса
func Size[T any](slice []T) int {
	return len(slice)
}

// Head возвращает первый элемент слайса
func Head[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Tail возвращает все элементы кроме первого
func Tail[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	result := make([]T, len(slice)-1)
	copy(result, slice[1:])
	return result
}

// Last возвращает последний элемент слайса
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Init возвращает все элементы кроме последнего
func Init[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	result := make([]T, len(slice)-1)
	copy(result, slice[:len(slice)-1])
	return result
}

// IsEqual проверяет равенство двух слайсов
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

// IsEqualBy проверяет равенство двух слайсов по функции сравнения
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

// DeepEqual проверяет глубокое равенство с помощью reflection
func DeepEqual[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// Clone создает поверхностную копию слайса
func Clone[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

// Words разбивает строку на слова
func Words(s string) []string {
	return strings.Fields(s)
}

// Lines разбивает строку на строки
func Lines(s string) []string {
	return strings.Split(s, "\n")
}

// Chars разбивает строку на символы
func Chars(s string) []string {
	return strings.Split(s, "")
}

// Runes разбивает строку на руны
func Runes(s string) []rune {
	return []rune(s)
}

// StringJoin объединяет строки
func StringJoin(strings []string, separator string) string {
	return strings.Join(strings, separator)
}

// Conditional возвращает одно из двух значений в зависимости от условия
func Conditional[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// ConditionalFunc возвращает результат одной из функций
func ConditionalFunc[T any](condition bool, ifTrue, ifFalse func() T) T {
	if condition {
		return ifTrue()
	}
	return ifFalse()
}
