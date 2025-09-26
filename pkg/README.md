# Functional Programming Library for Go

Мощная библиотека функционального программирования для Go с современными дженериками.

## Основные возможности

### 🚀 Базовые функции

- **Map, Filter, Reduce** - классические функции высшего порядка
- **Параллельная обработка** - автоматическое распараллеливание для больших коллекций
- **Композиция функций** - Pipe, Compose, Curry

### 🔧 Продвинутые утилиты

- **Optional/Result типы** - безопасная работа с nullable значениями
- **Коллекции** - GroupBy, Chunk, Partition, Zip, Flatten
- **Параллельные операции** - с контекстом и rate limiting

## Примеры использования

### Базовые операции

```go
import "github.com/malinatrash/funky/pkg/fp"

// Map - преобразование
numbers := []int{1, 2, 3, 4, 5}
doubled := fp.Map(numbers, func(x int) int { return x * 2 })
// [2, 4, 6, 8, 10]

// Filter - фильтрация
evens := fp.Filter(numbers, func(x int) bool { return x%2 == 0 })
// [2, 4]

// Reduce - свертка
sum := fp.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
// 15
```

### Pipeline обработки

```go
result := fp.NewPipeline([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Collect()
// [4, 16, 36, 64, 100]
```

### Optional/Result типы

```go
// Optional для безопасной работы с nil
user := fp.Some("John")
name := user.Map(strings.ToUpper).GetOrElse("Unknown")
// "JOHN"

// Result для обработки ошибок
result := fp.Try(func() (int, error) {
    return strconv.Atoi("42")
}).Map(func(x int) int { return x * 2 })

if result.IsOk() {
    fmt.Println(result.Unwrap()) // 84
}
```

### Композиция функций

```go
// Pipe - слева направо
result := fp.Pipe3("hello",
    strings.ToUpper,
    func(s string) string { return s + "!" },
    func(s string) string { return ">>> " + s })
// ">>> HELLO!"

// Compose - справа налево
transform := fp.Compose3(
    func(s string) string { return ">>> " + s },
    func(s string) string { return s + "!" },
    strings.ToUpper)
result := transform("hello")
// ">>> HELLO!"
```

### Работа с коллекциями

```go
// GroupBy
users := []User{{Name: "John", Age: 25}, {Name: "Jane", Age: 25}}
byAge := fp.GroupBy(users, func(u User) int { return u.Age })

// Chunk
numbers := []int{1, 2, 3, 4, 5, 6, 7}
chunks := fp.Chunk(numbers, 3)
// [[1, 2, 3], [4, 5, 6], [7]]

// Zip
names := []string{"John", "Jane"}
ages := []int{25, 30}
pairs := fp.Zip(names, ages)
// [{John, 25}, {Jane, 30}]
```

### Параллельная обработка

```go
// Параллельный Map
large := make([]int, 10000)
result := fp.MapParallel(large, heavyComputation)

// С контекстом
ctx := context.Background()
result, err := fp.MapWithContext(ctx, data, func(ctx context.Context, item int) (string, error) {
    return processItem(ctx, item)
}, fp.DefaultParallelConfig())
```

## Структура библиотеки

- `commonconst.go` - Общие типы и константы
- `map.go` - Функции преобразования
- `filter.go` - Функции фильтрации
- `reduce.go` - Функции свертки и агрегации
- `compose.go` - Композиция и каррирование функций
- `collections.go` - Утилиты для работы с коллекциями
- `optional.go` - Optional и Result типы
- `parallel.go` - Параллельная обработка
- `utils.go` - Дополнительные утилиты

## Производительность

Библиотека автоматически выбирает оптимальную стратегию:

- Для маленьких коллекций (<100 элементов) - последовательная обработка
- Для больших коллекций - параллельная обработка с воркер-пулом
- Настраиваемые параметры параллелизма

## Совместимость

- Go 1.18+ (требуются дженерики)
- Thread-safe операции
- Zero-dependency (только стандартная библиотека)
