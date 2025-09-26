package fp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ExampleUser для демонстрации
type ExampleUser struct {
	ID   int
	Name string
	Age  int
	City string
}

// ExampleOrder для демонстрации
type ExampleOrder struct {
	ID       int
	UserID   int
	Amount   float64
	Status   string
	Products []string
}

// DemoBasicOperations демонстрирует базовые операции
func DemoBasicOperations() {
	fmt.Println("=== Базовые операции ===")
	
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	// Map - удваиваем числа
	doubled := Map(numbers, func(x int) int { return x * 2 })
	fmt.Printf("Doubled: %v\n", doubled)
	
	// Filter - только четные
	evens := Filter(numbers, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Evens: %v\n", evens)
	
	// Reduce - сумма
	sum := Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
	fmt.Printf("Sum: %d\n", sum)
	
	// Комбинированная операция
	result := Pipe3(numbers,
		func(nums []int) []int { return Filter(nums, func(x int) bool { return x > 5 }) },
		func(nums []int) []int { return Map(nums, func(x int) int { return x * x }) },
		func(nums []int) int { return Reduce(nums, func(acc, x int) int { return acc + x }, 0) })
	
	fmt.Printf("Filtered > 5, squared and summed: %d\n", result)
}

// DemoCollections демонстрирует работу с коллекциями
func DemoCollections() {
	fmt.Println("\n=== Работа с коллекциями ===")
	
	users := []ExampleUser{
		{1, "Alice", 25, "Moscow"},
		{2, "Bob", 30, "SPB"},
		{3, "Charlie", 25, "Moscow"},
		{4, "Diana", 35, "Kazan"},
		{5, "Eve", 30, "Moscow"},
	}
	
	// GroupBy - группировка по возрасту
	byAge := GroupBy(users, func(u ExampleUser) int { return u.Age })
	fmt.Printf("Grouped by age: %v\n", Keys(byAge))
	
	// GroupBy - группировка по городу
	byCity := GroupBy(users, func(u ExampleUser) string { return u.City })
	for city, cityUsers := range byCity {
		names := Map(cityUsers, func(u ExampleUser) string { return u.Name })
		fmt.Printf("City %s: %v\n", city, names)
	}
	
	// Partition - разделение на молодых и старых
	young, old := Partition(users, func(u ExampleUser) bool { return u.Age < 30 })
	fmt.Printf("Young: %v\n", Map(young, func(u ExampleUser) string { return u.Name }))
	fmt.Printf("Old: %v\n", Map(old, func(u ExampleUser) string { return u.Name }))
	
	// Chunk - разбиение на группы по 2
	chunks := Chunk(users, 2)
	fmt.Printf("Chunks of 2: %d groups\n", len(chunks))
	
	// Unique names
	names := Map(users, func(u ExampleUser) string { return u.Name })
	cities := Map(users, func(u ExampleUser) string { return u.City })
	uniqueCities := Unique(cities)
	fmt.Printf("All names: %v\n", names)
	fmt.Printf("Unique cities: %v\n", uniqueCities)
}

// DemoOptionalAndResult демонстрирует Optional и Result
func DemoOptionalAndResult() {
	fmt.Println("\n=== Optional и Result ===")
	
	// Optional
	user := Some("John Doe")
	upperName := user.Map(strings.ToUpper)
	fmt.Printf("Upper name: %s\n", upperName.GetOrElse("Unknown"))
	
	emptyUser := Empty[string]()
	defaultName := emptyUser.GetOrElse("Default User")
	fmt.Printf("Default name: %s\n", defaultName)
	
	// Result
	parseNumber := func(s string) Result[int] {
		return Try(func() (int, error) {
			return strconv.Atoi(s)
		})
	}
	
	validResult := parseNumber("42")
	if validResult.IsOk() {
		doubled := validResult.Map(func(x int) int { return x * 2 })
		fmt.Printf("Valid number doubled: %d\n", doubled.Unwrap())
	}
	
	invalidResult := parseNumber("not-a-number")
	if invalidResult.IsErr() {
		fmt.Printf("Parse error: %v\n", invalidResult.Error())
	}
	
	// Цепочка операций с Result
	result := parseNumber("10").
		Map(func(x int) int { return x * 2 }).
		Map(func(x int) int { return x + 5 })
	
	fmt.Printf("Chain result: %d\n", result.UnwrapOr(-1))
}

// DemoStreams демонстрирует работу с потоками
func DemoStreams() {
	fmt.Println("\n=== Потоки данных ===")
	
	// Создание потока и ленивые вычисления
	numbers := RangeStream(1, 11).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Collect()
	
	fmt.Printf("Stream result: %v\n", numbers)
	
	// Бесконечный поток (ограничиваем Take)
	fibonacci := func() *Stream[int] {
		return NewStreamFromFunc(func() <-chan int {
			ch := make(chan int)
			go func() {
				defer close(ch)
				a, b := 0, 1
				for {
					ch <- a
					a, b = b, a+b
				}
			}()
			return ch
		})
	}
	
	fibNumbers := fibonacci().Take(10).Collect()
	fmt.Printf("Fibonacci (10): %v\n", fibNumbers)
	
	// Параллельная обработка в потоке
	heavyComputation := func(x int) int {
		// Имитация тяжелых вычислений
		time.Sleep(1 * time.Millisecond)
		return x * x
	}
	
	start := time.Now()
	parallelResult := RangeStream(1, 101).
		Parallel(4, heavyComputation).
		Take(10).
		Collect()
	duration := time.Since(start)
	
	fmt.Printf("Parallel stream (first 10): %v\n", parallelResult[:min(len(parallelResult), 5)])
	fmt.Printf("Parallel processing took: %v\n", duration)
}

// DemoAdvancedFeatures демонстрирует продвинутые возможности
func DemoAdvancedFeatures() {
	fmt.Println("\n=== Продвинутые возможности ===")
	
	// Композиция функций
	processString := Compose3(
		func(s string) string { return "Result: " + s },
		func(s string) string { return strings.ToUpper(s) },
		strings.TrimSpace,
	)
	
	result := processString("  hello world  ")
	fmt.Printf("Composed function: %s\n", result)
	
	// Каррирование
	add := func(a, b int) int { return a + b }
	curriedAdd := Curry2(add)
	add5 := curriedAdd(5)
	
	fmt.Printf("Curried add(5)(10): %d\n", add5(10))
	
	// Мемоизация
	expensiveFunc := func(n int) int {
		fmt.Printf("Computing for %d...\n", n)
		time.Sleep(10 * time.Millisecond)
		return n * n
	}
	
	memoized := Memoize(expensiveFunc)
	
	fmt.Printf("First call: %d\n", memoized(5))  // Вычисляется
	fmt.Printf("Second call: %d\n", memoized(5)) // Из кэша
	
	// Работа с контекстом
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	data := []int{1, 2, 3, 4, 5}
	slowProcessor := func(ctx context.Context, x int) (string, error) {
		select {
		case <-time.After(50 * time.Millisecond):
			return fmt.Sprintf("processed-%d", x), nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	
	processed, err := MapWithContext(ctx, data, slowProcessor, DefaultParallelConfig())
	if err != nil {
		fmt.Printf("Context processing error: %v\n", err)
	} else {
		fmt.Printf("Context processing result: %v\n", processed)
	}
}

// DemoRealWorldExample демонстрирует реальный пример использования
func DemoRealWorldExample() {
	fmt.Println("\n=== Реальный пример: Обработка заказов ===")
	
	orders := []ExampleOrder{
		{1, 1, 100.50, "completed", []string{"laptop", "mouse"}},
		{2, 2, 50.25, "pending", []string{"book"}},
		{3, 1, 200.00, "completed", []string{"monitor", "keyboard"}},
		{4, 3, 75.75, "cancelled", []string{"headphones"}},
		{5, 2, 150.00, "completed", []string{"tablet"}},
	}
	
	users := []ExampleUser{
		{1, "Alice", 25, "Moscow"},
		{2, "Bob", 30, "SPB"},
		{3, "Charlie", 25, "Moscow"},
	}
	
	// Создаем мапу пользователей для быстрого поиска
	userMap := ToMapBy(users, func(u ExampleUser) int { return u.ID })
	
	// Анализ заказов
	completedOrders := Filter(orders, func(o ExampleOrder) bool { return o.Status == "completed" })
	
	totalRevenue := Reduce(completedOrders, func(acc float64, o ExampleOrder) float64 {
		return acc + o.Amount
	}, 0.0)
	
	// Группировка по пользователям
	ordersByUser := GroupBy(completedOrders, func(o ExampleOrder) int { return o.UserID })
	
	// Статистика по пользователям
	userStats := Map(Keys(ordersByUser), func(userID int) map[string]interface{} {
		userOrders := ordersByUser[userID]
		user := userMap[userID]
		
		totalAmount := Reduce(userOrders, func(acc float64, o ExampleOrder) float64 {
			return acc + o.Amount
		}, 0.0)
		
		allProducts := FlatMap(userOrders, func(o ExampleOrder) []string {
			return o.Products
		})
		
		return map[string]interface{}{
			"user":         user.Name,
			"city":         user.City,
			"orderCount":   len(userOrders),
			"totalAmount":  totalAmount,
			"avgAmount":    totalAmount / float64(len(userOrders)),
			"products":     Unique(allProducts),
			"productCount": len(Unique(allProducts)),
		}
	})
	
	fmt.Printf("Total completed orders: %d\n", len(completedOrders))
	fmt.Printf("Total revenue: $%.2f\n", totalRevenue)
	fmt.Printf("Average order value: $%.2f\n", totalRevenue/float64(len(completedOrders)))
	
	fmt.Println("\nUser Statistics:")
	for _, stats := range userStats {
		fmt.Printf("- %s (%s): %d orders, $%.2f total, $%.2f avg, %d unique products\n",
			stats["user"], stats["city"], stats["orderCount"],
			stats["totalAmount"], stats["avgAmount"], stats["productCount"])
	}
	
	// Топ продукты
	allProducts := FlatMap(completedOrders, func(o ExampleOrder) []string {
		return o.Products
	})
	
	productCounts := GroupByCount(allProducts, Identity[string])
	topProducts := FromMap(productCounts, func(product string, count int) Pair[string, int] {
		return Pair[string, int]{First: product, Second: count}
	})
	
	// Сортируем по количеству (упрощенная сортировка)
	fmt.Println("\nTop Products:")
	for _, pair := range topProducts {
		fmt.Printf("- %s: %d orders\n", pair.First, pair.Second)
	}
}

// RunAllExamples запускает все примеры
func RunAllExamples() {
	fmt.Println("🚀 Демонстрация Functional Programming Library")
	fmt.Println("=" * 50)
	
	DemoBasicOperations()
	DemoCollections()
	DemoOptionalAndResult()
	DemoStreams()
	DemoAdvancedFeatures()
	DemoRealWorldExample()
	
	fmt.Println("\n✅ Все примеры выполнены успешно!")
}
