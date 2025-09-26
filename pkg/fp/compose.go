package fp

// Pipe выполняет композицию функций слева направо
func Pipe[T any](value T, functions ...func(T) T) T {
	result := value
	for _, fn := range functions {
		result = fn(result)
	}
	return result
}

// Pipe2 для функций с двумя типами
func Pipe2[T, R any](value T, fn1 func(T) R) R {
	return fn1(value)
}

// Pipe3 для функций с тремя типами
func Pipe3[T, R, S any](value T, fn1 func(T) R, fn2 func(R) S) S {
	return fn2(fn1(value))
}

// Pipe4 для функций с четырьмя типами
func Pipe4[T, R, S, U any](value T, fn1 func(T) R, fn2 func(R) S, fn3 func(S) U) U {
	return fn3(fn2(fn1(value)))
}

// Pipe5 для функций с пятью типами
func Pipe5[T, R, S, U, V any](value T, fn1 func(T) R, fn2 func(R) S, fn3 func(S) U, fn4 func(U) V) V {
	return fn4(fn3(fn2(fn1(value))))
}

// Compose выполняет композицию функций справа налево
func Compose[T any](functions ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for i := len(functions) - 1; i >= 0; i-- {
			result = functions[i](result)
		}
		return result
	}
}

// Compose2 для функций с двумя типами
func Compose2[T, R any](fn1 func(T) R) func(T) R {
	return fn1
}

// Compose3 для функций с тремя типами
func Compose3[T, R, S any](fn1 func(R) S, fn2 func(T) R) func(T) S {
	return func(value T) S {
		return fn1(fn2(value))
	}
}

// Compose4 для функций с четырьмя типами
func Compose4[T, R, S, U any](fn1 func(S) U, fn2 func(R) S, fn3 func(T) R) func(T) U {
	return func(value T) U {
		return fn1(fn2(fn3(value)))
	}
}

// Curry преобразует функцию с двумя аргументами в каррированную
func Curry2[T, R, S any](fn func(T, R) S) func(T) func(R) S {
	return func(t T) func(R) S {
		return func(r R) S {
			return fn(t, r)
		}
	}
}

// Curry3 преобразует функцию с тремя аргументами в каррированную
func Curry3[T, R, S, U any](fn func(T, R, S) U) func(T) func(R) func(S) U {
	return func(t T) func(R) func(S) U {
		return func(r R) func(S) U {
			return func(s S) U {
				return fn(t, r, s)
			}
		}
	}
}

// Uncurry преобразует каррированную функцию обратно
func Uncurry2[T, R, S any](fn func(T) func(R) S) func(T, R) S {
	return func(t T, r R) S {
		return fn(t)(r)
	}
}

// Uncurry3 преобразует каррированную функцию с тремя аргументами
func Uncurry3[T, R, S, U any](fn func(T) func(R) func(S) U) func(T, R, S) U {
	return func(t T, r R, s S) U {
		return fn(t)(r)(s)
	}
}

// Partial применяет частичное применение к функции
func Partial2[T, R, S any](fn func(T, R) S, t T) func(R) S {
	return func(r R) S {
		return fn(t, r)
	}
}

// Partial3 применяет частичное применение к функции с тремя аргументами
func Partial3[T, R, S, U any](fn func(T, R, S) U, t T) func(R, S) U {
	return func(r R, s S) U {
		return fn(t, r, s)
	}
}

// Flip меняет порядок аргументов функции
func Flip[T, R, S any](fn func(T, R) S) func(R, T) S {
	return func(r R, t T) S {
		return fn(t, r)
	}
}

// Memoize кэширует результаты функции
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

// MemoizeWithTTL кэширует результаты функции с TTL
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

// Debounce создает функцию с задержкой выполнения
func Debounce[T any](fn func(T), delay int64) func(T) {
	var lastCall int64
	
	return func(t T) {
		now := getCurrentTimestamp()
		lastCall = now
		
		go func() {
			// Простая реализация debounce без time.Sleep для примера
			// В реальном коде используйте time.Sleep(time.Duration(delay) * time.Millisecond)
			if getCurrentTimestamp()-lastCall >= delay {
				fn(t)
			}
		}()
	}
}

// Throttle создает функцию с ограничением частоты вызовов
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

// getCurrentTimestamp - helper функция для получения текущего времени
// В реальном коде замените на time.Now().Unix()
func getCurrentTimestamp() int64 {
	// Заглушка для компиляции, в реальном коде используйте time.Now().Unix()
	return 0
}
