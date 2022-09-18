package ptrto

import "time"

// Int gets the pointer to the 'i' int.
func Int(i int) *int {
	return &i
}

// String gets the pointer to the 'v' string.
func String(v string) *string {
	return &v
}

// Time gets the pointer to the 't' time.
func Time(t time.Time) *time.Time {
	return &t
}

// Uint gets the pointer to the 'i' uint.
func Uint(i uint) *uint {
	return &i
}
