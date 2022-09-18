package random

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

// Int will give you random number in range min max
func Int(min int, max int) int {
	return min + random.Intn(max-min)
}

// Uint will give you random number in range min max
func Uint(min int, max int) uint {
	return uint(Int(min, max))
}

// Float64 random float64
func Float64(fls ...float64) float64 {
	if len(fls) == 2 {
		return fls[0] + random.Float64()*(fls[1]-fls[0])
	}
	return random.Float64()
}

// Float64Decimal random float64 decimal
func Float64Decimal(fls ...float64) decimal.Decimal {
	ff := Float64(fls...)
	return decimal.NewFromFloat(ff)
}

// Float64DecimalFixed generates random decimal based on the provided numeric 'precision'
// which defines the number of the digits after comma. The 'fls' provides maximum and minimum value for the random float.
func Float64DecimalFixed(precision int, fls ...float64) decimal.Decimal {
	prec := strconv.FormatInt(int64(precision), 10)
	str := fmt.Sprintf("%."+prec+"f", Float64(fls...))
	f, _ := strconv.ParseFloat(str, 64)
	return decimal.NewFromFloat(f)
}

// Float64Fixed random float64
func Float64Fixed(precision int, fls ...float64) float64 {
	prec := strconv.FormatInt(int64(precision), 10)
	str := fmt.Sprintf("%."+prec+"f", Float64(fls...))
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
