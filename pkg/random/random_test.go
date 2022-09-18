package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEqualValuesInt(t *testing.T) {
	intr1 := Int(0, 20000000)
	intr2 := Int(0, 20000000)
	assert.Equal(t, true, intr1 != intr2)
}

func TestNoEqualValuesString(t *testing.T) {
	intr1 := String(100)
	intr2 := String(100)
	assert.Equal(t, true, intr1 != intr2)
}

func TestRandFloatInRange(t *testing.T) {
	f1 := Float64(10, 11)
	f2 := Float64(11, 12)
	assert.InDelta(t, f1, 10, 1)
	assert.InDelta(t, f2, 11, 1)
}
