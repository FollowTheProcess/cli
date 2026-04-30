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
	slice          = "[]"

	// Capacity hints used to pre-size []byte buffers in the slice formatters.
	// The "brackets" pair is the leading '[' and trailing ']'.
	//
	// The per-element hints are good enough default cases to minimise the
	// buffer growing and re-allocating.
	bracketsCap    = 2
	intElemHint    = 4 // "-12, "
	floatElemHint  = 8 // "-1.234, "
	boolElemHint   = 7 // "false, "
	stringElemHint = 4 // surrounding quotes plus ", "
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
	TypeURL          = "url"
	TypeBool         = "bool"
	TypeBytesHex     = "bytesHex"
	TypeTime         = "time"
	TypeDuration     = "duration"
	TypeIP           = "ip"
	TypeIntSlice     = slice + TypeInt
	TypeInt8Slice    = slice + TypeInt8
	TypeInt16Slice   = slice + TypeInt16
	TypeInt32Slice   = slice + TypeInt32
	TypeInt64Slice   = slice + TypeInt64
	TypeUintSlice    = slice + TypeUint
	TypeUint16Slice  = slice + TypeUint16
	TypeUint32Slice  = slice + TypeUint32
	TypeUint64Slice  = slice + TypeUint64
	TypeFloat32Slice = slice + TypeFloat32
	TypeFloat64Slice = slice + TypeFloat64
	TypeStringSlice  = slice + TypeString
	TypeURLSlice     = slice + TypeURL
)

// True is the literal boolean true as a string.
//
// We check for "true" when scanning boolean flags, interpreting
// default values and when flags have a NoArgValue.
const True = "true"

// Nil is the string representation of a Go nil value.
const Nil = "<nil>"

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

// Slice returns a string representation of a slice.
//
// It will return a bracketed, comma separated list of items. If T is
// a string, the items will be quoted.
//
//	Slice([]int{1, 2, 3, 4}) // "[1, 2, 3, 4]"
//	Slice([]string{"one", "two", "three"}) // `["one", "two", "three"]`
func Slice[T any](s []T) string {
	if len(s) == 0 {
		return slice
	}

	switch v := any(s).(type) {
	case []string:
		return formatStringSlice(v)
	case []bool:
		return formatBoolSlice(v)
	case []int:
		return formatSignedSlice(v)
	case []int8:
		return formatSignedSlice(v)
	case []int16:
		return formatSignedSlice(v)
	case []int32:
		return formatSignedSlice(v)
	case []int64:
		return formatSignedSlice(v)
	case []uint:
		return formatUnsignedSlice(v)
	case []uint16:
		return formatUnsignedSlice(v)
	case []uint32:
		return formatUnsignedSlice(v)
	case []uint64:
		return formatUnsignedSlice(v)
	case []float32:
		return formatFloat32Slice(v)
	case []float64:
		return formatFloat64Slice(v)
	default:
		return slice
	}
}

func formatSignedSlice[T constraints.Signed](s []T) string {
	buf := make([]byte, 0, bracketsCap+len(s)*intElemHint)
	buf = append(buf, '[')
	buf = strconv.AppendInt(buf, int64(s[0]), base10)

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendInt(buf, int64(e), base10)
	}

	buf = append(buf, ']')

	return string(buf)
}

func formatUnsignedSlice[T constraints.Unsigned](s []T) string {
	buf := make([]byte, 0, bracketsCap+len(s)*intElemHint)
	buf = append(buf, '[')
	buf = strconv.AppendUint(buf, uint64(s[0]), base10)

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendUint(buf, uint64(e), base10)
	}

	buf = append(buf, ']')

	return string(buf)
}

func formatFloat32Slice(s []float32) string {
	buf := make([]byte, 0, bracketsCap+len(s)*floatElemHint)
	buf = append(buf, '[')
	buf = strconv.AppendFloat(buf, float64(s[0]), floatFmt, floatPrecision, bits32)

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendFloat(buf, float64(e), floatFmt, floatPrecision, bits32)
	}

	buf = append(buf, ']')

	return string(buf)
}

func formatFloat64Slice(s []float64) string {
	buf := make([]byte, 0, bracketsCap+len(s)*floatElemHint)
	buf = append(buf, '[')
	buf = strconv.AppendFloat(buf, s[0], floatFmt, floatPrecision, bits64)

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendFloat(buf, e, floatFmt, floatPrecision, bits64)
	}

	buf = append(buf, ']')

	return string(buf)
}

func formatStringSlice(s []string) string {
	capacity := bracketsCap
	for _, e := range s {
		capacity += len(e) + stringElemHint
	}

	buf := make([]byte, 0, capacity)
	buf = append(buf, '[')
	buf = strconv.AppendQuote(buf, s[0])

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendQuote(buf, e)
	}

	buf = append(buf, ']')

	return string(buf)
}

func formatBoolSlice(s []bool) string {
	buf := make([]byte, 0, bracketsCap+len(s)*boolElemHint)
	buf = append(buf, '[')
	buf = strconv.AppendBool(buf, s[0])

	for _, e := range s[1:] {
		buf = append(buf, ", "...)
		buf = strconv.AppendBool(buf, e)
	}

	buf = append(buf, ']')

	return string(buf)
}
