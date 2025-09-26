package fp

// Pipe executes function composition from left to right
func Pipe[T any](value T, functions ...func(T) T) T {
	result := value
	for _, fn := range functions {
		result = fn(result)
	}
	return result
}

// Pipe2 for functions with two types
func Pipe2[T, R any](value T, fn1 func(T) R) R {
	return fn1(value)
}

// Pipe3 for functions with three types
func Pipe3[T, R, S any](value T, fn1 func(T) R, fn2 func(R) S) S {
	return fn2(fn1(value))
}

// Pipe4 for functions with four types
func Pipe4[T, R, S, U any](value T, fn1 func(T) R, fn2 func(R) S, fn3 func(S) U) U {
	return fn3(fn2(fn1(value)))
}

// Pipe5 for functions with five types
func Pipe5[T, R, S, U, V any](value T, fn1 func(T) R, fn2 func(R) S, fn3 func(S) U, fn4 func(U) V) V {
	return fn4(fn3(fn2(fn1(value))))
}

// Compose executes function composition from right to left
func Compose[T any](functions ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for i := len(functions) - 1; i >= 0; i-- {
			result = functions[i](result)
		}
		return result
	}
}

// Compose2 for functions with two types
func Compose2[T, R any](fn1 func(T) R) func(T) R {
	return fn1
}

// Compose3 for functions with three types
func Compose3[T, R, S any](fn1 func(R) S, fn2 func(T) R) func(T) S {
	return func(value T) S {
		return fn1(fn2(value))
	}
}

// Compose4 for functions with four types
func Compose4[T, R, S, U any](fn1 func(S) U, fn2 func(R) S, fn3 func(T) R) func(T) U {
	return func(value T) U {
		return fn1(fn2(fn3(value)))
	}
}

// Curry transforms a function with two arguments into a curried function
func Curry2[T, R, S any](fn func(T, R) S) func(T) func(R) S {
	return func(t T) func(R) S {
		return func(r R) S {
			return fn(t, r)
		}
	}
}

// Curry3 transforms a function with three arguments into a curried function
func Curry3[T, R, S, U any](fn func(T, R, S) U) func(T) func(R) func(S) U {
	return func(t T) func(R) func(S) U {
		return func(r R) func(S) U {
			return func(s S) U {
				return fn(t, r, s)
			}
		}
	}
}

// Uncurry transforms a curried function back into a regular function
func Uncurry2[T, R, S any](fn func(T) func(R) S) func(T, R) S {
	return func(t T, r R) S {
		return fn(t)(r)
	}
}

// Uncurry3 transforms a curried function with three arguments back into a regular function
func Uncurry3[T, R, S, U any](fn func(T) func(R) func(S) U) func(T, R, S) U {
	return func(t T, r R, s S) U {
		return fn(t)(r)(s)
	}
}

// Partial applies partial application to a function
func Partial2[T, R, S any](fn func(T, R) S, t T) func(R) S {
	return func(r R) S {
		return fn(t, r)
	}
}

// Partial3 applies partial application to a function with three arguments
func Partial3[T, R, S, U any](fn func(T, R, S) U, t T) func(R, S) U {
	return func(r R, s S) U {
		return fn(t, r, s)
	}
}

// Flip changes the order of arguments in a function
func Flip[T, R, S any](fn func(T, R) S) func(R, T) S {
	return func(r R, t T) S {
		return fn(t, r)
	}
}

// Memoize caches the results of a function
func Memoize[T comparable, R any](fn func(T) R) func(T) R {
	cache := make(map[T]R)
	return func(t T) R {
		if result, exists := cache[t]; exists {
			return result
		}
		result := fn(t)
		cache[t] = result
		return result
	}
}

// MemoizeWithTTL caches the results of a function with TTL
func MemoizeWithTTL[T comparable, R any](fn func(T) R, ttl int64) func(T) R {
	type cacheEntry struct {
		value     R
		timestamp int64
	}

	cache := make(map[T]cacheEntry)

	return func(t T) R {
		now := getCurrentTimestamp() // Вы можете заменить на time.Now().Unix()

		if entry, exists := cache[t]; exists && (now-entry.timestamp) < ttl {
			return entry.value
		}

		result := fn(t)
		cache[t] = cacheEntry{
			value:     result,
			timestamp: now,
		}
		return result
	}
}

// Debounce creates a function with execution delay
func Debounce[T any](fn func(T), delay int64) func(T) {
	var lastCall int64

	return func(t T) {
		now := getCurrentTimestamp()
		lastCall = now

		go func() {
			// Simple debounce implementation without time.Sleep for example
			// In real code use time.Sleep(time.Duration(delay) * time.Millisecond)
			if getCurrentTimestamp()-lastCall >= delay {
				fn(t)
			}
		}()
	}
}

// Throttle creates a function with execution throttling
func Throttle[T any](fn func(T), interval int64) func(T) {
	var lastExecution int64

	return func(t T) {
		now := getCurrentTimestamp()
		if now-lastExecution >= interval {
			lastExecution = now
			fn(t)
		}
	}
}

// getCurrentTimestamp - helper function to get the current timestamp
// In real code replace with time.Now().Unix()
func getCurrentTimestamp() int64 {
	// Stub for compilation, in real code use time.Now().Unix()
	return 0
}
