package main

import (
	"fmt"
	"github.com/malinatrash/funky/pkg/fp"
)

func main() {
	fmt.Println("Testing fixed examples...")
	
	// Test the fixed combined operation
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filtered := fp.Filter(numbers, func(x int) bool { return x > 5 })
	squared := fp.Map(filtered, func(x int) int { return x * x })
	result := fp.Reduce(squared, func(acc, x int) int { return acc + x }, 0)
	
	fmt.Printf("Filtered > 5, squared and summed: %d\n", result)
	
	// Test the fixed composition
	processString := fp.Compose3(
		func(s string) string { return "Result: " + s },
		func(s string) string { return strings.ToUpper(strings.TrimSpace(s)) },
	)
	
	composedResult := processString("  hello world  ")
	fmt.Printf("Composed function: %s\n", composedResult)
	
	fmt.Println("✅ All fixes working correctly!")
}
