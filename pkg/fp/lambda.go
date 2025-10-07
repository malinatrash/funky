package fp

import "github.com/google/uuid"

func UuidToString(uuid uuid.UUID) string {
	return uuid.String()
}

func StringToUuid(id string) *uuid.UUID {
	va, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &va
}


func IsEmptySlice[T any](s []T) bool {
	return len(s) == 0
}

func IsNotEmptySlice[T any](s []T) bool {
	return len(s) > 0
}

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func IsNotEmptyString(s string) bool {
	return len(s) > 0
}

func Not[T any](predicate func(T) bool) func(T) bool {
	return func(v T) bool {
		return !predicate(v)
	}
}
