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

// Type names.
const (
	TypeInt          = "int"
	TypeInt8         = "int8"
	TypeInt16        = "int16"
	TypeInt32        = "int32"
	TypeInt64        = "int64"
	TypeUint         = "uint"
	TypeCount        = "count"
	TypeUint8        = "uint8"
	TypeUint16       = "uint16"
	TypeUint32       = "uint32"
	TypeUint64       = "uint64"
	TypeUintptr      = "uintptr"
	TypeFloat32      = "float32"
	TypeFloat64      = "float64"
	TypeString       = "string"
	TypeBool         = "bool"
	TypeBytesHex     = "bytesHex"
	TypeTime         = "time"
	TypeDuration     = "duration"
	TypeIP           = "ip"
	TypeIntSlice     = "[]int"
	TypeInt8Slice    = "[]int8"
	TypeInt16Slice   = "[]int16"
	TypeInt32Slice   = "[]int32"
	TypeInt64Slice   = "[]int64"
	TypeUintSlice    = "[]uint"
	TypeUint16Slice  = "[]uint16"
	TypeUint32Slice  = "[]uint32"
	TypeUint64Slice  = "[]uint64"
	TypeFloat32Slice = "[]float32"
	TypeFloat64Slice = "[]float64"
	TypeStringSlice  = "[]string"
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
