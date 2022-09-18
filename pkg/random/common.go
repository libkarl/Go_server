package random

import (
	"math/rand"
)

// Bytes generates random bytes.
func Bytes(n int) []byte {
	temp := make([]byte, n)
	//nolint:
	rand.Read(temp)
	return temp
}

// OneOf gets randomly one of the provied values.
func OneOf(values ...interface{}) interface{} {
	return values[rand.Intn(len(values)-1)]
}

// OneOfStrings gets randomly one of the provided string values.
func OneOfStrings(values ...string) string {
	return values[rand.Intn(len(values)-1)]
}

// OneOfInts gets randomly one of the provided integer values.
func OneOfInts(values ...int) int {
	return values[rand.Intn(len(values)-1)]
}
