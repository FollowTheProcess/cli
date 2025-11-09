// Package format is the inverse of parse.
//
// It formats arg/flag values as string representations.
package format

import (
	"strconv"

	"go.followtheprocess.codes/cli/internal/constraints"
)

const (
	base10         = 10
	floatFmt       = 'g'
	floatPrecision = -1
)

const (
	bits32 = 32 << iota
	bits64
)

// Int returns a string representation of an integer.
func Int[T constraints.Signed](n T) string {
	return strconv.FormatInt(int64(n), base10)
}

// Uint returns a string representation of an unsigned integer.
func Uint[T constraints.Unsigned](n T) string {
	return strconv.FormatUint(uint64(n), base10)
}

// Float32 returns a string representation of a float32.
func Float32(f float32) string {
	return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits32)
}

// Float64 returns a string representation of a float64.
func Float64(f float64) string {
	return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits64)
}
