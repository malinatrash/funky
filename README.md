# Functional Programming Library for Go

Powerful library for functional programming in Go with modern generics.

## Main features

### 🚀 Base functions

- **Map, Filter, Reduce** - classic functions of higher order
- **Parallel processing** - automatic parallelization for large collections
- **Function composition** - Pipe, Compose, Curry

### 🔧 Advanced utilities

- **Optional/Result types** - safe work with nullable values
- **Collections** - GroupBy, Chunk, Partition, Zip, Flatten
- **Parallel operations** - with context and rate limiting

## Examples

### Base operations

```go
import "github.com/malinatrash/funky/pkg/fp"

// Map - transformation
numbers := []int{1, 2, 3, 4, 5}
doubled := fp.Map(numbers, func(x int) int { return x * 2 })
// [2, 4, 6, 8, 10]

// Filter - filtering
evens := fp.Filter(numbers, func(x int) bool { return x%2 == 0 })
// [2, 4]

// Reduce - reduction
sum := fp.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
// 15
```

### Pipeline processing

```go
result := fp.NewPipeline([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Collect()
// [4, 16, 36, 64, 100]
```

### Optional/Result types

```go
// Optional for safe work with nil
user := fp.Some("John")
name := user.Map(strings.ToUpper).GetOrElse("Unknown")
// "JOHN"

// Result for error handling
result := fp.Try(func() (int, error) {
    return strconv.Atoi("42")
}).Map(func(x int) int { return x * 2 })

if result.IsOk() {
    fmt.Println(result.Unwrap()) // 84
}
```

### Function composition

```go
// Pipe - from left to right
result := fp.Pipe3("hello",
    strings.ToUpper,
    func(s string) string { return s + "!" },
    func(s string) string { return ">>> " + s })
// ">>> HELLO!"

// Compose - from right to left
transform := fp.Compose3(
    func(s string) string { return ">>> " + s },
    func(s string) string { return s + "!" },
    strings.ToUpper)
result := transform("hello")
// ">>> HELLO!"
```

### Collections

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

### Parallel processing

```go
// Parallel Map
large := make([]int, 10000)
result := fp.MapParallel(large, heavyComputation)

// With context
ctx := context.Background()
result, err := fp.MapWithContext(ctx, data, func(ctx context.Context, item int) (string, error) {
    return processItem(ctx, item)
}, fp.DefaultParallelConfig())
```

## Library structure

- `commonconst.go` - Common types and constants
- `map.go` - Mapping functions
- `filter.go` - Filtering functions
- `reduce.go` - Reduction functions
- `compose.go` - Function composition and currying
- `collections.go` - Collection utilities
- `optional.go` - Optional and Result types
- `parallel.go` - Parallel processing
- `utils.go` - Additional utilities

## Performance

The library automatically selects the optimal strategy:

- For small collections (<100 elements) - sequential processing
- For large collections - parallel processing with worker pool
- Configurable parallelism parameters

## Compatibility

- Go 1.18+ (requires generics)
- Thread-safe operations
- Zero-dependency (only standard library)
