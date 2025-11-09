// Package format is the inverse of parse.
//
// It formats arg/flag values as string representations.
package format

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.followtheprocess.codes/cli/internal/constraints"
)

const (
	base10         = 10
	floatFmt       = 'g'
	floatPrecision = -1
	slice          = "[]"
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
)

// True is the literal boolean true as a string.
//
// We check for "true" when scanning boolean flags, interpreting
// default values and when flags have a NoArgValue.
const True = "true"

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
	length := len(s)

	if length == 0 {
		// If it's empty or nil, avoid doing the work below
		// and just return "[]"
		return slice
	}

	builder := &strings.Builder{}
	builder.WriteByte('[')

	typ := reflect.TypeFor[T]().Kind()

	first := fmt.Sprintf("%v", s[0])
	if typ == reflect.String {
		first = strconv.Quote(first)
	}

	builder.WriteString(first)

	for _, element := range s[1:] {
		builder.WriteString(", ")

		str := fmt.Sprintf("%v", element)
		if typ == reflect.String {
			// If it's a string, quote it
			str = strconv.Quote(str)
		}

		builder.WriteString(str)
	}

	builder.WriteByte(']')

	return builder.String()
}
